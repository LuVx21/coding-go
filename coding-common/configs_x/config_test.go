package configs_x

import (
	"fmt"
	"testing"

	"github.com/luvx21/coding-go/coding-common/jsons"
)

func Test_00(t *testing.T) {
	viper := LoadConfig("config-dev", "$HOME/OneDrive/Code/coding-go/luvx/config")
	fmt.Println(jsons.ToJsonString(viper.AllSettings()))
	fmt.Println(viper.GetString("switch.a"))

	var AppConfig config
	viper.Unmarshal(&AppConfig)
	fmt.Println(jsons.ToJsonString(AppConfig))

	subv := viper.Sub("webclient")
	fmt.Println("sub", jsons.ToJsonString(subv.AllSettings()))

	webclient := make(map[string]any)
	viper.UnmarshalKey("webclient", &webclient)
	fmt.Println("UnmarshalKey", jsons.ToJsonString(webclient))
}

func Test_01(t *testing.T) {
}

type config struct {
	Server    serverConfig
	Webclient map[string]map[string]any
}

type serverConfig struct {
	Port  string
	Debug bool
}
