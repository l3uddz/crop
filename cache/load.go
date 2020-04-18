package cache

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func loadFromFile(cachePath string) error {
	// does cache exist on disk?
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		// cache did not exist, so we will let it be created on demand
		return nil
	}

	// open cache on disk
	cacheFile, err := os.Open(cachePath)
	if err != nil {
		return errors.Wrapf(err, "failed loading cache file: %q", cachePath)
	}
	defer cacheFile.Close()

	// read cache data
	cacheBytes, err := ioutil.ReadAll(cacheFile)
	if err != nil {
		return errors.Wrapf(err, "failed reading bytes from cache file: %q", cachePath)
	}

	// unmarshal cache data
	if err := json.Unmarshal(cacheBytes, &vault); err != nil {
		return errors.Wrapf(err, "failed to unmarshal bytes from cache file: %q", cachePath)
	}

	return nil
}
