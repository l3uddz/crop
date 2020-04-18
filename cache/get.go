package cache

import (
	"time"
)

func Get(key string) bool {
	changes := false
	exists := false

	// acquire mutex
	mtx.Lock()
	defer mtx.Unlock()

	// does the item exist in cache
	if expiry, ok := vault[key]; ok {
		// has the item expired?
		if expiry.Before(time.Now().UTC()) {
			log.Warnf("Expired %q: %v", key, expiry)
			changes = true

			// remove item from cache
			delete(vault, key)
		} else {
			// the item has not expired
			exists = true
		}
	}

	// changes made?
	if changes {
		// update the cache file
		if err := dumpToFile(); err != nil {
			log.WithError(err).Error("Failed storing updated cache")
		}
	}

	return exists
}
