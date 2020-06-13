package pathutils

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"runtime"
)

/* Public */

func GetCurrentBinaryPath() string {
	// get current binary path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// get current working dir
		if dir, err = os.Getwd(); err != nil {
			panic("failed to determine current binary location")
		}
	}

	return dir
}

func GetDefaultConfigPath() string {
	// get binary path
	bp := GetCurrentBinaryPath()
	if dirIsWriteable(bp) == nil {
		return bp
	}

	// binary path is not write-able, use alternative path
	uhp, err := os.UserHomeDir()
	if err != nil {
		panic("failed to determine current user home directory")
	}

	// set crop path inside user home dir
	chp := filepath.Join(uhp, ".config", "crop")
	if _, err := os.Stat(chp); os.IsNotExist(err) {
		if e := os.MkdirAll(chp, os.ModePerm); e != nil {
			panic("failed to create crop config directory")
		}
	}

	return chp
}

/* Private */

func dirIsWriteable(dir string) error {
	// credits: https://stackoverflow.com/questions/20026320/how-to-tell-if-folder-exists-and-is-writable
	var err error

	if runtime.GOOS != "windows" {
		err = unix.Access(dir, unix.W_OK)
	} else {
		f, e := os.Stat(dir)
		if e != nil {
			return e
		}

		switch {
		case !f.IsDir():
			err = errors.New("dir is not a directory")
		case f.Mode().Perm()&(1<<(uint(7))) == 0:
			err = errors.New("dir is not writeable")
		default:
			break
		}
	}

	return err
}
