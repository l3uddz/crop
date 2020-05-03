package config

type UploaderCheck struct {
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
	MoveServerSide []RcloneServerSide `mapstructure:"move_server_side"`
	Dedupe         []string
}

type UploaderRcloneParams struct {
	Copy                 []string
	GlobalCopy           string `mapstructure:"global_copy"`
	Move                 []string
	GlobalMove           string   `mapstructure:"global_move"`
	MoveServerSide       []string `mapstructure:"move_server_side"`
	GlobalMoveServerSide string   `mapstructure:"global_move_server_side"`
	Dedupe               []string
	GlobalDedupe         string `mapstructure:"global_dedupe"`
}

type UploaderConfig struct {
	Name         string
	Enabled      bool
	Check        UploaderCheck
	Hidden       UploaderHidden
	LocalFolder  string `mapstructure:"local_folder"`
	Remotes      UploaderRemotes
	RcloneParams UploaderRcloneParams `mapstructure:"rclone_params"`
}
