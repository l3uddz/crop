package cmd

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/syncer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var (
	flagSyncer string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Perform syncer task(s)",
	Long:  `This command can be used to trigger a sync.`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)
		defer cache.Close()

		// iterate syncer's
		for _, syncerConfig := range config.Config.Syncer {
			log := log.WithField("syncer", syncerConfig.Name)

			// skip disabled syncer(s)
			if !syncerConfig.Enabled {
				log.Debug("Skipping disabled syncer")
				continue
			}

			// skip syncer specific chosen
			if flagSyncer != "" && !strings.EqualFold(syncerConfig.Name, flagSyncer) {
				log.Debugf("Skipping syncer as not: %q", flagSyncer)
				continue
			}

			// create syncer
			sync, err := syncer.New(config.Config, &syncerConfig, syncerConfig.Name)
			if err != nil {
				log.WithError(err).Error("Failed initializing syncer, skipping...")
				continue
			}

			serviceAccountCount := sync.RemoteServiceAccountFiles.ServiceAccountsCount()
			if serviceAccountCount > 0 {
				sync.Log.WithField("found_files", serviceAccountCount).Info("Loaded service accounts")
			} else {
				// no service accounts were loaded
				// check to see if any of the copy or move remote(s) are banned
				banned, expiry := rclone.AnyRemotesBanned(sync.Config.Remotes.Copy)
				if banned && !expiry.IsZero() {
					// one of the copy remotes is banned, abort
					sync.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with sync as a copy remote is banned")
					continue
				}

				banned, expiry = rclone.AnyRemotesBanned(sync.Config.Remotes.Sync)
				if banned && !expiry.IsZero() {
					// one of the sync remotes is banned, abort
					sync.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with sync as a sync remote is banned")
					continue
				}
			}

			log.Info("Syncer commencing...")

			// perform sync
			if err := performSync(sync); err != nil {
				sync.Log.WithError(err).Error("Error occurred while running syncer, skipping...")
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringVarP(&flagSyncer, "syncer", "s", "", "Run for a specific syncer")
}

func performSync(s *syncer.Syncer) error {
	s.Log.Info("Running...")

	/* Copies */
	if len(s.Config.Remotes.Copy) > 0 {
		s.Log.Info("Running copies...")

		if err := s.Copy(nil); err != nil {
			return errors.WithMessage(err, "failed performing all copies")
		}

		s.Log.Info("Finished copies!")
	}

	/* Sync */
	if len(s.Config.Remotes.Sync) > 0 {
		s.Log.Info("Running syncs...")

		if err := s.Sync(nil); err != nil {
			return errors.WithMessage(err, "failed performing all syncs")
		}

		s.Log.Info("Finished syncs!")
	}

	/* Move Server Side */
	if len(s.Config.Remotes.MoveServerSide) > 0 {
		s.Log.Info("Running move server-sides...")

		if err := s.Move(nil); err != nil {
			return errors.WithMessage(err, "failed performing server-side moves")
		}

		s.Log.Info("Finished move server-sides!")
	}

	/* Dedupe */
	if len(s.Config.Remotes.Dedupe) > 0 {
		s.Log.Info("Running dedupes...")

		if err := s.Dedupe(nil); err != nil {
			return errors.WithMessage(err, "failed performing all dedupes")
		}

		s.Log.Info("Finished dedupes!")
	}

	s.Log.Info("Finished!")
	return nil
}
