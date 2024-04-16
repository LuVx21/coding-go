package config

import (
    "flag"
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
    "os"
)

var AppConfig Config

func config(configName string) Config {
    //viper := viper.New()
    viper.SetConfigName(configName)
    viper.SetConfigType("yml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath("$GOPATH/config")
    viper.AddConfigPath("$HOME/OneDrive/Code/gist/Go/gists/config")
    err := viper.ReadInConfig()
    if err != nil {
        panic(err)
    }

    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    })
    viper.WatchConfig()

    var c Config
    viper.Unmarshal(&c)
    return c
}

func Exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

func init() {
    var env = *flag.String("env", "dev", "go run main.go -env dev")

    //if !flag.Parsed() {
    //  测试时候会出现问题: flag provided but not defined
    //    flag.Parse()
    //}

    configName := ""
    switch env {
    case "test":
        configName = "config-test"
    case "prd":
        configName = "config-prd"
    //case "dev":
    //    configName = "config-dev"
    default:
        configName = "config-dev"
    }

    AppConfig = config(configName)
}
