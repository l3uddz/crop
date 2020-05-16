package rclone

import "strings"

// credits: https://github.com/rclone/rclone/blob/master/fs/config.go
func ConfigToEnv(section, name string) string {
	return "RCLONE_CONFIG_" + strings.ToUpper(strings.Replace(section+"_"+name, "-", "_", -1))
}
