package config

type RcloneConfig struct {
	Path                  string                  `koanf:"path"`
	Config                string                  `koanf:"config"`
	Stats                 string                  `koanf:"stats"`
	LiveRotate            bool                    `koanf:"live_rotate"`
	DryRun                bool                    `koanf:"dry_run"`
	ServiceAccountRemotes map[string][]string     `koanf:"service_account_remotes"`
	GlobalParams          map[string]RcloneParams `koanf:"global_params"`
}

type RcloneServerSide struct {
	From string
	To   string
}

type RcloneParams struct {
	Copy           []string
	Move           []string
	MoveServerSide []string `koanf:"move_server_side"`
	Sync           []string
	Dedupe         []string
}
