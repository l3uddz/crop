package rclone

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Public */

func Move(from string, to string, serviceAccounts []*RemoteServiceAccount, serverSide bool,
	additionalRcloneParams []string) (bool, int, error) {
	// set variables
	rLog := log.WithFields(logrus.Fields{
		"action": CmdMove,
		"from":   from,
		"to":     to,
	})
	result := false

	// generate required rclone parameters
	params := []string{
		CmdMove,
		from,
		to,
	}

	if baseParams, err := getBaseParams(); err != nil {
		return false, 1, errors.WithMessagef(err, "failed generating baseParams to %q: %q -> %q",
			CmdMove, from, to)
	} else {
		params = append(params, baseParams...)
	}

	extraParams := additionalRcloneParams
	if serverSide {
		// add server side parameter
		extraParams = append(extraParams, "--drive-server-side-across-configs")
	}

	if additionalParams, err := getAdditionalParams(CmdMove, extraParams); err != nil {
		return false, 1, errors.WithMessagef(err, "failed generating additionalParams to %q: %q -> %q",
			CmdMove, from, to)
	} else {
		params = append(params, additionalParams...)
	}

	rLog.Debugf("Generated params: %v", params)

	// generate required rclone env
	var rcloneEnv []string
	if len(serviceAccounts) > 0 {
		// iterate service accounts, creating env
		for _, env := range serviceAccounts {
			if env == nil {
				continue
			}

			v := env
			rcloneEnv = append(rcloneEnv, fmt.Sprintf("%s=%s", v.RemoteEnvVar, v.ServiceAccountPath))
		}
	}
	rLog.Debugf("Generated rclone env: %v", rcloneEnv)

	// setup cmd
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	rcloneCmd := cmd.NewCmdOptions(cmdOptions, cfg.Rclone.Path, params...)
	rcloneCmd.Env = rcloneEnv

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
