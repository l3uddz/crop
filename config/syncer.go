package config

type SyncerRemotes struct {
	Copy           []string
	Sync           []string
	MoveServerSide []UploaderRemotesMoveServerSide `mapstructure:"move_server_side"`
}

type SyncerRcloneParams struct {
	Copy           []string
	Sync           []string
	MoveServerSide []string `mapstructure:"move_server_side"`
}

type SyncerConfig struct {
	Enabled              bool
	ServiceAccountFolder string `mapstructure:"sa_folder"`
	SourceRemote         string `mapstructure:"source_remote"`
	Remotes              SyncerRemotes
	RcloneParams         SyncerRcloneParams `mapstructure:"rclone_params"`
}
