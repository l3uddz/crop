package uploader

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/uploader/cleaner"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	supportedCleaners = map[string]interface{}{
		"unionfs": cleaner.Unionfs{},
	}
)

func (u *Uploader) RefreshLocalFiles() error {
	// retrieve files
	u.LocalFiles, u.LocalFilesSize = pathutils.GetPathsInFolder(u.Config.LocalFolder, true, false,
		func(path string) *string {
			rcloneStylePath := strings.TrimLeft(strings.Replace(path, u.Config.LocalFolder, "", 1), "/")

			// should this path be excluded?
			if len(u.ExcludePatterns) > 0 {
				for _, excludePattern := range u.ExcludePatterns {
					if excludePattern.MatchString(rcloneStylePath) {
						// this path matches an exclude pattern
						return nil
					}
				}
			}

			// should this path be included?
			if len(u.Config.Check.Include) > 0 {
				for _, includePattern := range u.IncludePatterns {
					if includePattern.MatchString(rcloneStylePath) {
						// this path matches an include pattern
						return &path
					}
				}

				return nil
			}

			// we are interested in all these files
			return &path
		})

	// log results
	u.Log.WithFields(logrus.Fields{
		"found_files":  len(u.LocalFiles),
		"files_size":   humanize.Bytes(u.LocalFilesSize),
		"local_folder": u.Config.LocalFolder,
	}).Info("Refreshed local files")

	return nil
}

func (u *Uploader) RefreshHiddenPaths() error {
	var err error

	// Retrieve hidden files/folders
	u.HiddenFiles, u.HiddenFolders, err = u.Cleaner.FindHidden(&u.Config.Hidden, u.Log)
	if err != nil {
		return fmt.Errorf("failed refreshing hidden paths for: %q", u.Config.Hidden.Folder)
	}

	return nil
}
