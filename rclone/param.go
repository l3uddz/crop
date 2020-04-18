package rclone

import (
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
		// stats
		"--stats", "30s",
		// user-agent
		"--user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/74.0.3729.131 Safari/537.36",
	)

	return params, nil
}

func getAdditionalParams(cmd string, extraParams []string) ([]string, error) {
	var params []string

	// additional params based on the rclone command being used
	switch cmd {
	case CMD_COPY:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CMD_MOVE:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CMD_SYNC:
		params = append(params,
			// stop on upload limit
			"--drive-stop-on-upload-limit",
		)
	case CMD_DELETE_FILE:
		break
	case CMD_DELETE_DIR:
		break
	case CMD_DELETE_DIRS:
		break
	case CMD_DEDUPE:
		params = append(params,
			// keep newest duplicate file
			"--dedupe-mode", "newest",
		)
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
