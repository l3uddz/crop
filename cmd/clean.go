package cmd

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/uploader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yale8848/gorpool"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [UPLOADER]",
	Short: "Perform uploader clean",
	Long:  `This command can be used to trigger an uploader clean.`,

	Args: cobra.ExactArgs(1),
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
