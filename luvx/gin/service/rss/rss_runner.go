package rss

import (
	"time"

	"luvx/gin/dao/freshrss_dao"
	"luvx/gin/db"
	"luvx/gin/service"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		service.NewRunner("重置rss", "23 */5 * * * *", time.Minute*10, reset),
		service.NewRunner("rss_spider", "0 7 0/2 * * *", time.Minute*60, PullByKey),
		service.NewRunner("新版本时间拉新", "23 29 4/12 * * *", time.Minute*10, pullLatest),
	}
}

func reset() {
	db.FreshrssDb.Exec(`update ` + freshrss_dao.Prefix + `feed set "lastUpdate" = "lastUpdate"-30*60 where url like '%/weibo/rss/%'`)
}
