package configs_x

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/luvx21/coding-go/coding-common/os_x"
	viper_p "github.com/spf13/viper"
)

var (
	defaultAppName = "go_app"
	defaultViper   *viper_p.Viper

	initOnce sync.Once
)

func init() {
	const defaultConfigName = "config"

	slog.Info("初始化默认配置", "appName", defaultAppName, "配置文件名", defaultConfigName)
	initOnce.Do(func() {
		defaultViper, _ = LoadConfig(defaultConfigName)
	})
}

func RegisterAppName(n string)         { defaultAppName = n }
func GetDefaultConfig() *viper_p.Viper { return defaultViper }
func GetDefaultConfigByKey(viper *viper_p.Viper, key string) *viper_p.Viper {
	if viper == nil {
		return nil
	}
	return viper.Sub(key)
}

// LoadConfig 也可执行RegisterAppName和RegisterConfigName后, 再执行
func LoadConfig(configName string, paths ...string) (*viper_p.Viper, error) {
	viper := viper_p.New()
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	for _, path := range configDefaultPath(paths...) {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		slog.Info("加载配置文件异常", "Error", err)
		return nil, errors.New("加载配置文件异常")
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	return viper, nil
}

// configDefaultPath 默认查找配置所在的目录
func configDefaultPath(paths ...string) []string {
	// 当前目录
	r := []string{".", "./config"}
	dir, err := os.Executable()
	if err != nil && dir != "" {
		dir := strings.TrimSpace(filepath.Dir(dir))
		// 项目根目录下
		r = append(r, dir, filepath.Join(dir, "config"))
	}
	// 用户主目录下
	r = append(r, "$HOME/.config/"+defaultAppName, "$GOPATH/config")
	// 自定义目录下
	for _, path := range paths {
		if !os_x.Exists(os.ExpandEnv(path)) {
			continue
		}
		r = append(r, path)
	}
	return r
}
