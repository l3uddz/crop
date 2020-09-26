package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/uploader"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/disk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var (
	flagNoCheck bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Perform uploader task(s)",
	Long:  `This command can be used to trigger an uploader check, clean & upload.`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(true)
		defer cache.Close()
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

			serviceAccountCount := upload.RemoteServiceAccountFiles.ServiceAccountsCount()
			if serviceAccountCount > 0 {
				upload.Log.WithField("found_files", serviceAccountCount).Info("Loaded service accounts")
			} else {
				// no service accounts were loaded
				// check to see if any of the copy or move remote(s) are banned
				banned, expiry := rclone.AnyRemotesBanned(upload.Config.Remotes.Copy)
				if banned && !expiry.IsZero() {
					// one of the copy remotes is banned, abort
					upload.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with upload as a copy remote is banned")
					continue
				}

				banned, expiry = rclone.AnyRemotesBanned([]string{upload.Config.Remotes.Move})
				if banned && !expiry.IsZero() {
					// the move remote is banned, abort
					upload.Log.WithFields(logrus.Fields{
						"expires_time": expiry,
						"expires_in":   humanize.Time(expiry),
					}).Warn("Cannot proceed with upload as the move remote is banned")
					continue
				}
			}

			log.Info("Uploader commencing...")

			// refresh details about files to upload
			if err := upload.RefreshLocalFiles(); err != nil {
				upload.Log.WithError(err).Error("Failed refreshing details of files to upload")
				continue
			}

			if len(upload.LocalFiles) == 0 {
				// there are no files to upload
				upload.Log.Info("There were no files found, skipping...")
				continue
			}

			// check if upload criteria met
			forced := false

			if !flagNoCheck {
				// no check was not enabled
				res, err := upload.Check()
				if err != nil {
					upload.Log.WithError(err).Error("Failed checking if uploader check conditions met, skipping...")
					continue
				}

				if !res.Passed {
					// get free disk space
					freeDiskSpace := "Unknown"
					du, err := disk.Usage(upload.Config.LocalFolder)
					if err == nil {
						freeDiskSpace = humanize.IBytes(du.Free)
					}

					// check available disk space
					switch {
					case err != nil && upload.Config.Check.MinFreeSpace > 0:
						// error checking free space
						upload.Log.WithError(err).Errorf("Failed checking available free space for: %q",
							upload.Config.LocalFolder)
					case err == nil && du.Free < upload.Config.Check.MinFreeSpace:
						// free space has gone below the free space threshold
						forced = true
						upload.Log.WithFields(logrus.Fields{
							"until":     res.Info,
							"free_disk": freeDiskSpace,
						}).Infof("Upload conditions not met, however, proceeding as free space below %s",
							humanize.IBytes(upload.Config.Check.MinFreeSpace))
					default:
						break
					}

					if !forced {
						upload.Log.WithFields(logrus.Fields{
							"until":     res.Info,
							"free_disk": freeDiskSpace,
						}).Info("Upload conditions not met, skipping...")
						continue
					}

					// the upload was forced as min_free_size was met
				}
			}

			// perform upload
			if err := performUpload(upload, forced); err != nil {
				upload.Log.WithError(err).Error("Error occurred while running uploader, skipping...")
				continue
			}
		}

		log.Infof("Finished in: %v", humanize.RelTime(started, time.Now().UTC(), "", ""))
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&flagUploader, "uploader", "u", "", "Run for a specific uploader")

	uploadCmd.Flags().BoolVar(&flagNoCheck, "no-check", false, "Ignore check and run")
	uploadCmd.Flags().BoolVar(&flagNoDedupe, "no-dedupe", false, "Ignore dedupe tasks for uploader")
}

func performUpload(u *uploader.Uploader, forced bool) error {
	u.Log.Info("Running...")

	var liveRotateParams []string

	if u.GlobalConfig.Rclone.LiveRotate && u.RemoteServiceAccountFiles.ServiceAccountsCount() > 0 {
		// start web-server
		u.Ws.Run()
		defer u.Ws.Stop()

		liveRotateParams = append(liveRotateParams,
			"--drive-service-account-url",
			fmt.Sprintf("http://%s:%d", u.Ws.Host, u.Ws.Port),
		)
	}

	/* Cleans */
	if u.Config.Hidden.Enabled {
		err := performClean(u)
		if err != nil {
			return errors.Wrap(err, "failed clearing remotes")
		}
	}

	/* Generate Additional Rclone Params */
	var additionalRcloneParams []string

	switch forced {
	case false:
		if !flagNoCheck || u.Config.Check.Forced {
			// if no-check is false (default) or check is forced via config, include check params
			additionalRcloneParams = u.CheckRcloneParams()
		}
	default:
		break
	}

	// add live rotate params set
	if len(liveRotateParams) > 0 {
		additionalRcloneParams = append(additionalRcloneParams, liveRotateParams...)
	}

	/* Copies */
	if len(u.Config.Remotes.Copy) > 0 {
		u.Log.Info("Running copies...")

		if err := u.Copy(additionalRcloneParams); err != nil {
			return errors.WithMessage(err, "failed performing all copies")
		}

		u.Log.Info("Finished copies!")
	}

	/* Move */
	if len(u.Config.Remotes.Move) > 0 {
		u.Log.Info("Running move...")

		if err := u.Move(false, additionalRcloneParams); err != nil {
			return errors.WithMessage(err, "failed performing move")
		}

		u.Log.Info("Finished move!")
	}

	/* Move Server Side */
	if len(u.Config.Remotes.MoveServerSide) > 0 {
		u.Log.Info("Running move server-sides...")

		if err := u.Move(true, nil); err != nil {
			return errors.WithMessage(err, "failed performing server-side moves")
		}

		u.Log.Info("Finished move server-sides!")
	}

	/* Dedupe */
	if !flagNoDedupe && len(u.Config.Remotes.Dedupe) > 0 {
		u.Log.Info("Running dedupes...")

		if err := u.Dedupe(nil); err != nil {
			return errors.WithMessage(err, "failed performing dedupes")
		}

		u.Log.Info("Finished dupes!")
	}

	u.Log.Info("Finished!")
	return nil
}
