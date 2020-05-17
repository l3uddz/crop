package web

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/recover"
	"github.com/l3uddz/crop/rclone"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

/* Const */

const (
	maxSaCacheHits       int           = 4
	durationSaCacheEntry time.Duration = 10 * time.Second
)

/* Var */

var (
	fpc *FreePortCache
)

/* Private */

func init() {
	fpc = &FreePortCache{
		pCache: make(map[int]int),
		Mutex:  sync.Mutex{},
	}
}

/* Public */

func New(host string, log *logrus.Entry, name string, sa *rclone.ServiceAccountManager) *Server {
	// get free port
	fpc.Lock()
	defer fpc.Unlock()
	port := 0

	for {
		p, err := freeport.GetFreePort()
		if err != nil {
			log.WithError(err).Fatal("Failed locating free port for the service account server")
		}

		if _, exists := fpc.pCache[p]; !exists {
			fpc.pCache[p] = p
			port = p
			log.Debugf("Found free port for service account server: %d", port)
			break
		}
	}

	// create ws object
	ws := &Server{
		Host: host,
		Port: port,
		app:  fiber.New(),
		log:  log,
		name: name,
		sa:   sa,
		saCache: &ServiceAccountCache{
			cache: make(map[string]*ServiceAccountCacheEntry),
			Mutex: sync.Mutex{},
		},
	}

	// setup app
	ws.app.Settings.DisableStartupMessage = true

	// middleware(s)
	ws.app.Use(recover.New())

	// route(s)
	ws.app.Post("*", ws.ServiceAccountHandler)

	return ws
}

func (ws *Server) Run() {
	go func() {
		ws.log.Infof("Starting service account server: %s:%d", ws.Host, ws.Port)
		ws.Running = true

		if err := ws.app.Listen(fmt.Sprintf("%s:%d", ws.Host, ws.Port)); err != nil {
			ws.log.WithError(err).Error("Service account server failed...")
		}

		ws.Running = false
	}()
}

func (ws *Server) Stop() {
	if err := ws.app.Shutdown(); err != nil {
		ws.log.WithError(err).Error("Failed shutting down service account server...")
	}
}
