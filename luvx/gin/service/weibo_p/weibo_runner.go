package weibo_p

import (
    "context"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/logs"
    "go.mongodb.org/mongo-driver/bson"
    "luvx/gin/db"
    "luvx/gin/service"
)

func RunnerRegister() []*service.Runner {
    return []*service.Runner{
        {Name: "拉取微博热搜", Crontab: "0 7/10 * * * *", Fn: func() { common_x.RunCatching(PullHotBand) }},
        {Name: "拉取分组微博", Crontab: "0 9/4 * * * *", Fn: func() { common_x.RunCatching(PullByGroup) }},
        {Name: "删除已读", Crontab: "0 1/3 * * * *", Fn: func() { common_x.RunCatching(Delete) }},
    }
}

func Delete() {
    var feeds []map[string]any
    db.MySQLClient.Table("freshrss.t_admin_feed").
        Select("id").
        Find(&feeds, "url like '%/weibo/rss/%'")

    sql := `
 select guid
 from freshrss.t_admin_entry
 where true
    and id <= (select id
              from freshrss.t_admin_entry
              where true
                and id_feed = ?
              order by guid desc
              limit 1,1)
   and id_feed = ?
   and is_read = 1
   and is_favorite = 0
 order by guid
 limit 100
`
    mysqlGuids, guids := make([]string, 0), make([]int64, 0)
    for _, rss := range feeds {
        rows, _ := db.MySQLClient.Raw(sql, rss["id"], rss["id"]).Rows()
        for rows.Next() {
            var guid string
            _ = rows.Scan(&guid)
            mysqlGuids = append(mysqlGuids, guid)
            guids = append(guids, cast_x.ToInt64(guid))
        }
    }
    if len(guids) == 0 {
        return
    }

    filter := bson.D{bson.E{Key: "_id", Value: bson.M{"$in": guids}}}
    update := bson.D{{"$set",
        bson.D{{"invalid", 1}},
    }}
    dr, err := collection.UpdateMany(context.TODO(), filter, update)

    //dr, err := collection.DeleteMany(context.TODO(), filter)
    if err != nil {
        return
    }
    logs.Log.Infoln("mongodb删除数量:", dr.ModifiedCount)
    db.MySQLClient.Table("freshrss.t_admin_entry").Delete(nil, "guid in ?", mysqlGuids)
}
