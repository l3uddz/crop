package rclone

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
)

var (
	log = logger.GetLogger("rclone")

	// init
	cfg *config.Configuration
)

/* Struct */

type RemoteInstruction struct {
	From       string
	To         string
	ServerSide bool
}

/* Public */

func Init(c *config.Configuration) error {
	// set required globals
	cfg = c

	// load service files for all uploader(s)
	return nil
}
