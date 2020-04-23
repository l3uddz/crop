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

func (_ Size) RcloneParams(cfg *config.UploaderCheck, log *logrus.Entry) []string {
	params := make([]string, 0)

	// add filters
	for _, include := range cfg.Include {
		params = append(params,
			"--include", include)
	}

	for _, exclude := range cfg.Exclude {
		params = append(params,
			"--exclude", exclude)
	}

	return params
}
