package syncer

import (
	"fmt"
	"github.com/l3uddz/crop/rclone"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Syncer) Dedupe(additionalRcloneParams []string) error {
	extraParams := rclone.FormattedParams(s.Config.RcloneParams.Dedupe)
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	// iterate all remotes and run dedupe
	for _, dedupeRemote := range s.Config.Remotes.Dedupe {
		// set variables
		rLog := s.Log.WithFields(logrus.Fields{
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
