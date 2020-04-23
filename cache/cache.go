package cache

import (
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/stringutils"
	"github.com/pkg/errors"
	"github.com/zippoxer/bow"
	"github.com/zippoxer/bow/codec/msgp"
)

var (
	log           = logger.GetLogger("cache")
	cacheFilePath string

	// Internal
	db *bow.DB
)

/* Public */

func Init(cachePath string, logLevel int) error {
	// set globals
	cacheFilePath = cachePath

	// set badger options
	opts := []bow.Option{
		// codec
		bow.SetCodec(msgp.Codec{}),
	}

	if logLevel < 2 {
		// disable badger logging for non trace log level
		opts = append(opts, bow.SetLogger(nil))
	}

	// init database
	v, err := bow.Open(cachePath, opts...)
	if err != nil {
		return errors.WithMessage(err, "failed opening cache")
	}

	db = v

	return nil
}

func Close() {
	if err := db.Close(); err != nil {
		log.WithError(err).Error("Failed closing cache gracefully...")
	}
}

func ShowUsing() {
	log.Infof("Using %s = %q", stringutils.LeftJust("CACHE", " ", 10), cacheFilePath)
}
