package config

type RcloneConfig struct {
	Path   string
	Config string
	DryRun bool `mapstructure:"dry_run"`
}

type RcloneServerSide struct {
	From string
	To   string
}
