package rclone

import (
	"github.com/go-cmd/cmd"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Public */

func Move(u *config.UploaderConfig, localPath string, remotePath string, serviceAccountFile *pathutils.Path) (bool, int, error) {
	// set variables
	rLog := log.WithFields(logrus.Fields{
		"action":      CMD_MOVE,
		"local_path":  localPath,
		"remote_path": remotePath,
	})
	result := false

	// generate required rclone parameters
	params := []string{
		CMD_MOVE,
		localPath,
		remotePath,
	}

	if baseParams, err := getBaseParams(); err != nil {
		return false, 1, errors.Wrapf(err, "failed generating baseParams to %q: %q -> %q",
			CMD_MOVE, localPath, remotePath)
	} else {
		params = append(params, baseParams...)
	}

	if additionalParams, err := getAdditionalParams(CMD_MOVE, u.RcloneParams.Move); err != nil {
		return false, 1, errors.Wrapf(err, "failed generating additionalParams to %q: %q -> %q",
			CMD_MOVE, localPath, remotePath)
	} else {
		params = append(params, additionalParams...)
	}

	if serviceAccountFile != nil {
		params = append(params, getServiceAccountParams(serviceAccountFile)...)
	}

	rLog.Tracef("Generated params: %v", params)

	// remove file
	rcloneCmd := cmd.NewCmd(cfg.Rclone.Path, params...)
	status := <-rcloneCmd.Start()

	// check status
	switch status.Exit {
	case EXIT_SUCCESS:
		result = true
	default:
		break
	}

	rLog.WithField("exit_code", status.Exit).Debug("Finished")
	return result, status.Exit, status.Error
}
