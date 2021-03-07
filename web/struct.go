package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ServiceAccountCache struct {
	cache map[string]*ServiceAccountCacheEntry
	sync.Mutex
}

type ServiceAccountCacheEntry struct {
	ResponseServiceAccount string
	Expires                time.Time
	Hits                   int
}

type Server struct {
	Host    string
	Port    int
	Running bool
	app     *fiber.App
	log     *logrus.Entry
	name    string
	sa      *rclone.ServiceAccountManager
	saCache *ServiceAccountCache
}

type FreePortCache struct {
	pCache map[int]int
	sync.Mutex
}

type ServiceAccountRequest struct {
	OldServiceAccount string `json:"old"`
	Remote            string `json:"remote"`
}
