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

type MoveInstruction struct {
	From       string
	To         string
	ServerSide bool
}

func (s *Syncer) Move(additionalRcloneParams []string) error {
	var moveRemotes []MoveInstruction

	// set variables
	for _, remote := range s.Config.Remotes.MoveServerSide {
		moveRemotes = append(moveRemotes, MoveInstruction{
			From:       remote.From,
			To:         remote.To,
			ServerSide: true,
		})
	}

	extraParams := s.Config.RcloneParams.MoveServerSide
	if additionalRcloneParams != nil {
		extraParams = append(extraParams, additionalRcloneParams...)
	}

	// iterate all remotes and run copy
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

		// ban remove
		if err := cache.Set(stringutils.FromLeftUntil(move.To, ":"),
			time.Now().UTC().Add(25*time.Hour)); err != nil {
			rLog.WithError(err).Errorf("Failed banning remote")
		}

		return fmt.Errorf("move failed with exit code: %v", exitCode)
	}

	return nil
}
