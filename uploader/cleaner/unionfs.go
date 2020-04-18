package cleaner

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
	"strings"
)

type Unionfs struct{}

func (_ Unionfs) FindHidden(cfg *config.UploaderHidden, log *logrus.Entry) ([]pathutils.Path, []pathutils.Path, error) {
	tLog := log.WithField("cleaner", "unionfs")

	// retrieve files
	files, _ := pathutils.GetPathsInFolder(cfg.Folder, true, true,
		func(path string) *string {
			if strings.HasSuffix(path, "_HIDDEN~") {
				// we are interested in hidden files/folders
				newPath := strings.ReplaceAll(path, "_HIDDEN~", "")
				return &newPath
			}

			// we are not interested in non-hidden files/folders
			return nil
		})

	// create hidden variables
	hiddenFiles := make([]pathutils.Path, 0)
	hiddenFolders := make([]pathutils.Path, 0)
	for _, path := range files {
		if !path.IsDir {
			// this is a hidden file
			hiddenFiles = append(hiddenFiles, path)
		} else {
			// this is a hidden folder
			hiddenFolders = append(hiddenFolders, path)
		}
	}

	// sort results

	// log results
	tLog.WithFields(logrus.Fields{
		"found_files":   len(hiddenFiles),
		"found_folders": len(hiddenFolders),
		"hidden_folder": cfg.Folder,
	}).Info("Refreshed hidden files/folders")
	return hiddenFiles, hiddenFolders, nil
}
