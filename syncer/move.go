package syncer

import (
	"fmt"
	"github.com/l3uddz/crop/rclone"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Syncer) Move(additionalRcloneParams []string) error {
	moveRemotes := make([]rclone.RemoteInstruction, 0)

	// set variables
	for _, remote := range s.Config.Remotes.MoveServerSide {
		moveRemotes = append(moveRemotes, rclone.RemoteInstruction{
			From:       remote.From,
			To:         remote.To,
			ServerSide: true,
		})
	}

	extraParams := s.Config.RcloneParams.MoveServerSide
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	if globalParams := rclone.GetGlobalParams(rclone.GlobalMoveServerSideParams, s.Config.RcloneParams.GlobalMoveServerSide); globalParams != nil {
		extraParams = append(extraParams, globalParams...)
	}

	// iterate remotes and run move
	for _, move := range moveRemotes {
		// set variables
		attempts := 1
		rLog := s.Log.WithFields(logrus.Fields{
			"move_to":   move.To,
			"move_from": move.From,
			"attempts":  attempts,
		})

		// move to remote
		rLog.Info("Moving...")
		success, exitCode, err := rclone.Move(move.From, move.To, nil, true, extraParams)

		// check result
		if err != nil {
			rLog.WithError(err).Errorf("Failed unexpectedly...")
			return errors.WithMessagef(err, "move failed unexpectedly with exit code: %v", exitCode)
		} else if success {
			// successful exit code
			continue
		}

		return fmt.Errorf("move failed with exit code: %v", exitCode)
	}

	return nil
}
