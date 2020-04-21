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

func (s *Syncer) Sync(additionalRcloneParams []string) error {
	// set variables
	extraParams := s.Config.RcloneParams.Sync
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	// add server side parameter
	extraParams = append(extraParams, "--drive-server-side-across-configs")

	// iterate all remotes and run sync
	for _, remotePath := range s.Config.Remotes.Sync {
		// set variables
		attempts := 1
		rLog := s.Log.WithFields(logrus.Fields{
			"sync_remote":   remotePath,
			"source_remote": s.Config.SourceRemote,
			"attempts":      attempts,
		})

		// sync to remote
		for {
			// get service account file(s)
			serviceAccounts, err := s.RemoteServiceAccountFiles.GetServiceAccount(s.Config.SourceRemote, remotePath)
			if err != nil {
				return errors.WithMessagef(err,
					"aborting further sync attempts of %q due to serviceAccount exhaustion",
					s.Config.SourceRemote)
			} else if len(serviceAccounts) > 0 {
				// reset log
				rLog = s.Log.WithFields(logrus.Fields{
					"sync_remote":   remotePath,
					"source_remote": s.Config.SourceRemote,
					"attempts":      attempts,
				})

				for _, sa := range serviceAccounts {
					rLog.Infof("Using service account %q: %v", sa.RemoteEnvVar, sa.ServiceAccountPath)
				}
			}

			// sync
			rLog.Info("Syncing...")
			success, exitCode, err := rclone.Sync(s.Config.SourceRemote, remotePath, serviceAccounts, extraParams)

			// check result
			if err != nil {
				rLog.WithError(err).Errorf("Failed unexpectedly...")
				return errors.WithMessagef(err, "sync failed unexpectedly with exit code: %v", exitCode)
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
					if err := cache.Set(stringutils.FromLeftUntil(remotePath, ":"),
						time.Now().UTC().Add(25*time.Hour)); err != nil {
						rLog.WithError(err).Errorf("Failed banning remote")
					}

					return fmt.Errorf("sync failed with exit code: %v", exitCode)
				}

				// ban this service account
				for _, sa := range serviceAccounts {
					if err := cache.Set(sa.ServiceAccountPath, time.Now().UTC().Add(25*time.Hour)); err != nil {
						rLog.WithError(err).Error("Failed banning service account, cannot try again...")
						return fmt.Errorf("failed banning service account: %v", sa.ServiceAccountPath)
					}
				}

				// attempt sync again
				rLog.Warnf("Sync failed with retryable exit code %v, trying again...", exitCode)
				attempts++
				continue
			default:
				return fmt.Errorf("failed and cannot proceed with exit code: %v", exitCode)
			}
		}
	}

	return nil
}
