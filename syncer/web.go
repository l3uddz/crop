package syncer

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/recover"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
)

type WebServer struct {
	Host string
	Port int

	app *fiber.App
	log *logrus.Entry
	sa  *rclone.ServiceAccountManager
}

type ServiceAccountRequest struct {
	OldServiceAccount string `json:"old"`
	Remote            string `json:"remote"`
}

func newWebServer(host string, port int, log *logrus.Entry, sa *rclone.ServiceAccountManager) *WebServer {
	ws := &WebServer{
		Host: host,
		Port: port,
		app:  fiber.New(),
		log:  log,
		sa:   sa,
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
