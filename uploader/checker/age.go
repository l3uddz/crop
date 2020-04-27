package checker

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/sirupsen/logrus"
	"time"
)

type Age struct{}

func (Age) Check(cfg *config.UploaderCheck, log *logrus.Entry, paths []pathutils.Path, size uint64) (bool, error) {
	var checkPassed bool
	var filesPassed int
	var filesSize int64

	// Check File Ages
	maxFileAge := time.Now().Add(time.Duration(-cfg.Limit) * time.Minute)

	for _, path := range paths {
		// skip directories
		if path.IsDir {
			continue
		}

		// was this file modified after our max file age?
		if path.ModifiedTime.Before(maxFileAge) {
			filesPassed++
			filesSize += path.Size

			log.WithFields(logrus.Fields{
				"max_age":   humanize.Time(maxFileAge),
				"file_time": path.ModifiedTime,
				"file_path": path.Path,
				"over_age":  humanize.RelTime(maxFileAge, path.ModifiedTime, "", ""),
			}).Trace("Age is greater than specified limit")

			checkPassed = true
		}
	}

	if checkPassed {
		log.WithFields(logrus.Fields{
			"files_passed": filesPassed,
			"files_size":   humanize.Bytes(uint64(filesSize)),
		}).Info("Local files matching check criteria")
	}

	return checkPassed, nil
}

func (Age) CheckFile(cfg *config.UploaderCheck, log *logrus.Entry, path pathutils.Path, size uint64) (bool, error) {
	maxFileAge := time.Now().Add(time.Duration(-cfg.Limit) * time.Minute)

	// Check File Age
	if path.ModifiedTime.Before(maxFileAge) {
		return true, nil
	}

	return false, nil
}

func (Age) RcloneParams(cfg *config.UploaderCheck, log *logrus.Entry) []string {
	params := []string{
		"--min-age",
		fmt.Sprintf("%dm", cfg.Limit),
	}

	// add filters
	params = append(params, rclone.IncludeExcludeToFilters(cfg.Include, cfg.Exclude)...)

	return params
}
