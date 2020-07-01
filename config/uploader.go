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
	MoveServerSide []RcloneServerSide `yaml:"move_server_side"`
	Dedupe         []string
}

type UploaderRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `yaml:"global_copy"`
	Move                 []string
	GlobalMove           string   `yaml:"global_move"`
	MoveServerSide       []string `yaml:"move_server_side"`
	GlobalMoveServerSide string   `yaml:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `yaml:"global_dedupe"`
}

type UploaderConfig struct {
	Name         string
	Enabled      bool
	Check        UploaderCheck
	Hidden       UploaderHidden
	LocalFolder  string `yaml:"local_folder"`
	Remotes      UploaderRemotes
	RcloneParams UploaderRcloneParams `yaml:"rclone_params"`
}
