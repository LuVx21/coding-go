package configs_x

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/luvx21/coding-go/coding-common/os_x"
	viper_p "github.com/spf13/viper"
)

func LoadConfig(configName string, paths ...string) *viper_p.Viper {
	viper := viper_p.New()
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	for _, path := range configPath(paths...) {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("加载配置文件异常", "Error", err)
		return nil
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

func configPath(paths ...string) []string {
	// 当前目录
	r := []string{".", "./config"}
	dir, b := os_x.Command("sh", "-c", "go list -m -f {{.Dir}}")
	if b && dir != "" {
		dir := strings.TrimSpace(dir)
		// 项目根目录下
		r = append(r, dir, filepath.Join(dir, "config"))
	}
	// 用户主目录下
	r = append(r, "$HOME/.config", "$GOPATH/config")
	// 自定义目录下
	for _, path := range paths {
		if !Exists(os.ExpandEnv(path)) {
			continue
		}
		r = append(r, path)
	}
	return r
}
