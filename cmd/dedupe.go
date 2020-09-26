package cmd

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/uploader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "Perform dedupe associated with uploader(s)",
	Long:  `This command can be used to trigger a dedupe associated with uploader(s).`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)
		defer releaseFileLock()

		// iterate uploader's
		started := time.Now().UTC()

		for _, uploaderConfig := range config.Config.Uploader {
			log := log.WithField("uploader", uploaderConfig.Name)

			// skip disabled uploader(s)
			if !uploaderConfig.Enabled {
				log.Debug("Skipping disabled uploader")
				continue
			}

			// skip uploader specific chosen
			if flagUploader != "" && !strings.EqualFold(uploaderConfig.Name, flagUploader) {
				log.Debugf("Skipping uploader as not: %q", flagUploader)
				continue
			}

			// create uploader
			upload, err := uploader.New(config.Config, &uploaderConfig, uploaderConfig.Name)
			if err != nil {
				log.WithError(err).Error("Failed initializing uploader, skipping...")
				continue
			}

			log.Info("Dedupe commencing...")

			// perform upload
			if err := performDedupe(upload); err != nil {
				upload.Log.WithError(err).Error("Error occurred while running dedupe, skipping...")
				continue
			}
		}

		log.Infof("Finished in: %v", humanize.RelTime(started, time.Now().UTC(), "", ""))
	},
}

func init() {
	rootCmd.AddCommand(dedupeCmd)

	dedupeCmd.Flags().StringVarP(&flagUploader, "uploader", "u", "", "Run for a specific uploader")
}

func performDedupe(u *uploader.Uploader) error {
	u.Log.Info("Running dedupe...")

	/* Dedupe */
	err := u.Dedupe(nil)
	if err != nil {
		return errors.Wrap(err, "failed dedupe remotes")
	}

	u.Log.Info("Finished dedupe!")
	return nil
}
