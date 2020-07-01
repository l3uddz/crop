package config

import (
	"fmt"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/stringutils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Rclone   RcloneConfig
	Uploader []UploaderConfig
	Syncer   []SyncerConfig
}

/* Vars */

var (
	Config *Configuration

	// internal
	cfgPath = ""
	log     = logger.GetLogger("cfg")
)

/* Public */

func Init(configFilePath string) error {
	// set package variables
	cfgPath = configFilePath

	// read config file
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed reading config file: %w", err)
	}

	// decode config file
	if err := yaml.Unmarshal(b, &Config); err != nil {
		return fmt.Errorf("failed decoding config file: %w", err)
	}

	return nil
}

func ShowUsing() {
	log.Infof("Using %s = %q", stringutils.LeftJust("CONFIG", " ", 10), cfgPath)
}
