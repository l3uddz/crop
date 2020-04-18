package cmd

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/uploader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yale8848/gorpool"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [UPLOADER]",
	Short: "Perform uploader task",
	Long:  `This command can be used to trigger an uploader check / upload.`,

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
			log.Info("Uploader commencing...")

			// create uploader
			upload, err := uploader.New(config.Config, &uploaderConfig, uploaderName)
			if err != nil {
				log.WithField("uploader", uploaderName).WithError(err).
					Error("Failed initializing uploader, skipping...")
				continue
			}

			if upload.ServiceAccountCount > 0 {
				upload.Log.WithField("found_files", upload.ServiceAccountCount).
					Info("Loaded service accounts")
			}

			// refresh details about files to upload
			if err := upload.RefreshLocalFiles(); err != nil {
				upload.Log.WithError(err).Error("Failed refreshing details of files to upload")
				continue
			}

			// check if upload criteria met
			if shouldUpload, err := upload.Check(); err != nil {
				upload.Log.WithError(err).Error("Failed checking if uploader check conditions met, skipping...")
				continue
			} else if !shouldUpload {
				upload.Log.Info("Upload conditions not met, skipping...")
				continue
			}

			// perform upload
			if err := performUpload(upload); err != nil {
				upload.Log.WithError(err).Error("Error occurred while running uploader, skipping...")
				continue
			}

			// clean local upload folder of empty directories
			upload.Log.Debug("Cleaning empty local directories...")
		}

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func performUpload(u *uploader.Uploader) error {
	/* Cleans */
	if u.Config.Hidden.Enabled {
		gp := gorpool.NewPool(config.Config.Core.Workers, 0).
			Start().
			EnableWaitForAll(true)

		err := u.PerformCleans(gp)
		if err != nil {
			return errors.Wrap(err, "failed clearing remotes")
		}
	}

	/* Copies */
	if len(u.Config.Remotes.Copy) > 0 {
		if err := u.Copy(); err != nil {
			return errors.WithMessage(err, "failed performing copies")
		}
	}

	/* Move */

	/* Move Server Side */

	return nil
}
