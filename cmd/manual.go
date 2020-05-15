package cmd

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/stringutils"
	"github.com/l3uddz/crop/syncer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var (
	flagSrc      string
	flagDest     string
	flagSaFolder string
	flagDedupe   bool
	flagCopy     bool
	flagSync     bool
)

var manualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Perform a manual copy/sync task",
	Long:  `This command can be used to trigger a copy/sync without requiring configuration changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)
		defer cache.Close()

		// determine destination remotes
		syncRemotes := make([]string, 0)
		copyRemotes := make([]string, 0)

		switch {
		case flagCopy && flagSync:
			log.Fatal("You should must a single mode to use, --sync / --copy")
		case flagCopy:
			copyRemotes = append(copyRemotes, flagDest)
		case flagSync:
			syncRemotes = append(syncRemotes, flagDest)
		default:
			log.Fatal("You must specify a mode to use, --sync / --copy")
		}

		// create remote to service account map
		remoteSaFolders := make(map[string]string)

		switch flagSaFolder != "" {
		case true:
			if strings.Contains(flagSrc, ":") {
				// source is a remote
				srcRemote := stringutils.FromLeftUntil(flagSrc, ":")
				log.Debugf("Using service account folder for %q: %v", srcRemote, flagSaFolder)
				remoteSaFolders[srcRemote] = flagSaFolder
			}

			if strings.Contains(flagDest, ":") {
				// dest is a remote
				dstRemote := stringutils.FromLeftUntil(flagDest, ":")
				log.Debugf("Using service account folder for %q: %v", dstRemote, flagSaFolder)
				remoteSaFolders[dstRemote] = flagSaFolder
			}

		default:
			break
		}

		// create syncer config
		syncerConfig := config.SyncerConfig{
			Name:         "manual",
			Enabled:      true,
			SourceRemote: flagSrc,
			Remotes: config.SyncerRemotes{
				Copy: copyRemotes,
				Sync: syncRemotes,
			},
			RcloneParams: config.SyncerRcloneParams{
				Copy: args,
				Sync: args,
				Dedupe: []string{
					"--tpslimit=5",
				},
			},
		}

		if flagDedupe {
			// dedupe was enabled
			syncerConfig.Remotes.Dedupe = []string{
				flagDest,
			}
		}

		// create a config structure for manual sync
		cfg := config.Configuration{
			Rclone: config.RcloneConfig{
				Path:                  config.Config.Rclone.Path,
				Config:                config.Config.Rclone.Config,
				Stats:                 config.Config.Rclone.Stats,
				DryRun:                config.Config.Rclone.DryRun,
				ServiceAccountRemotes: remoteSaFolders,
			},
			Uploader: nil,
			Syncer: []config.SyncerConfig{
				syncerConfig,
			},
		}

		// create syncer
		sync, err := syncer.New(&cfg, &syncerConfig, syncerConfig.Name, 1)
		if err != nil {
			log.WithError(err).Fatal("Failed initializing syncer, skipping...")
		}

		// load service accounts
		serviceAccountCount := sync.RemoteServiceAccountFiles.ServiceAccountsCount()
		if serviceAccountCount > 0 {
			sync.Log.WithField("found_files", serviceAccountCount).Info("Loaded service accounts")
		} else {
			// no service accounts were loaded
			// check to see if any of the copy or sync remote(s) are banned
			banned, expiry := rclone.AnyRemotesBanned(sync.Config.Remotes.Copy)
			if banned && !expiry.IsZero() {
				// one of the copy remotes is banned, abort
				sync.Log.WithFields(logrus.Fields{
					"expires_time": expiry,
					"expires_in":   humanize.Time(expiry),
				}).Fatal("Cannot proceed as a copy remote is banned")
			}

			banned, expiry = rclone.AnyRemotesBanned(sync.Config.Remotes.Sync)
			if banned && !expiry.IsZero() {
				// one of the sync remotes is banned, abort
				sync.Log.WithFields(logrus.Fields{
					"expires_time": expiry,
					"expires_in":   humanize.Time(expiry),
				}).Fatal("Cannot proceed as a sync remote is banned")
			}
		}

		log.Info("Syncer commencing...")

		// perform sync
		if err := performSync(sync); err != nil {
			sync.Log.WithError(err).Fatal("Error occurred while running syncer, skipping...")
		}

		log.Info("Finished!")
	},
}

func init() {
	rootCmd.AddCommand(manualCmd)

	manualCmd.Flags().StringVar(&flagSrc, "src", "", "Source")
	manualCmd.Flags().StringVar(&flagDest, "dst", "", "Destination")

	_ = manualCmd.MarkFlagRequired("from")
	_ = manualCmd.MarkFlagRequired("dest")

	manualCmd.Flags().StringVar(&flagSaFolder, "sa", "", "Service account folder")

	manualCmd.Flags().BoolVar(&flagCopy, "copy", false, "Copy to destination")
	manualCmd.Flags().BoolVar(&flagSync, "sync", false, "Sync to destination")
	manualCmd.Flags().BoolVar(&flagDedupe, "dedupe", false, "Dedupe destination")
}
