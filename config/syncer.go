package config

type SyncerRemotes struct {
	Copy           []string
	Sync           []string
	MoveServerSide []RcloneServerSide `mapstructure:"move_server_side"`
	Dedupe         []string
}

type SyncerRcloneParams struct {
	Copy           []string
	Sync           []string
	MoveServerSide []string `mapstructure:"move_server_side"`
	Dedupe         []string
}

type SyncerConfig struct {
	Name         string
	Enabled      bool
	SourceRemote string `mapstructure:"source_remote"`
	Remotes      SyncerRemotes
	RcloneParams SyncerRcloneParams `mapstructure:"rclone_params"`
}
