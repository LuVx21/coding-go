package keeplive

import (
	"context"
	"crypto/tls"
	"fmt"
	"luvx/gin/config"
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/redis/go-redis/v9"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		{Name: "保活", Crontab: "23 23 4/12 * * *", Fn: func() { common_x.RunCatching(aaa) }},
	}
}

func aaa() {
	for _, config := range config.Viper.GetStringMap("keeplive") {
		m := config.(map[string]any)
		_type := m["type"]
		if _type == "redis" {
			klRedisCli := redis.NewClient(&redis.Options{
				Addr:     m["host"].(string),
				Username: m["username"].(string),
				Password: m["password"].(string),
				DB:       0,
				TLSConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			})
			s, e := klRedisCli.Get(context.Background(), "foo").Result()
			fmt.Println(s, e == nil)
		}
	}
}
