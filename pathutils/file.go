package pathutils

import (
	"os"
	"path/filepath"
)

/* Public */

func GetCurrentBinaryPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// get current working dir
		if dir, err = os.Getwd(); err != nil {
			os.Exit(1)
		}
	}
	return dir
}
