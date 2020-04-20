package rclone

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
)

/* Private */

func getBaseParams() ([]string, error) {
	var params []string

	// dry run
	if cfg.Rclone.DryRun {
		params = append(params, "--dry-run")
	}

	// defaults
	params = append(params,
		// config
		"--config", cfg.Rclone.Config,
		// verbose
		"-v",
		// user-agent
		"--user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/74.0.3729.131 Safari/537.36",
	)

	// add stats
	if config.Config.Rclone.Stats != "" {
		params = append(params,
			// stats
			"--stats", config.Config.Rclone.Stats)
	}

	return params, nil
}

func getAdditionalParams(cmd string, extraParams []string) ([]string, error) {
	var params []string

	// additional params based on the rclone command being used
	switch cmd {
	case CmdCopy:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CmdMove:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CmdSync:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CmdDeleteFile:
		break
	case CmdDeleteDir:
		break
	case CmdDeleteDirs:
		break
	case CmdDedupe:
		break
	default:
		break
	}

	// add any additional params
	params = append(params, extraParams...)
	return params, nil
}

func getServiceAccountParams(serviceAccountFile *pathutils.Path) []string {
	// service account params
	params := []string{
		"--drive-service-account-file",
		serviceAccountFile.RealPath,
	}

	return params
}
