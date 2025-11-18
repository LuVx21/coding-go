package configs_x

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	viper_p "github.com/spf13/viper"
)

func LoadConfig(configName string, paths ...string) *viper_p.Viper {
	viper := viper_p.New()
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$GOPATH/config")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("加载配置文件异常: %s", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	return viper
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
