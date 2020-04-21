package config

type RcloneConfig struct {
	Path                  string
	Config                string
	Stats                 string
	DryRun                bool              `mapstructure:"dry_run"`
	ServiceAccountRemotes map[string]string `mapstructure:"service_account_remotes"`
}

type RcloneServerSide struct {
	From string
	To   string
}
