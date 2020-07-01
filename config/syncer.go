package config

type SyncerRemotes struct {
	Copy           []string
	Sync           []string
	MoveServerSide []RcloneServerSide `yaml:"move_server_side"`
	Dedupe         []string
}

type SyncerRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `yaml:"global_copy"`
	Sync                 []string
	GlobalSync           string   `yaml:"global_sync"`
	MoveServerSide       []string `yaml:"move_server_side"`
	GlobalMoveServerSide string   `yaml:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `yaml:"global_dedupe"`
}

type SyncerConfig struct {
	Name         string
	Enabled      bool
	SourceRemote string `yaml:"source_remote"`
	Remotes      SyncerRemotes
	RcloneParams SyncerRcloneParams `yaml:"rclone_params"`
}
