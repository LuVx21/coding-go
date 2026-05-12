package rss

import (
	"log/slog"
	"slices"
	"strings"
	"time"

	"luvx/gin/dao/freshrss_dao"
	"luvx/gin/db"
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/maths_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

var (
	maxEntryID int64 = -1
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		service.NewRunner("重置rss", "0 3/5 * * * *", time.Minute*10, reset),
		service.NewRunner("rss_spider", "0 7 0/2 * * *", time.Minute*10, PullByKey),
		service.NewRunner("新版本时间拉新", "23 29 4/12 * * *", time.Minute*10, pullLatest),
		service.NewRunner("bingTag", "33 33/5 * * * *", time.Minute*10, bindTag),
	}
}

func reset() {
	db.FreshrssDb.Exec(`update ` + freshrss_dao.Prefix + `feed set lastUpdate = lastUpdate-30*60 where url like '%/weibo/rss/%'`)
}

func bindTag() {
	if maxEntryID <= 0 {
		db.FreshrssDb.Table(freshrss_dao.Prefix + "entrytag").
			Select("id_entry").
			Order("id_entry desc").
			Limit(1).
			Scan(&maxEntryID)
	}

	type aa struct {
		ID       int64
		Tags     string
		tagNames []string
	}
	var entries []aa
	db.FreshrssDb.Raw(`select id, tags from `+freshrss_dao.Prefix+`entry
where id >= ?
and is_read = 0
and tags IS NOT NULL
and tags != ''
and id_feed in (select id from `+freshrss_dao.Prefix+`feed where url like ? and category != 2 and category != 11)
order by id limit ?`, maxEntryID+1, "%/weibo/rss/%", 100).
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
	getIdByTag := func(names []string) a.Rows {
		var rows a.Rows
		db.FreshrssDb.Table(freshrss_dao.Prefix+"tag").
			Select("name, id").
			Where("name in ?", names).
			Scan(&rows)
		return rows
	}
	rows := getIdByTag(alltags)
	existTagMap := slices_x.ToMap(rows, func(_ int, r a.Row) string { return r["name"].(string) }, func(_ int, r a.Row) int64 { return cast_x.ToInt64(r["id"]) })
	needInsertTags := slices_x.Filter(alltags, func(s string) bool {
		_, ok := existTagMap[s]
		return !ok
	})

	if len(needInsertTags) > 0 {
		for _, tag := range needInsertTags {
			db.FreshrssDb.Table(freshrss_dao.Prefix + "tag").
				Create(a.Row{
					"name":       tag,
					"attributes": "[]",
				})
		}
		ttt := getIdByTag(needInsertTags)
		for _, row := range ttt {
			existTagMap[row["name"].(string)] = cast_x.ToInt64(row["id"])
		}
	}

	for _, entry := range entries {
		var rows a.Rows
		db.FreshrssDb.Table(freshrss_dao.Prefix+"entrytag").
			Where("id_entry = ?", entry.ID).
			Scan(&rows)
		bindTagIds := slices_x.Transfer(func(r a.Row) int64 { return cast_x.ToInt64(r["id_tag"]) }, rows...)
		for _, tagName := range entry.tagNames {
			tagID := existTagMap[tagName]
			if tagID == 0 {
				slog.Error("tag不存在", "tagName", tagName)
			}
			if tagID == 0 || slices.Contains(bindTagIds, tagID) {
				continue
			}

			result := db.FreshrssDb.Table(freshrss_dao.Prefix + "entrytag").
				Create(a.Row{
					"id_entry": entry.ID,
					"id_tag":   tagID,
				})
			if result.Error != nil {
				slog.Error("绑定tag错误", "entryId", entry.ID, "tagName", tagName)
			}
		}
		maxEntryID = maths_x.Max(maxEntryID, cast_x.ToInt64(entry.ID))
	}
}
