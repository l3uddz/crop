package uploader

import (
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/pathutils"
	"github.com/pkg/errors"
)

/* Private */
func (u *Uploader) getAvailableServiceAccount() (*pathutils.Path, error) {
	var serviceAccountFile *pathutils.Path = nil

	// find an unbanned service account
	for _, sa := range u.ServiceAccountFiles {
		// does the cache already contain this service account?
		if cache.Get(sa.RealPath) {
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
