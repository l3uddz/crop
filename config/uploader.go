package config

type UploaderCheck struct {
	Forced  bool
	Type    string
	Limit   uint64
	Exclude []string
	Include []string
}

type UploaderHidden struct {
	Enabled bool
	Type    string
	Folder  string
	Cleanup bool
	Workers int
}

type UploaderRemotes struct {
	Clean          []string
	Copy           []string
	Move           string
	MoveServerSide []RcloneServerSide `koanf:"move_server_side"`
	Dedupe         []string
}

type UploaderRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `koanf:"global_copy"`
	Move                 []string
	GlobalMove           string   `koanf:"global_move"`
	MoveServerSide       []string `koanf:"move_server_side"`
	GlobalMoveServerSide string   `koanf:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `koanf:"global_dedupe"`
}

type UploaderConfig struct {
	Name         string
	Enabled      bool
	Check        UploaderCheck
	Hidden       UploaderHidden
	LocalFolder  string `koanf:"local_folder"`
	Remotes      UploaderRemotes
	RcloneParams UploaderRcloneParams `koanf:"rclone_params"`
}
