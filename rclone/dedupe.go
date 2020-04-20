package rclone

import (
	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Public */

func Dedupe(remotePath string, additionalRcloneParams []string) (bool, int, error) {
	// set variables
	rLog := log.WithFields(logrus.Fields{
		"action":      CmdDedupe,
		"remote_path": remotePath,
	})
	result := false

	// generate required rclone parameters
	params := []string{
		CmdDedupe,
		remotePath,
	}

	if baseParams, err := getBaseParams(); err != nil {
		return false, 1, errors.Wrapf(err, "failed generating baseParams to %q: %q", CmdDedupe,
			remotePath)
	} else {
		params = append(params, baseParams...)
	}

	if additionalParams, err := getAdditionalParams(CmdDedupe, additionalRcloneParams); err != nil {
		return false, 1, errors.Wrapf(err, "failed generating additionalParams to %s: %q",
			CmdDedupe, remotePath)
	} else {
		params = append(params, additionalParams...)
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
	case ExitSuccess:
		result = true
	default:
		break
	}

	rLog.WithField("exit_code", status.Exit).Debug("Finished")
	return result, status.Exit, status.Error
}
