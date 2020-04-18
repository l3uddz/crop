package syncer

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Syncer) Copy(additionalRcloneParams []string) error {
	// set variables
	copyParams := s.Config.RcloneParams.Copy
	if additionalRcloneParams != nil {
		copyParams = append(copyParams, additionalRcloneParams...)
	}

	// iterate all remotes and run copy
	for _, remotePath := range s.Config.Remotes.Copy {
		// set variables
		attempts := 1
		rLog := s.Log.WithFields(logrus.Fields{
			"copy_remote":   remotePath,
			"source_remote": s.Config.SourceRemote,
			"attempts":      attempts,
		})

		// copy to remote
		for {
			// get service account file
			var serviceAccount *pathutils.Path
			var err error

			if s.ServiceAccountCount > 0 {
				serviceAccount, err = rclone.GetAvailableServiceAccount(s.ServiceAccountFiles)
				if err != nil {
					return errors.WithMessagef(err,
						"aborting further copy attempts of %q due to serviceAccount exhaustion",
						s.Config.SourceRemote)
				}

				// reset log
				rLog = s.Log.WithFields(logrus.Fields{
					"copy_remote":     remotePath,
					"source_remote":   s.Config.SourceRemote,
					"attempts":        attempts,
					"service_account": serviceAccount.RealPath,
				})
			}

			// copy
			rLog.Info("Copying...")
			success, exitCode, err := rclone.Copy(s.Config.SourceRemote, remotePath, serviceAccount, copyParams)

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
			case rclone.EXIT_FATAL_ERROR:
				// are we using service accounts?
				if s.ServiceAccountCount == 0 {
					// we are not using service accounts, so mark this remote as banned
					if err := cache.Set(stringutils.FromLeftUntil(remotePath, ":"),
						time.Now().UTC().Add(25*time.Hour)); err != nil {
						rLog.WithError(err).Errorf("Failed banning remote")
					}

					return fmt.Errorf("copy failed with exit code: %v", exitCode)
				}

				// ban this service account
				if err := cache.Set(serviceAccount.RealPath, time.Now().UTC().Add(25*time.Hour)); err != nil {
					rLog.WithError(err).Error("Failed banning service account, cannot try again...")
					return fmt.Errorf("failed banning service account: %v", serviceAccount.RealPath)
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
