package cmd

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/uploader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yale8848/gorpool"
	"strings"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Perform cleans associated to uploader(s)",
	Long:  `This command can be used to trigger a clean associated with uploader(s).`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)

		// iterate uploader's
		for uploaderName, uploaderConfig := range config.Config.Uploader {
			// skip disabled uploader(s)
			if !uploaderConfig.Enabled {
				log.WithField("uploader", uploaderName).Trace("Skipping disabled uploader")
				continue
			}

			// skip uploader specific chosen
			if flagUploader != "" && !strings.EqualFold(uploaderName, flagUploader) {
				log.WithField("uploader", uploaderName).Tracef("Skipping uploader as not: %q",
					flagUploader)
				continue
			}

			log := log.WithField("uploader", uploaderName)

			// create uploader
			upload, err := uploader.New(config.Config, &uploaderConfig, uploaderName)
			if err != nil {
				log.WithField("uploader", uploaderName).WithError(err).
					Error("Failed initializing uploader, skipping...")
				continue
			}

			log.Info("Clean commencing...")

			// perform upload
			if err := performClean(upload); err != nil {
				upload.Log.WithError(err).Error("Error occurred while running clean, skipping...")
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().StringVarP(&flagUploader, "uploader", "u", "", "Run for a specific uploader.")
}

func performClean(u *uploader.Uploader) error {
	u.Log.Info("Running cleans...")

	/* Cleans */
	if u.Config.Hidden.Enabled {
		// set worker count
		workers := u.Config.Hidden.Workers
		if workers == 0 {
			workers = 8
		}

		// create worker pool
		gp := gorpool.NewPool(workers, 0).
			Start().
			EnableWaitForAll(true)

		// queue clean tasks
		err := u.PerformCleans(gp)
		if err != nil {
			return errors.Wrap(err, "failed clearing remotes")
		}
	}

	u.Log.Info("Finished cleans!")
	return nil
}
