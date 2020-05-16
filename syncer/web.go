package syncer

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/recover"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/rclone"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"sync"
)

type WebServer struct {
	Host       string
	Port       int
	app        *fiber.App
	log        *logrus.Entry
	syncerName string
	sa         *rclone.ServiceAccountManager
}

type FreePortCache struct {
	pCache map[int]int
	sync.Mutex
}

type ServiceAccountRequest struct {
	OldServiceAccount string `json:"old"`
	Remote            string `json:"remote"`
}

var (
	fpc *FreePortCache
)

func init() {
	fpc = &FreePortCache{
		pCache: make(map[int]int),
		Mutex:  sync.Mutex{},
	}
}

func newWebServer(host string, log *logrus.Entry, syncerName string, sa *rclone.ServiceAccountManager) *WebServer {
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
	ws := &WebServer{
		Host:       host,
		Port:       port,
		app:        fiber.New(),
		log:        log,
		syncerName: syncerName,
		sa:         sa,
	}

	// setup app
	ws.app.Settings.DisableStartupMessage = true

	// middleware(s)
	ws.app.Use(recover.New())

	// route(s)
	ws.app.Post("*", ws.ServiceAccountHandler)

	return ws
}

func (ws *WebServer) Run() {
	go func() {
		ws.log.Infof("Starting service account server on %s:%d", ws.Host, ws.Port)

		if err := ws.app.Listen(fmt.Sprintf("%s:%d", ws.Host, ws.Port)); err != nil {
			ws.log.WithError(err).Error("Service account server failed...")
		}
	}()
}

func (ws *WebServer) Stop() {
	if err := ws.app.Shutdown(); err != nil {
		ws.log.WithError(err).Error("Failed shutting down service account server...")
	}
}

func (ws *WebServer) ServiceAccountHandler(c *fiber.Ctx) {
	// only accept json
	c.Accepts("application/json")

	// parse body
	req := new(ServiceAccountRequest)
	if err := c.BodyParser(req); err != nil {
		ws.log.WithError(err).Error("Failed parsing service account request from gclone...")
		c.SendStatus(500)
		return
	}

	// handle response
	ws.log.Warnf("Service account limit reached for remote %q, sa: %v", req.Remote, req.OldServiceAccount)

	// ban this service account
	if err := cache.SetBanned(req.OldServiceAccount, 25); err != nil {
		ws.log.WithError(err).Error("Failed banning service account, cannot try again...")
		c.SendStatus(500)
		return
	}

	// get service account for this remote
	sa, err := ws.sa.GetServiceAccount(req.Remote)
	switch {
	case err != nil:
		ws.log.WithError(err).Error("Failed retrieving service account for remote: %q", req.Remote)
		c.SendStatus(500)
		return
	case len(sa) < 1:
		ws.log.Error("Failed finding service account for remote: %q", req.Remote)
		c.SendStatus(500)
		return
	default:
		break
	}

	c.SendString(sa[0].ServiceAccountPath)
}
