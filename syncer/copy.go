package syncer

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Syncer) Copy(additionalRcloneParams []string, daisyChain bool) error {
	// set variables
	extraParams := s.Config.RcloneParams.Copy
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	if globalParams := rclone.GetGlobalParams(rclone.GlobalCopyParams, s.Config.RcloneParams.GlobalCopy); globalParams != nil {
		extraParams = append(extraParams, globalParams...)
	}

	// add server side parameter
	extraParams = append(extraParams, "--drive-server-side-across-configs")

	pos := 0
	srcRemote := s.Config.SourceRemote

	// iterate all remotes and run copy
	for _, remotePath := range s.Config.Remotes.Copy {
		// set variables
		attempts := 1

		// daisy
		if daisyChain && pos > 0 {
			srcRemote = s.Config.Remotes.Copy[pos-1]
		}
		pos++

		// copy to remote
		for {
			// set log
			rLog := s.Log.WithFields(logrus.Fields{
				"copy_remote":   remotePath,
				"source_remote": srcRemote,
				"attempts":      attempts,
			})

			// get service account file(s)
			serviceAccounts, err := s.RemoteServiceAccountFiles.GetServiceAccount(srcRemote, remotePath)
			if err != nil {
				return errors.WithMessagef(err,
					"aborting further copy attempts of %q due to serviceAccount exhaustion",
					srcRemote)
			}

			// display service account(s) being used
			if len(serviceAccounts) > 0 {
				for _, sa := range serviceAccounts {
					rLog.Infof("Using service account %q: %v", sa.RemoteEnvVar, sa.ServiceAccountPath)
				}
			}

			// copy
			rLog.Info("Copying...")
			success, exitCode, err := rclone.Copy(srcRemote, remotePath, serviceAccounts, extraParams)

			// check result
			if err != nil {
				rLog.WithError(err).Errorf("Failed unexpectedly...")
				return errors.WithMessagef(err, "copy failed unexpectedly with exit code: %v", exitCode)
			} else if success {
				// successful exit code
				if !s.Ws.Running {
					// web service is not running (no live rotate)
					rclone.RemoveServiceAccountsFromTempCache(serviceAccounts)
				}
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

				// ban this service account
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

		// sleep before moving on
		if daisyChain && pos < len(s.Config.Remotes.Copy) {
			s.Log.Info("Waiting 60 seconds before continuing...")
			time.Sleep(60 * time.Second)
		}
	}

	return nil
}
