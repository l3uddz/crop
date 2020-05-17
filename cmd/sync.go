package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/syncer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"sync"
	"time"
)

var (
	flagSyncer      string
	flagParallelism int
	flagNoDedupe    bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Perform syncer task(s)",
	Long:  `This command can be used to trigger a sync.`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)
		defer cache.Close()

		// create workers
		var wg sync.WaitGroup
		jobs := make(chan *syncer.Syncer, len(config.Config.Syncer))

		for w := 1; w <= flagParallelism; w++ {
			wg.Add(1)
			go worker(&wg, jobs)
		}

		// iterate syncer's
		started := time.Now().UTC()

		for _, syncerConfig := range config.Config.Syncer {
			syncerConfig := syncerConfig

			slog := log.WithField("syncer", syncerConfig.Name)

			// skip disabled syncer(s)
			if !syncerConfig.Enabled {
				slog.Debug("Skipping disabled syncer")
				continue
			}

			// skip syncer specific chosen
			if flagSyncer != "" && !strings.EqualFold(syncerConfig.Name, flagSyncer) {
				slog.Debugf("Skipping syncer as not: %q", flagSyncer)
				continue
			}

			// create syncer
			syncr, err := syncer.New(config.Config, &syncerConfig, syncerConfig.Name, flagParallelism)
			if err != nil {
				slog.WithError(err).Error("Failed initializing syncer, skipping...")
				continue
			}

			serviceAccountCount := syncr.RemoteServiceAccountFiles.ServiceAccountsCount()
			if serviceAccountCount > 0 {
				syncr.Log.WithField("found_files", serviceAccountCount).Info("Loaded service accounts")
			} else {
				// no service accounts were loaded
				// check to see if any of the copy or sync remote(s) are banned
				banned, expiry := rclone.AnyRemotesBanned(syncr.Config.Remotes.Copy)
				if banned && !expiry.IsZero() {
					// one of the copy remotes is banned, abort
					syncr.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with sync as a copy remote is banned")
					continue
				}

				banned, expiry = rclone.AnyRemotesBanned(syncr.Config.Remotes.Sync)
				if banned && !expiry.IsZero() {
					// one of the sync remotes is banned, abort
					syncr.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with sync as a sync remote is banned")
					continue
				}
			}

			// queue sync job
			jobs <- syncr
		}

		// wait for all syncers to finish
		log.Info("Waiting for syncer(s) to finish")
		close(jobs)
		wg.Wait()

		log.Infof("Finished in: %v", humanize.RelTime(started, time.Now().UTC(), "", ""))
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringVarP(&flagSyncer, "syncer", "s", "", "Run for a specific syncer")
	syncCmd.Flags().IntVarP(&flagParallelism, "parallelism", "p", 1, "Max parallel syncers")

	syncCmd.Flags().BoolVar(&flagNoDedupe, "no-dedupe", false, "Ignore dedupe tasks for syncer")
}

func worker(wg *sync.WaitGroup, jobs <-chan *syncer.Syncer) {
	defer wg.Done()

	for j := range jobs {
		// perform syncer job
		if err := performSync(j); err != nil {
			j.Log.WithError(err).Error("Error occurred while running syncer, skipping...")
		}
	}
}

func performSync(s *syncer.Syncer) error {
	s.Log.Info("Running...")

	var gcloneParams []string
	if strings.Contains(s.GlobalConfig.Rclone.Path, "gclone") &&
		s.RemoteServiceAccountFiles.ServiceAccountsCount() > 0 {
		// start web-server
		s.Ws.Run()
		defer s.Ws.Stop()

		gcloneParams = append(gcloneParams,
			"--drive-service-account-url",
			fmt.Sprintf("http://%s:%d", s.Ws.Host, s.Ws.Port),
		)
	}

	/* Copies */
	if len(s.Config.Remotes.Copy) > 0 {
		s.Log.Info("Running copies...")

		if err := s.Copy(gcloneParams); err != nil {
			return errors.WithMessage(err, "failed performing all copies")
		}

		s.Log.Info("Finished copies!")
	}

	/* Sync */
	if len(s.Config.Remotes.Sync) > 0 {
		s.Log.Info("Running syncs...")

		if err := s.Sync(gcloneParams); err != nil {
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
	if !flagNoDedupe && len(s.Config.Remotes.Dedupe) > 0 {
		s.Log.Info("Running dedupes...")

		if err := s.Dedupe(nil); err != nil {
			return errors.WithMessage(err, "failed performing all dedupes")
		}

		s.Log.Info("Finished dedupes!")
	}

	s.Log.Info("Finished!")
	return nil
}
