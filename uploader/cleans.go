package uploader

import (
	"github.com/pkg/errors"
	"github.com/yale8848/gorpool"
	"time"
)

func (u *Uploader) PerformCleans(gp *gorpool.Pool) error {
	// refresh details about hidden files/folders to remove
	if err := u.RefreshHiddenPaths(); err != nil {
		u.Log.WithError(err).Error("Failed refreshing details of hidden files/folders to clean")
		return errors.Wrap(err, "failed refreshing details of hidden files/folders")
	}

	// perform clean files
	if len(u.HiddenFiles) > 0 {
		u.Log.Info("Performing clean of hidden files...")
		for _, path := range u.HiddenFiles {
			p := path

			gp.AddJob(func() {
				_ = u.Clean(&p)
			})
		}

		u.Log.Debug("Waiting for queued jobs to finish")
		time.Sleep(2 * time.Second)
		gp.WaitForAll()
		u.Log.Info("Finished cleaning hidden files!")
	}

	// perform clean folders
	if len(u.HiddenFolders) > 0 {
		u.Log.Info("Performing clean of hidden folders...")
		for _, path := range u.HiddenFolders {
			p := path

			gp.AddJob(func() {
				_ = u.Clean(&p)
			})
		}

		u.Log.Debug("Waiting for queued jobs to finish")
		time.Sleep(2 * time.Second)
		gp.WaitForAll()
		u.Log.Info("Finished cleaning hidden folders!")
	}

	return nil
}
