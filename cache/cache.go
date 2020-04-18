package cache

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/stringutils"
	"sync"
	"time"
)

var (
	log           = logger.GetLogger("cache")
	cacheFilePath string
	mtx           sync.Mutex
	vault         map[string]time.Time

	// Internal
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

/* Public */

func Init(cachePath string) error {
	// set globals
	cacheFilePath = cachePath
	vault = make(map[string]time.Time)

	// load cache from disk
	if err := loadFromFile(cachePath); err != nil {
		return err
	}

	return nil
}

func ShowUsing() {
	log.Infof("Using %s = %q", stringutils.LeftJust("CACHE", " ", 10), cacheFilePath)
}
