package cache

import "time"

func Set(key string, expires time.Time) error {
	// acquire mutex
	mtx.Lock()
	defer mtx.Unlock()

	// set cache
	vault[key] = expires

	// log
	log.Warnf("Banned %q: %v", key, expires)

	// dump cache to disk
	return dumpToFile()
}
