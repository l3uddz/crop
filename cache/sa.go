package cache

import (
	"github.com/zippoxer/bow"
	"time"
)

type Banned struct {
	Path    string `bow:"key"`
	Expires time.Time
}

func IsBanned(key string) (bool, time.Time) {
	// check if key was found in banned bucket
	var item Banned
	err := db.Bucket("banned").Get(key, &item)

	// was key not found
	if err == bow.ErrNotFound {
		// this key is not banned
		return false, time.Time{}
	} else if err != nil {
		log.WithError(err).Errorf("Failed checking banned bucket for: %q", key)
		return false, time.Time{}
	}

	// check if the ban has expired
	if item.Expires.Before(time.Now().UTC()) {
		// the ban has expired, remove
		log.Warnf("Expired %q: %v", key, item.Expires)

		err := db.Bucket("banned").Delete(key)
		if err != nil {
			log.WithError(err).Errorf("Failed removing from banned bucket: %q", key)
			return false, time.Time{}
		}

		return false, time.Time{}
	}

	// this key is still banned
	return true, item.Expires
}

func SetBanned(key string, hours int) error {
	expiry := time.Now().UTC().Add(time.Duration(hours) * time.Hour)

	return db.Bucket("banned").Put(Banned{
		Path:    key,
		Expires: expiry,
	})
}
