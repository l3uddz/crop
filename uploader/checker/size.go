package checker

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
)

type Size struct{}

func (_ Size) Check(cfg *config.UploaderCheck, log *logrus.Entry, paths []pathutils.Path, size uint64) (bool, error) {

	// Check Total Size
	if size > cfg.Limit {
		log.WithFields(logrus.Fields{
			"max_size":     humanize.Bytes(cfg.Limit),
			"current_size": humanize.Bytes(size),
			"over_size":    humanize.Bytes(size - cfg.Limit),
		}).Info("Size is greater than specified limit")
		return true, nil
	}

	return false, nil
}

func (_ Size) CheckFile(cfg *config.UploaderCheck, log *logrus.Entry, path pathutils.Path, size uint64) (bool, error) {

	// Check Total Size
	if size > cfg.Limit {
		return true, nil
	}

	return false, nil
}
