package weibo_p

import (
	"context"
	"time"

	"luvx/gin/dao/redis_dao"
	"luvx/gin/db"
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func RunnerRegister() []*service.Runner {
	result := redis_dao.GetSwitch("runner_weibo")
	if !result {
		return make([]*service.Runner, 0)
	}
	return []*service.Runner{
		// {Name: "拉取微博热搜", Crontab: "0 7/10 * * * *", Fn: func() {
		// 	common_x.RunCatching(func() {
		// 		service.RunnerLocker.LockRun("拉取微博热搜", time.Minute*10, PullHotBand)
		// 	})
		// }},
		service.NewRunner("拉取微博热搜", "0 7/10 * * * *", time.Minute*7, PullHotBand),
		{Name: "拉取分组微博-日", Crontab: "0 4/4 7-23 * * *", Fn: func() { common_x.RunCatching(PullByGroupLock) }},
		{Name: "拉取分组微博-夜", Crontab: "0 4/20 0-6 * * *", Fn: func() { common_x.RunCatching(PullByGroupLock) }},
		{Name: "删除weibo已读", Crontab: "0 1/2 * * * *", Fn: func() { common_x.RunCatching(DeleteLock) }},
	}
}

func DeleteLock() {
	service.RunnerLocker.LockRun("删除weibo已读", time.Minute*2, Delete)
}
func Delete() {
	var feeds []map[string]any
	db.FreshrssDb.Table("feed").
		Select("id").
		Find(&feeds, "url like '%/weibo/rss/%'")

	sql := `
 select guid
 from entry
 where true
    and guid <= (select guid
              from entry
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
		rows, _ := db.FreshrssDb.Raw(sql, rss["id"], rss["id"]).Rows()
		for rows.Next() {
			var guid string
			_ = rows.Scan(&guid)
			guids, mysqlGuids = append(guids, cast_x.ToInt64(guid)), append(mysqlGuids, guid)
		}
	}
	if len(guids) == 0 {
		return
	}

	filter := bson.M{"_id": bson.M{"$in": guids}}
	opts := options.Find().
		SetProjection(bson.M{"_id": 1, "retweeted_status": 1}).
		SetLimit(300)
	rowsMap, _ := mongodb.RowsMap(context.TODO(), collection, filter, opts)
	for _, row := range *rowsMap {
		if cell, ok := row["retweeted_status"]; ok {
			guids, mysqlGuids = append(guids, cast_x.ToInt64(cell)), append(mysqlGuids, cast_x.ToString(cell))
		}
	}

	go db.FreshrssDb.Table("entry").Delete(nil, "guid in ? and is_favorite = 0", mysqlGuids)

	if len(guids) > 0 {
		filter = bson.M{"_id": bson.M{"$in": guids}}
		update := bson.M{"$set": bson.M{
			"invalid": 1,
			"read":    1,
		}}
		dr, err := collection.UpdateMany(context.TODO(), filter, update)

		// dr, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			return
		}
		log.Infoln("mongodb删除数量:", dr.ModifiedCount)
	}

	go func() {
		collection.UpdateMany(context.TODO(), bson.M{"groupId": 3639801313908027, "invalid": 0, "pic_ids": bson.M{"$size": 0}}, bson.M{"$set": bson.M{"invalid": 1, "read": 1}})
		collection.UpdateMany(context.TODO(), bson.M{"groupId": 3639801313908027, "invalid": 1, "read": 0}, bson.M{"$set": bson.M{"invalid": 0}})
	}()
}
