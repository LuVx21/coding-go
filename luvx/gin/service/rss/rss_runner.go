package rss

import (
    "github.com/luvx21/coding-go/coding-common/common_x"
    "luvx/gin/db"
    "luvx/gin/service"
)

func RunnerRegister() []*service.Runner {
    return []*service.Runner{
        {Name: "重置rss", Crontab: "0 3/5 * * * *", Fn: func() { common_x.RunCatching(reset) }},
        {Name: "rss_spider", Crontab: "0 7 0/2 * * *", Fn: func() { common_x.RunCatching(PullByKey) }},
    }
}

func reset() {
    var feeds []map[string]any
    db.MySQLClient.Table("freshrss.t_admin_feed").
        Select("id").
        Find(&feeds, "url like '%/weibo/rss/%'")
    for _, rss := range feeds {
        db.MySQLClient.Exec("update freshrss.t_admin_feed set lastUpdate = lastUpdate-30*60 where id = ?", rss["id"])
    }
}
