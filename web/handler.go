package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/l3uddz/crop/cache"
	"time"
)

func (ws *Server) ServiceAccountHandler(c *fiber.Ctx) error {
	// only accept json
	c.Accepts("application/json")

	// acquire cache lock
	ws.saCache.Lock()
	defer ws.saCache.Unlock()

	// parse body
	req := new(ServiceAccountRequest)
	if err := c.BodyParser(req); err != nil {
		ws.log.WithError(err).Error("Failed parsing service account request from rclone...")
		return c.SendStatus(500)
	}

	// have we issued a replacement sa for this banned sa?
	now := time.Now().UTC()
	nsa, ok := ws.saCache.cache[req.OldServiceAccount]
	switch {
	case ok && now.Before(nsa.Expires):
		// we issued a replacement sa for this one already
		nsa.Hits++
		if nsa.Hits <= maxSaCacheHits {
			// return last response
			return c.SendString(nsa.ResponseServiceAccount)
		}

		// remove entries that have exceeded max hits
		delete(ws.saCache.cache, req.OldServiceAccount)
	case ok:
		// we issued a replacement sa for this one already, but it has expired
		delete(ws.saCache.cache, req.OldServiceAccount)
	default:
		break
	}

	// handle response
	ws.log.Warnf("Service account limit reached for remote %q, sa: %v", req.Remote, req.OldServiceAccount)

	// ban this service account
	if err := cache.SetBanned(req.OldServiceAccount, 25); err != nil {
		ws.log.WithError(err).Error("Failed banning service account, cannot try again...")
		return c.SendStatus(500)
	}

	// get service account for this remote
	sa, err := ws.sa.GetServiceAccount(req.Remote)
	switch {
	case err != nil:
		ws.log.WithError(err).Errorf("Failed retrieving service account for remote: %q", req.Remote)
		return c.SendStatus(500)
	case len(sa) < 1:
		ws.log.Errorf("Failed finding service account for remote: %q", req.Remote)
		return c.SendStatus(500)
	default:
		break
	}

	// create cache entry
	cacheEntry := &ServiceAccountCacheEntry{
		ResponseServiceAccount: sa[0].ServiceAccountPath,
		Expires:                time.Now().UTC().Add(durationSaCacheEntry),
		Hits:                   0,
	}

	// store cache entry for the old account
	ws.saCache.cache[req.OldServiceAccount] = cacheEntry

	// store cache entry for the new account
	// (so if another transfer routine requests within N duration, re-issue the same sa)
	ws.saCache.cache[sa[0].ServiceAccountPath] = cacheEntry

	// return service account
	ws.log.Warnf("New service account for remote %q, sa: %v", req.Remote, sa[0].ServiceAccountPath)
	return c.SendString(sa[0].ServiceAccountPath)
}
