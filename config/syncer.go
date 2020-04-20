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
	Enabled              bool
	ServiceAccountFolder string `mapstructure:"sa_folder"`
	SourceRemote         string `mapstructure:"source_remote"`
	Remotes              SyncerRemotes
	RcloneParams         SyncerRcloneParams `mapstructure:"rclone_params"`
}
