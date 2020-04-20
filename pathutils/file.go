package pathutils

import (
	"github.com/l3uddz/crop/logger"
	"os"
	"path/filepath"
)

var log = logger.GetLogger("paths")

/* Public */

func GetCurrentBinaryPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// get current working dir
		if dir, err = os.Getwd(); err != nil {
			// TODO: better handling here, this should never occur but still..
			os.Exit(1)
		}
	}
	return dir
}
