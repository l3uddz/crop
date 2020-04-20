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
	Copy           []string
	Move           []string
	MoveServerSide []string `mapstructure:"move_server_side"`
	Dedupe         []string
}

type UploaderConfig struct {
	Enabled              bool
	Check                UploaderCheck
	Hidden               UploaderHidden
	LocalFolder          string `mapstructure:"local_folder"`
	ServiceAccountFolder string `mapstructure:"sa_folder"`
	Remotes              UploaderRemotes
	RcloneParams         UploaderRcloneParams `mapstructure:"rclone_params"`
}
