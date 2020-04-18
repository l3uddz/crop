package uploader

import (
	"fmt"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/reutils"
	"github.com/l3uddz/crop/stringutils"
	"github.com/l3uddz/crop/uploader/checker"
	"github.com/l3uddz/crop/uploader/cleaner"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
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

	ServiceAccountFiles []pathutils.Path
	ServiceAccountCount int

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

	// - service account files
	var serviceAccountFiles []pathutils.Path
	if uploaderConfig.ServiceAccountFolder != "" {
		serviceAccountFiles, _ = pathutils.GetPathsInFolder(uploaderConfig.ServiceAccountFolder, true,
			false, func(path string) *string {
				lowerPath := strings.ToLower(path)

				// ignore non json files
				if !strings.HasSuffix(lowerPath, ".json") {
					return nil
				}

				return &path
			})

		// sort service files
		if len(serviceAccountFiles) > 0 {
			re := regexp.MustCompile("[0-9]+")
			sort.SliceStable(serviceAccountFiles, func(i, j int) bool {
				is := stringutils.NewOrExisting(re.FindString(serviceAccountFiles[i].RealPath), "0")
				js := stringutils.NewOrExisting(re.FindString(serviceAccountFiles[j].RealPath), "0")

				in, err := strconv.Atoi(is)
				if err != nil {
					return false
				}
				jn, err := strconv.Atoi(js)
				if err != nil {
					return false
				}

				return in < jn
			})
		}
	}

	// init uploader
	uploader := &Uploader{
		Log:                 logger.GetLogger(uploaderName),
		GlobalConfig:        config,
		Config:              uploaderConfig,
		Name:                uploaderName,
		Checker:             chk,
		Cleaner:             cln,
		IncludePatterns:     includePatterns,
		ExcludePatterns:     excludePatterns,
		ServiceAccountFiles: serviceAccountFiles,
		ServiceAccountCount: len(serviceAccountFiles),
	}

	return uploader, nil
}
