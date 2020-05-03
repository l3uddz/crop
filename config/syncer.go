package config

type SyncerRemotes struct {
	Copy           []string
	Sync           []string
	MoveServerSide []RcloneServerSide `mapstructure:"move_server_side"`
	Dedupe         []string
}

type SyncerRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `mapstructure:"global_copy"`
	Sync                 []string
	GlobalSync           string   `mapstructure:"global_sync"`
	MoveServerSide       []string `mapstructure:"move_server_side"`
	GlobalMoveServerSide string   `mapstructure:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `mapstructure:"global_dedupe"`
}

type SyncerConfig struct {
	Name         string
	Enabled      bool
	SourceRemote string `mapstructure:"source_remote"`
	Remotes      SyncerRemotes
	RcloneParams SyncerRcloneParams `mapstructure:"rclone_params"`
}
