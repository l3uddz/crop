package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/stringutils"
)

type Configuration struct {
	Rclone   RcloneConfig
	Uploader []UploaderConfig
	Syncer   []SyncerConfig
}

/* Vars */

var (
	cfgPath = ""

	// Config exports the config object
	Config *Configuration

	// Internal
	delimiter = "."
	k         = koanf.New(delimiter)

	log = logger.GetLogger("cfg")
)

/* Public */

func Init(configFilePath string) error {
	// set package variables
	cfgPath = configFilePath

	// load config file
	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return err
	}

	// unmarshal into struct
	if err := k.Unmarshal("", &Config); err != nil {
		return err
	}

	return nil
}

func ShowUsing() {
	log.Infof("Using %s = %q", stringutils.LeftJust("CONFIG", " ", 10), cfgPath)
}
