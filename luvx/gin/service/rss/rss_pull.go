package rss

import (
	"context"
	"fmt"
	"luvx/gin/common/consts"
	"luvx/gin/config"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/db"
	"luvx/gin/service/webhook"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/luvx21/coding-go/coding-common/brace"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

var (
	re = regexp.MustCompile(`^(?:([^/]+)/)?([^:]+)(?::(.+))?$`)
)

const (
	redis_key_last_sync_time = "app:rss:last_sync_time"
	redis_key_time           = "app:rss"
)

func a() {
	rss := config.Viper.GetStringMap("rss")
	for _, cate := range []string{"docker", "github"} {
		m := rss[cate].(map[string]any)
		pullFrom := m[strings.ToLower("pullFrom")].(string)
		remote, keys, keyUrl := m["remote"].(bool), m["keys"].(string), m[strings.ToLower("keyUrl")].(string)
		parsePath := m[strings.ToLower("parsePath")].(string)

		keySli := common_x.IfThenGet(remote,
			func() []string { return parseRemoteUrl(keyUrl) },
			func() []string { return strings.Split(keys, ",") },
		)
		keySli = slices_x.FlatMap(keySli, func(s string) []string { return brace.Expand(s) })
		lastSyncTimeStr := db.RedisClient.Get(context.Background(), redis_key_last_sync_time+":"+cate).Val()
		lastSyncTime, _ := time.Parse(time.DateTime, strings_x.FirstNonEmpty(lastSyncTimeStr, "2020-01-01 00:00:00"))
		values := make([]any, 0)
		for _, key := range keySli {
			var args []any
			switch cate {
			case "docker":
				owner, imageName, tag := parseDockerImage(key)
				if imageName == "" {
					continue
				}
				args = append(args, owner, imageName, tag)
			case "github":
				args = append(args, key)
			default:
				continue
			}
			_url := fmt.Sprintf(pullFrom, args...)
			_ = consts.GetRateLimiter(_url).Wait(context.TODO())
			r, body, errs := consts.GoRequest.Get(_url).End()
			if errs != nil || r.StatusCode/100 != 2 {
				continue
			}
			pubDate := parseXml(body, parsePath)
			latest, _ := times_x.StringToDate(pubDate)
			if latest.After(lastSyncTime) {
				values = append(values, key, latest.Format(time.DateTime))
			}
			log.Infof("%-40s %s", key, latest.Format(time.DateTime))
		}
		db.RedisClient.HMSet(context.Background(), redis_key_time+":"+cate, values...)
		mm := db.RedisClient.HGetAll(context.Background(), redis_key_time+":"+cate).Val()
		if len(mm) > 0 {
			log.Infoln("有新版的组件", strings.Join(maps.Keys(mm), ","))
			webhook.SendMessage(webhook.TO_USER, mongo_dao.DynamicConfig.Get()["template_id_2"].(string),
				"", map[string]webhook.WeixinMsg{
					"cate":    {Value: cate},
					"content": {Value: strings.Join(maps.Keys(mm), "\n")},
				},
			)
		}
	}
}

func parseRemoteUrl(_url string) []string {
	_, body, _ := consts.GoRequest.Get(_url).End()
	return strings.Split(body, "\n")[3:]
}

func parseDockerImage(image string) (owner string, imageName string, tag string) {
	matches := re.FindStringSubmatch(image)
	if len(matches) < 4 {
		return "", "", ""
	}
	owner = common_x.IfThen(matches[1] == "", "library", matches[1])
	imageName = matches[2]
	tag = common_x.IfThen(matches[3] == "", "latest", matches[3])
	return
}

func parseXml(xmlData, xmlPath string) string {
	doc, err := xmlquery.Parse(strings.NewReader(xmlData))
	if err != nil {
		log.Fatal(err)
	}
	// 使用XPath查询所有pubDate节点
	nodes := xmlquery.Find(doc, xmlPath)
	if len(nodes) < 1 {
		return ""
	}
	return nodes[0].InnerText()
}
