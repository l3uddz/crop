package rclone

import (
	"github.com/go-cmd/cmd"
	"github.com/l3uddz/crop/pathutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Public */

func Move(from string, to string, serviceAccountFile *pathutils.Path, serverSide bool,
	additionalRcloneParams []string) (bool, int, error) {
	// set variables
	rLog := log.WithFields(logrus.Fields{
		"action": CMD_MOVE,
		"from":   from,
		"to":     to,
	})
	result := false

	// generate required rclone parameters
	params := []string{
		CMD_MOVE,
		from,
		to,
	}

	if baseParams, err := getBaseParams(); err != nil {
		return false, 1, errors.WithMessagef(err, "failed generating baseParams to %q: %q -> %q",
			CMD_MOVE, from, to)
	} else {
		params = append(params, baseParams...)
	}

	extraParams := additionalRcloneParams
	if serverSide {
		// add server side parameter
		extraParams = append(extraParams, "--drive-server-side-across-configs")
	}

	if additionalParams, err := getAdditionalParams(CMD_MOVE, extraParams); err != nil {
		return false, 1, errors.WithMessagef(err, "failed generating additionalParams to %q: %q -> %q",
			CMD_MOVE, from, to)
	} else {
		params = append(params, additionalParams...)
	}

	if serviceAccountFile != nil && !serverSide {
		// service account file provided but this is a server side move, so ignore it
		saParams := getServiceAccountParams(serviceAccountFile)
		params = append(params, saParams...)
	}

	rLog.Debugf("Generated params: %v", params)

	// setup cmd
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	rcloneCmd := cmd.NewCmdOptions(cmdOptions, cfg.Rclone.Path, params...)

	// live stream logs
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)

		for rcloneCmd.Stdout != nil || rcloneCmd.Stderr != nil {
			select {
			case line, open := <-rcloneCmd.Stdout:
				if !open {
					rcloneCmd.Stdout = nil
					continue
				}
				log.Info(line)
			case line, open := <-rcloneCmd.Stderr:
				if !open {
					rcloneCmd.Stderr = nil
					continue
				}
				log.Info(line)
			}
		}
	}()

	// run command
	rLog.Debug("Starting...")

	status := <-rcloneCmd.Start()
	<-doneChan

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
