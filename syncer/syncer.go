package syncer

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/stringutils"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Syncer struct {
	// Public
	Log          *logrus.Entry
	GlobalConfig *config.Configuration
	Config       *config.SyncerConfig
	Name         string

	ServiceAccountFiles []pathutils.Path
	ServiceAccountCount int
}

func New(config *config.Configuration, syncerConfig *config.SyncerConfig, syncerName string) (*Syncer, error) {
	// init syncer dependencies
	// - service account files
	var serviceAccountFiles []pathutils.Path
	if syncerConfig.ServiceAccountFolder != "" {
		serviceAccountFiles, _ = pathutils.GetPathsInFolder(syncerConfig.ServiceAccountFolder, true,
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
	syncer := &Syncer{
		Log:                 logger.GetLogger(syncerName),
		GlobalConfig:        config,
		Config:              syncerConfig,
		Name:                syncerName,
		ServiceAccountFiles: serviceAccountFiles,
		ServiceAccountCount: len(serviceAccountFiles),
	}

	return syncer, nil
}
