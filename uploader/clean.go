package uploader

import (
	"github.com/l3uddz/crop/pathutils"

	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (u *Uploader) Clean(path *pathutils.Path) error {
	// iterate all remotes and remove the file/folder
	for _, remotePath := range u.Config.Remotes.Clean {
		// transform remotePath to a path that can be removed
		cleanRemotePath := strings.Replace(path.Path, u.Config.Hidden.Folder, remotePath, 1)

		// set log
		rLog := u.Log.WithFields(logrus.Fields{
			"clean_remote":      remotePath,
			"clean_local_path":  path.RealPath,
			"clean_remote_path": cleanRemotePath,
		})

		// remove from remote
		var success bool
		var exitCode int
		var err error

		rLog.Debug("Removing...")
		if path.IsDir {
			// remove directory
			success, exitCode, err = rclone.RmDir(cleanRemotePath)
		} else {
			// remove file
			success, exitCode, err = rclone.DeleteFile(cleanRemotePath)
		}

		// handle response
		if err != nil {
			// error removing
			rLog.WithError(err).WithField("exit_code", exitCode).Error("Error removing remotely")
		} else if !success {
			// failed
			rLog.WithField("exit_code", exitCode).Debug("Failed removing remotely")
		} else {
			// cleaned
			rLog.Info("Removed remotely")
		}
	}

	// cleanup cleaned path locally
	if !u.GlobalConfig.Rclone.DryRun && u.Config.Hidden.Cleanup {
		if err := os.Remove(path.RealPath); err != nil {
			u.Log.
				WithField("clean_local_path", path.RealPath).
				WithError(err).
				Error("Failed removing locally")
		} else {
			u.Log.
				WithField("clean_local_path", path.RealPath).
				Debug("Removed locally")
		}
	}
	return nil
}
