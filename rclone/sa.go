package rclone

import (
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"time"
)

/* Private */
func GetAvailableServiceAccount(serviceAccountFiles []pathutils.Path) (*pathutils.Path, error) {
	var serviceAccountFile *pathutils.Path = nil

	// find an unbanned service account
	for _, sa := range serviceAccountFiles {
		// does the cache already contain this service account?
		if exists, _ := cache.Get(sa.RealPath); exists {
			// service account is currently banned
			continue
		}

		// this service account does not exist in the cache (therefore it is unbanned)
		serviceAccountFile = &sa
		break
	}

	// was an unbanned service account found?
	if serviceAccountFile == nil {
		return nil, errors.New("no unbanned service account was available")
	}

	return serviceAccountFile, nil
}

/* Public */
func AnyRemotesBanned(remotes []string) (bool, time.Time) {
	var banned bool
	var expires time.Time

	// ignore empty remotes slice
	if remotes == nil {
		return banned, expires
	}

	// format remotes into remote names if possible
	var checkRemotes []string
	for _, remote := range remotes {
		checkRemotes = append(checkRemotes, stringutils.FromLeftUntil(remote, ":"))
	}

	// iterate remotes
	for _, remote := range checkRemotes {
		banned, expires = cache.Get(remote)
		if banned {
			break
		}
	}

	return banned, expires
}
