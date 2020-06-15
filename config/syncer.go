package config

type SyncerRemotes struct {
	Copy           []string
	Sync           []string
	MoveServerSide []RcloneServerSide `koanf:"move_server_side"`
	Dedupe         []string
}

type SyncerRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `koanf:"global_copy"`
	Sync                 []string
	GlobalSync           string   `koanf:"global_sync"`
	MoveServerSide       []string `koanf:"move_server_side"`
	GlobalMoveServerSide string   `koanf:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `koanf:"global_dedupe"`
}

type SyncerConfig struct {
	Name         string
	Enabled      bool
	SourceRemote string `koanf:"source_remote"`
	Remotes      SyncerRemotes
	RcloneParams SyncerRcloneParams `koanf:"rclone_params"`
}
