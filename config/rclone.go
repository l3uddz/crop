package config

type RcloneConfig struct {
	Path                  string
	Config                string
	Stats                 string
	LiveRotate            bool                    `mapstructure:"live_rotate"`
	DryRun                bool                    `mapstructure:"dry_run"`
	ServiceAccountRemotes map[string]string       `mapstructure:"service_account_remotes"`
	GlobalParams          map[string]RcloneParams `mapstructure:"global_params"`
}

type RcloneServerSide struct {
	From string
	To   string
}

type RcloneParams struct {
	Copy           []string
	Move           []string
	MoveServerSide []string `mapstructure:"move_server_side"`
	Sync           []string
	Dedupe         []string
}
