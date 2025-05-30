package weibo_p

import (
	"context"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/db"
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunnerRegister() []*service.Runner {
	result, _ := db.RedisClient.HGet(context.TODO(), consts.AppSwitchKey, "runner_weibo").Result()
	if !cast_x.ToBool(result) {
		return make([]*service.Runner, 0)
	}
	return []*service.Runner{
		// {Name: "拉取微博热搜", Crontab: "0 7/10 * * * *", Fn: func() {
		// 	common_x.RunCatching(func() {
		// 		service.RunnerLocker.LockRun("拉取微博热搜", time.Minute*10, PullHotBand)
		// 	})
		// }},
		service.NewRunner("拉取微博热搜", "0 7/10 * * * *", time.Minute*7, PullHotBand),
		{Name: "拉取分组微博", Crontab: "0 4/4 * * * *", Fn: func() { common_x.RunCatching(PullByGroupLock) }},
		{Name: "删除weibo已读", Crontab: "0 1/2 * * * *", Fn: func() { common_x.RunCatching(DeleteLock) }},
	}
}

func DeleteLock() {
	service.RunnerLocker.LockRun("删除weibo已读", time.Minute*2, Delete)
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
    and guid <= (select guid
              from freshrss.t_admin_entry
              where true
                and id_feed = ?
              order by guid desc
              limit 1)
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
	opts := options.Find().
		SetProjection(bson.D{{Key: "_id", Value: 1}, {Key: "retweeted_status", Value: 1}}).
		SetLimit(300)
	rowsMap, _ := mongodb.RowsMap(context.TODO(), collection, filter, opts)
	ids := slices_x.Transfer(func(m bson.M) int64 { return cast_x.ToInt64(m["retweeted_status"]) }, *rowsMap...)
	idsStr := slices_x.Transfer(func(m bson.M) string { return cast_x.ToString(m["retweeted_status"]) }, *rowsMap...)
	guids = append(guids, ids...)
	mysqlGuids = append(mysqlGuids, idsStr...)

	if len(guids) > 0 {
		filter = bson.D{bson.E{Key: "_id", Value: bson.M{"$in": guids}}}
		update := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "invalid", Value: 1},
				{Key: "read", Value: 1},
			},
		}}
		dr, err := collection.UpdateMany(context.TODO(), filter, update)

		// dr, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			return
		}
		logs.Log.Infoln("mongodb删除数量:", dr.ModifiedCount)
	}

	db.MySQLClient.Table("freshrss.t_admin_entry").Delete(nil, "guid in ? and is_favorite = 0", mysqlGuids)
}
