package pathutils

import (
	"github.com/l3uddz/crop/logger"

	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	log = logger.GetLogger("paths")
)

type Path struct {
	Path             string
	RealPath         string
	RelativeRealPath string
	FileName         string
	Directory        string
	IsDir            bool
	Size             int64
	ModifiedTime     time.Time
}

type callbackAllowed func(string) *string

func GetPathsInFolder(folder string, includeFiles bool, includeFolders bool, acceptFn callbackAllowed) ([]Path,
	uint64) {
	var paths []Path
	var size uint64 = 0

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.WithError(err).Error("Failed finding paths within folder")
		return paths, size
	}

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		// handle err
		if err != nil {
			return err
		}

		// skip files if not wanted
		if !includeFiles && !info.IsDir() {
			log.Tracef("Skipping file: %s", path)
			return nil
		}

		// skip folders if not wanted
		if !includeFolders && info.IsDir() {
			log.Tracef("Skipping folder: %s", path)
			return nil
		}

		// skip paths rejected by accept callback
		realPath := path
		finalPath := path
		relativeRealPath := strings.Replace(realPath, folder, "", 1)

		if strings.HasPrefix(relativeRealPath, "/") {
			relativeRealPath = strings.Replace(relativeRealPath, "/", "", 1)
		}

		if acceptFn != nil {
			acceptedPath := acceptFn(path)
			if acceptedPath == nil {
				log.Tracef("Skipping rejected path: %s", path)
				return nil
			}

			finalPath = *acceptedPath
		}

		foundPath := Path{
			Path:             finalPath,
			RealPath:         realPath,
			RelativeRealPath: relativeRealPath,
			FileName:         info.Name(),
			Directory:        filepath.Dir(path),
			IsDir:            info.IsDir(),
			Size:             info.Size(),
			ModifiedTime:     info.ModTime(),
		}

		paths = append(paths, foundPath)
		size += uint64(info.Size())

		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve paths from: %s", folder)
	}

	return paths, size
}
