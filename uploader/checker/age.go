package checker

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
	"time"
)

type Age struct{}

func (_ Age) Check(cfg *config.UploaderCheck, log *logrus.Entry, paths []pathutils.Path, size uint64) (bool, error) {
	maxFileAge := time.Now().Add(time.Duration(-cfg.Limit) * time.Minute)

	// Check File Ages
	for _, path := range paths {
		// skip directories
		if path.IsDir {
			continue
		}

		// was this file modified after our max file age?
		if path.ModifiedTime.Before(maxFileAge) {
			log.WithFields(logrus.Fields{
				"max_age":   humanize.Time(maxFileAge),
				"file_time": path.ModifiedTime,
				"file_path": path.Path,
				"over_age":  humanize.RelTime(maxFileAge, path.ModifiedTime, "", ""),
			}).Info("Age is greater than specified limit")
			return true, nil
		}
	}

	return false, nil
}

func (_ Age) CheckFile(cfg *config.UploaderCheck, log *logrus.Entry, path pathutils.Path, size uint64) (bool, error) {
	maxFileAge := time.Now().Add(time.Duration(-cfg.Limit) * time.Minute)

	// Check File Age
	if path.ModifiedTime.Before(maxFileAge) {
		return true, nil
	}

	return false, nil
}

func (_ Age) RcloneParams(cfg *config.UploaderCheck, log *logrus.Entry) []string {
	params := []string{
		"--min-age",
		fmt.Sprintf("%dm", cfg.Limit),
	}

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
