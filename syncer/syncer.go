package syncer

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/rclone"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Syncer struct {
	// Public
	Log                       *logrus.Entry
	GlobalConfig              *config.Configuration
	Config                    *config.SyncerConfig
	Name                      string
	RemoteServiceAccountFiles *rclone.ServiceAccountManager
	Ws                        *WebServer
}

func New(config *config.Configuration, syncerConfig *config.SyncerConfig, syncerName string, parallelism int) (*Syncer, error) {
	// init syncer dependencies
	// - service account manager
	sam := rclone.NewServiceAccountManager(config.Rclone.ServiceAccountRemotes, parallelism)

	remotePaths := append([]string{}, syncerConfig.Remotes.Copy...)
	remotePaths = append(remotePaths, syncerConfig.Remotes.Sync...)
	remotePaths = append(remotePaths, syncerConfig.SourceRemote)

	if err := sam.LoadServiceAccounts(remotePaths); err != nil {
		return nil, errors.WithMessage(err, "failed initializing associated remote service accounts")
	}

	// init syncer
	l := logger.GetLogger(syncerName)

	syncer := &Syncer{
		Log:                       l,
		GlobalConfig:              config,
		Config:                    syncerConfig,
		Name:                      syncerName,
		RemoteServiceAccountFiles: sam,
		Ws:                        newWebServer("127.0.0.1", 3263, l, sam),
	}

	return syncer, nil
}
