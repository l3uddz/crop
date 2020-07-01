package config

type RcloneConfig struct {
	Path                  string                  `yaml:"path"`
	Config                string                  `yaml:"config"`
	Stats                 string                  `yaml:"stats"`
	LiveRotate            bool                    `yaml:"live_rotate"`
	DryRun                bool                    `yaml:"dry_run"`
	ServiceAccountRemotes map[string][]string     `yaml:"service_account_remotes"`
	GlobalParams          map[string]RcloneParams `yaml:"global_params"`
}

type RcloneServerSide struct {
	From string
	To   string
}

type RcloneParams struct {
	Copy           []string
	Move           []string
	MoveServerSide []string `yaml:"move_server_side"`
	Sync           []string
	Dedupe         []string
}
