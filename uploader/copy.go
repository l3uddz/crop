package uploader

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (u *Uploader) Copy(additionalRcloneParams []string) error {
	// set variables
	extraParams := u.Config.RcloneParams.Copy
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	if globalParams := rclone.GetGlobalParams(rclone.GlobalCopyParams, u.Config.RcloneParams.GlobalCopy); globalParams != nil {
		extraParams = append(extraParams, globalParams...)
	}

	// iterate all remotes and run copy
	for _, remotePath := range u.Config.Remotes.Copy {
		// set variables
		attempts := 1

		// copy to remote
		for {
			// set log
			rLog := u.Log.WithFields(logrus.Fields{
				"copy_remote":     remotePath,
				"copy_local_path": u.Config.LocalFolder,
				"attempts":        attempts,
			})

			// get service account(s)
			serviceAccounts, err := u.RemoteServiceAccountFiles.GetServiceAccount(remotePath)
			if err != nil {
				return errors.WithMessagef(err,
					"aborting further copy attempts of %q due to serviceAccount exhaustion",
					u.Config.LocalFolder)
			}

			// display service account(s) being used
			if len(serviceAccounts) > 0 {
				for _, sa := range serviceAccounts {
					rLog.Infof("Using service account %q: %v", sa.RemoteEnvVar, sa.ServiceAccountPath)
				}
			}

			// copy
			rLog.Info("Copying...")
			success, exitCode, err := rclone.Copy(u.Config.LocalFolder, remotePath, serviceAccounts, extraParams)

			// check result
			if err != nil {
				rLog.WithError(err).Errorf("Failed unexpectedly...")
				return errors.WithMessagef(err, "copy failed unexpectedly with exit code: %v", exitCode)
			} else if success {
				// successful exit code
				break
			}

			// is this an exit code we can retry?
			switch exitCode {
			case rclone.ExitFatalError:
				// are we using service accounts?
				if len(serviceAccounts) == 0 {
					// we are not using service accounts, so mark this remote as banned
					if err := cache.SetBanned(stringutils.FromLeftUntil(remotePath, ":"), 25); err != nil {
						rLog.WithError(err).Errorf("Failed banning remote")
					}

					return fmt.Errorf("copy failed with exit code: %v", exitCode)
				}

				// ban service account(s) used
				for _, sa := range serviceAccounts {
					if err := cache.SetBanned(sa.ServiceAccountPath, 25); err != nil {
						rLog.WithError(err).Error("Failed banning service account, cannot try again...")
						return fmt.Errorf("failed banning service account: %v", sa.ServiceAccountPath)
					}
				}

				// attempt copy again
				rLog.Warnf("Copy failed with retryable exit code %v, trying again...", exitCode)
				attempts++
				continue
			default:
				return fmt.Errorf("failed and cannot proceed with exit code: %v", exitCode)
			}
		}
	}

	return nil
}
