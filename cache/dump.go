package cache

import (
	"github.com/pkg/errors"
	"io/ioutil"
)

func dumpToFile() error {
	// marshal vault
	jsonData, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache data")
	}

	// write json data to file
	if err := ioutil.WriteFile(cacheFilePath, jsonData, 0644); err != nil {
		return errors.Wrapf(err, "failed to write marshalled cache data to: %q", cacheFilePath)
	}

	return nil
}
