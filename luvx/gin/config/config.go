package config

import (
	"flag"

	"github.com/luvx21/coding-go/coding-common/configs_x"
	"github.com/luvx21/coding-go/infra/logs"
)

var AppConfig Config

func init() {
	var env = *flag.String("env", "dev", "go run main.go -env dev")

	// if !flag.Parsed() {
	//  测试时候会出现问题: flag provided but not defined
	//    flag.Parse()
	// }

	configName := ""
	switch env {
	case "test", "prd":
		configName = "config-" + env
	default:
		configName = "config-dev"
	}

	logs.Log.Infoln("加载配置文件...", configName)
	viper := configs_x.LoadConfig(configName, "$HOME/OneDrive/Code/coding-go/luvx/config")
	viper.Unmarshal(&AppConfig)
}
