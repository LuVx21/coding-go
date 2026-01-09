package config

// GetSwitch 本地配置文件的开关
func GetSwitch(key string) bool {
	if key == "" {
		return false
	}
	return Viper.GetBool("switch." + key)
}
