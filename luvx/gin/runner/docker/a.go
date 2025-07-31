package docker

import (
	"context"
	"fmt"
	"log"
	"luvx/gin/common/consts"
	"luvx/gin/config"
	"luvx/gin/db"
	"luvx/gin/service"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/luvx21/coding-go/coding-common/brace"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/infra/logs"
	"golang.org/x/exp/maps"
)

var (
	re          = regexp.MustCompile(`^(?:([^/]+)/)?([^:]+)(?::(.+))?$`)
	rateLimiter = consts.GetRateLimiter("http://docker.mini.rx")
)

const (
	redis_key_last_sync_time = "app:docker:last_sync_time"
	redis_key_image_time     = "app:docker"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		{Name: "docker image拉新", Crontab: "23 29 4/12 * * *", Fn: func() { common_x.RunCatching(a) }},
	}
}

func a() {
	rss := config.Viper.GetStringMap("docker.rss")
	tul := rss["rsshub"].(string) + "/dockerhub/build"
	lines := strings.Split(rss["images"].(string), ",")
	if rss["remote"].(bool) {
		lines = remote(rss["url"].(string))
	}
	lastSyncTimeStr := db.RedisClient.Get(context.Background(), redis_key_last_sync_time).Val()
	lastSyncTime, _ := time.Parse(time.DateTime, strings_x.FirstNonEmpty(lastSyncTimeStr, "2020-01-01 00:00:00"))
	values := make([]any, 0)
	images := slices_x.FlatMap(lines, func(s string) []string { return brace.Expand(s) })
	for _, image := range images {
		owner, imageName, tag := parseImage(image)
		if imageName == "" {
			continue
		}

		_url := tul + "/" + owner + "/" + imageName + "/" + tag
		_ = rateLimiter.Wait(context.TODO())
		_, body, errs := consts.GoRequest.Get(_url).End()
		if errs != nil {
			continue
		}
		pubDate := parseXml(body)
		latest, _ := time.Parse(time.RFC1123, pubDate)
		fmt.Println("docker 镜像", image, pubDate, latest)
		if latest.After(lastSyncTime) {
			values = append(values, image, latest.Format(time.DateTime))
		}
	}
	db.RedisClient.HMSet(context.Background(), redis_key_image_time, values...)
	m := db.RedisClient.HGetAll(context.Background(), redis_key_image_time).Val()
	logs.Log.Infoln("需要同步的image", strings.Join(maps.Keys(m), ","))
}

func remote(_url string) []string {
	_, body, _ := consts.GoRequest.Get(_url).End()
	return strings.Split(body, "\n")[3:]
}

func parseImage(image string) (owner string, imageName string, tag string) {
	matches := re.FindStringSubmatch(image)
	if len(matches) < 4 {
		return "", "", ""
	}
	owner = common_x.IfThen(matches[1] == "", "library", matches[1])
	imageName = matches[2]
	tag = common_x.IfThen(matches[3] == "", "latest", matches[3])
	return
}

func parseXml(xmlData string) string {
	doc, err := xmlquery.Parse(strings.NewReader(xmlData))
	if err != nil {
		log.Fatal(err)
	}
	// 使用XPath查询所有pubDate节点
	nodes := xmlquery.Find(doc, "//pubDate")
	if len(nodes) < 1 {
		return ""
	}
	return nodes[0].InnerText()
}
