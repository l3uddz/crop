package checker

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
)

type Size struct{}

func (Size) Check(cfg *config.UploaderCheck, log *logrus.Entry, paths []pathutils.Path, size uint64) (bool, error) {
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

func (Size) CheckFile(cfg *config.UploaderCheck, log *logrus.Entry, path pathutils.Path, size uint64) (bool, error) {
	// Check Total Size
	if size > cfg.Limit {
		return true, nil
	}

	return false, nil
}

func (Size) RcloneParams(cfg *config.UploaderCheck, log *logrus.Entry) []string {
	return rclone.IncludeExcludeToFilters(cfg.Include, cfg.Exclude)
}
