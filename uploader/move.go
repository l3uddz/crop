package uploader

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (u *Uploader) Move(serverSide bool, additionalRcloneParams []string) error {
	var moveRemotes []rclone.RemoteInstruction
	var extraParams []string

	// create move instructions
	if serverSide {
		// this is a server side move
		for _, remote := range u.Config.Remotes.MoveServerSide {
			moveRemotes = append(moveRemotes, rclone.RemoteInstruction{
				From:       remote.From,
				To:         remote.To,
				ServerSide: true,
			})
		}

		extraParams = u.Config.RcloneParams.MoveServerSide
	} else {
		// this is a normal move (to only one location)
		moveRemotes = append(moveRemotes, rclone.RemoteInstruction{
			From:       u.Config.LocalFolder,
			To:         u.Config.Remotes.Move,
			ServerSide: false,
		})

		extraParams = u.Config.RcloneParams.Move
	}

	// set variables
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	// iterate all remotes and run move
	for _, move := range moveRemotes {
		// set variables
		attempts := 1

		// move to remote
		for {
			var serviceAccounts []*rclone.RemoteServiceAccount
			var err error

			// set log
			rLog := u.Log.WithFields(logrus.Fields{
				"move_to":   move.To,
				"move_from": move.From,
				"attempts":  attempts,
			})

			// get service account(s) for non server side move
			if !serverSide {
				serviceAccounts, err = u.RemoteServiceAccountFiles.GetServiceAccount(move.To)
				if err != nil {
					return errors.WithMessagef(err,
						"aborting further move attempts of %q due to serviceAccount exhaustion",
						move.From)
				}

				// display service accounts being used
				if len(serviceAccounts) > 0 {
					for _, sa := range serviceAccounts {
						rLog.Infof("Using service account %q: %v", sa.RemoteEnvVar, sa.ServiceAccountPath)
					}
				}
			}

			// move
			rLog.Info("Moving...")
			success, exitCode, err := rclone.Move(move.From, move.To, serviceAccounts, serverSide, extraParams)

			// check result
			if err != nil {
				rLog.WithError(err).Errorf("Failed unexpectedly...")
				return errors.WithMessagef(err, "move failed unexpectedly with exit code: %v", exitCode)
			}

			if success {
				// successful exit code
				break
			} else if serverSide {
				// server side moves will not use service accounts, so we will not retry...
				return fmt.Errorf("failed and cannot proceed with exit code: %v", exitCode)
			}

			// is this an exit code we can retry?
			switch exitCode {
			case rclone.ExitFatalError:
				// are we using service accounts?
				if len(serviceAccounts) == 0 {
					// we are not using service accounts, so mark this remote as banned (if non server side move)
					if !serverSide {
						// this was not a server side move, so lets ban the remote we are moving too
						if err := cache.Set(stringutils.FromLeftUntil(move.To, ":"),
							time.Now().UTC().Add(25*time.Hour)); err != nil {
							rLog.WithError(err).Errorf("Failed banning remote")
						}
					}

					return fmt.Errorf("move failed with exit code: %v", exitCode)
				}

				// ban the service account(s) used
				for _, sa := range serviceAccounts {
					if err := cache.Set(sa.ServiceAccountPath, time.Now().UTC().Add(25*time.Hour)); err != nil {
						rLog.WithError(err).Error("Failed banning service account, cannot try again...")
						return fmt.Errorf("failed banning service account: %v", sa.ServiceAccountPath)
					}
				}

				// attempt move again
				rLog.Warnf("Move failed with retryable exit code %v, trying again...", exitCode)
				attempts++
				continue
			default:
				return fmt.Errorf("failed and cannot proceed with exit code: %v", exitCode)
			}
		}
	}

	return nil
}
