package weibo_p

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/dao/freshrss_dao"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/dao/redis_dao"
	"luvx/gin/db"
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/common_x/t"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
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
		// service.NewRunner("拉取微博热搜", "0 7/10 * * * *", time.Minute*7, PullHotBand),
		{Name: "拉取分组微博-日", Crontab: "0 */4 7-23 * * *", Fn: func() { common_x.RunCatching(PullByGroupLock) }},
		{Name: "拉取分组微博-夜", Crontab: "0 */20 0-6 * * *", Fn: func() { common_x.RunCatching(PullByGroupLock) }},
		{Name: "删除weibo已读", Crontab: "0 */2 * * * *", Fn: func() { common_x.RunCatching(DeleteLock) }},
		service.NewRunner("bingTag", "33 */5 * * * *", time.Minute*10, bindTag),
	}
}

func DeleteLock() {
	service.RunnerLocker().LockRun("删除weibo已读", time.Minute*2, delete)
}
func delete() {
	go func() {
		freshrss_dao.DeleteEntry(slices_x.Transfer(func(i int64) string { return cast_x.ToString(i) }, mongo_dao.IgnoreRetweet()...))
	}()

	go func() {
		collection.UpdateMany(context.TODO(), bson.M{"groupId": 3639801313908027, "invalid": 0, "pic_ids": bson.M{"$size": 0}}, bson.M{"$set": bson.M{"invalid": 1, "read": 1}})
		collection.UpdateMany(context.TODO(), bson.M{"groupId": 3639801313908027, "invalid": 1, "read": 0}, bson.M{"$set": bson.M{"invalid": 0}})
	}()

	feedIds := freshrss_dao.FeedIds()
	mysqlGuids, guids := make([]string, 0), make([]int64, 0)
	for _, feedId := range feedIds {
		guidsById := freshrss_dao.SelectGuid(feedId)
		for _, guid := range guidsById {
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

	go freshrss_dao.DeleteEntry(mysqlGuids)

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
}

func bindTag() {
	freshrss_dao.DeleteLessTag()
	freshrss_dao.DeleteDanglingTag()

	maxEntryID := freshrss_dao.MaxEntryId()

	type aa struct {
		ID       int64
		Content  string
		Tags     string
		tagNames []string
	}
	var entries []aa
	db.FreshrssDb.Raw(`select id, content, tags from `+freshrss_dao.Prefix+`entry
where id >= ?
  and is_read = 0
  and tags IS NOT NULL
  and tags != ''
  and id_feed in (select id from `+freshrss_dao.Prefix+`feed where url like ? and category != 2 and category != 11)
order by id limit ?`, maxEntryID+1, "%/weibo/rss/%", 200).
		Scan(&entries)
	if len(entries) == 0 {
		return
	}

	set := sets.NewSet[string]()
	for i := range entries {
		entry := &entries[i]
		entry.tagNames = slices_x.Transfer(func(s string) string {
			return strings.ToLower(strings.TrimPrefix(s, "#"))
		}, strings.Split(entry.Tags, " ")...)
		set.Add(entry.tagNames...)
	}
	alltags := set.ToSlice()
	existTagMap := slices_x.ToMap(freshrss_dao.SelectTagByNames(alltags),
		func(_ int, r a.Row) string { return r["name"].(string) },
		func(_ int, r a.Row) int64 { return cast_x.ToInt64(r["id"]) },
	)
	needInsertTags := slices_x.Filter(alltags, func(s string) bool {
		_, ok := existTagMap[s]
		return !ok
	})

	// 存新Tag
	if len(needInsertTags) > 0 {
		for _, tag := range needInsertTags {
			pk := freshrss_dao.SaveTag(tag)
			existTagMap[tag] = pk
		}
		// ttt := freshrss_dao.SelectTagByNames(needInsertTags)
		// for _, row := range ttt {
		// 	existTagMap[row["name"].(string)] = cast_x.ToInt64(row["id"])
		// }
	}

	prefix := "标签:"
	for _, entry := range entries {
		needAdd := len(entry.tagNames) > 0 && !strings.HasPrefix(entry.Content, prefix)
		tagIds := make([]t.Pair[int64, string], 0)

		rows := freshrss_dao.SelectEntryTag(entry.ID, -1)
		bindedTagIds := slices_x.Transfer(func(r a.Row) int64 { return cast_x.ToInt64(r["id_tag"]) }, rows...)
		for _, tagName := range entry.tagNames {
			tagID, ok := existTagMap[tagName]
			if !ok || tagID == 0 {
				tagID = freshrss_dao.SaveTag(tagName)
			}
			if needAdd {
				tagIds = append(tagIds, t.NewPair(tagID, tagName))
			}
			if tagID == 0 || slices.Contains(bindedTagIds, tagID) {
				continue
			}

			if freshrss_dao.SaveEntryTag(entry.ID, tagID) != nil {
				slog.Error("绑定tag错误", "entryId", entry.ID, "tagName", tagName)
			}
		}
		if needAdd {
			var sb strings.Builder
			sb.WriteString(prefix)
			for i, p := range tagIds {
				fmt.Fprintf(&sb, `<a href="http://localhost:50080/i/?a=normal&get=t_%d">#%s<a/>`+common_x.IfThen(i < len(tagIds)-1, strings.Repeat(consts.Nbsp, 4), ""), p.K, p.V)
			}
			sb.WriteString("<br/>")
			db.FreshrssDb.Exec(`update `+freshrss_dao.Prefix+`entry set content = concat(?::text, content) where id = ? and content not like ?`, sb.String(), entry.ID, prefix+"%")
		}

	}
}
