package uploader

import (
	"fmt"
	"github.com/l3uddz/crop/rclone"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (u *Uploader) Dedupe(additionalRcloneParams []string) error {
	extraParams := u.Config.RcloneParams.Dedupe
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	if globalParams := rclone.GetGlobalParams(rclone.GlobalDedupeParams, u.Config.RcloneParams.GlobalDedupe); globalParams != nil {
		extraParams = append(extraParams, globalParams...)
	}

	// iterate all remotes and run dedupe
	for _, dedupeRemote := range u.Config.Remotes.Dedupe {
		// set variables
		rLog := u.Log.WithFields(logrus.Fields{
			"dedupe_remote": dedupeRemote,
		})

		// dedupe remote
		rLog.Info("Deduping...")
		success, exitCode, err := rclone.Dedupe(dedupeRemote, extraParams)

		// check result
		if err != nil {
			rLog.WithError(err).Errorf("Failed unexpectedly...")
			return errors.WithMessagef(err, "dedupe failed unexpectedly with exit code: %v", exitCode)
		} else if success {
			// successful exit code
			continue
		}

		return fmt.Errorf("dedupe failed with exit code: %v", exitCode)
	}

	return nil
}
