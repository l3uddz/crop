package checker

import (
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
)

type Size struct{}

func (Size) Check(cfg *config.UploaderCheck, log *logrus.Entry, paths []pathutils.Path, size uint64) (*Result, error) {
	// Check Total Size
	s := humanize.Bytes(size)

	if size > cfg.Limit {
		log.WithFields(logrus.Fields{
			"max_size":     humanize.Bytes(cfg.Limit),
			"current_size": s,
			"over_size":    humanize.Bytes(size - cfg.Limit),
		}).Info("Size is greater than specified limit")
		return &Result{
			Passed: true,
			Info:   s,
		}, nil
	} else {
		return &Result{
			Passed: false,
			Info:   humanize.Bytes(cfg.Limit - size),
		}, nil
	}
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
