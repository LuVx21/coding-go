package config

import "strings"

// GetSwitch 本地配置文件的开关
func GetSwitch(keys []string) bool {
	if len(keys) == 0 {
		return false
	}
	return Viper.GetBool("switch." + strings.Join(keys, "."))
}
