package uploader

import (
	"fmt"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/reutils"
	"github.com/l3uddz/crop/uploader/checker"
	"github.com/l3uddz/crop/uploader/cleaner"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type Uploader struct {
	// Public
	Log          *logrus.Entry
	GlobalConfig *config.Configuration
	Config       *config.UploaderConfig
	Name         string

	Checker checker.Interface
	Cleaner cleaner.Interface

	IncludePatterns []*regexp.Regexp
	ExcludePatterns []*regexp.Regexp

	RemoteServiceAccountFiles *rclone.ServiceAccountManager

	LocalFiles     []pathutils.Path
	LocalFilesSize uint64
	HiddenFiles    []pathutils.Path
	HiddenFolders  []pathutils.Path
}

func New(config *config.Configuration, uploaderConfig *config.UploaderConfig, uploaderName string) (*Uploader, error) {
	// init uploader dependencies
	// - checker
	c, found := supportedCheckers[strings.ToLower(uploaderConfig.Check.Type)]
	if !found {
		return nil, fmt.Errorf("unknown check type specified: %q", uploaderConfig.Check.Type)
	}

	chk, ok := c.(checker.Interface)
	if !ok {
		return nil, fmt.Errorf("failed typecasting to checker interface for: %q", uploaderConfig.Check.Type)
	}

	// - cleaner
	var cln cleaner.Interface = nil
	if uploaderConfig.Hidden.Enabled {
		c, found := supportedCleaners[strings.ToLower(uploaderConfig.Hidden.Type)]
		if !found {
			// checker was not found
			return nil, fmt.Errorf("unknown cleaner type specified: %q", uploaderConfig.Hidden.Type)
		}

		// Typecast found cleaner
		cln, ok = c.(cleaner.Interface)
		if !ok {
			return nil, fmt.Errorf("failed typecasting to cleaner interface for: %q", uploaderConfig.Hidden.Type)
		}
	}

	// - include patterns
	var includePatterns []*regexp.Regexp

	for _, includePattern := range uploaderConfig.Check.Include {
		if g, err := reutils.GlobToRegexp(includePattern, false); err != nil {
			return nil, fmt.Errorf("invalid include pattern: %q", includePattern)
		} else {
			includePatterns = append(includePatterns, g)
		}
	}

	// - exclude patterns
	var excludePatterns []*regexp.Regexp

	for _, excludePattern := range uploaderConfig.Check.Exclude {
		if g, err := reutils.GlobToRegexp(excludePattern, false); err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %q", excludePattern)
		} else {
			excludePatterns = append(excludePatterns, g)
		}
	}

	// - service account manager
	sam := rclone.NewServiceAccountManager(config.Rclone.ServiceAccountRemotes)

	remotePaths := append([]string{}, uploaderConfig.Remotes.Copy...)
	remotePaths = append(remotePaths, uploaderConfig.Remotes.Move)

	if err := sam.LoadServiceAccounts(remotePaths); err != nil {
		return nil, errors.WithMessage(err, "failed initializing associated remote service accounts")
	}

	// init uploader
	uploader := &Uploader{
		Log:                       logger.GetLogger(uploaderName),
		GlobalConfig:              config,
		Config:                    uploaderConfig,
		Name:                      uploaderName,
		Checker:                   chk,
		Cleaner:                   cln,
		IncludePatterns:           includePatterns,
		ExcludePatterns:           excludePatterns,
		RemoteServiceAccountFiles: sam,
	}

	return uploader, nil
}
