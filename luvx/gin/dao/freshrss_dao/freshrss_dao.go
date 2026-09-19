package freshrss_dao

import (
	"log/slog"
	"luvx/gin/config"
	"luvx/gin/db"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

const (
	mysql_prefix    = "freshrss.t_admin_"
	postgres_prefix = "t_admin_"
)

var (
	Prefix = common_x.IfThen(config.GetSwitch([]string{"rss", "freshrssPostgres"}), postgres_prefix, "")
)

func ExistedGuids(path string, guids []string) []string {
	if len(guids) == 0 || path == "" {
		return []string{}
	}
	_sql := `
select guid
from ` + Prefix + `entry
where true
  and id_feed in (
    select id
    from ` + Prefix + `feed
    where true
    and url like ?
)
and guid in ?
limit 200
`
	var r []string
	db.FreshrssDb.Raw(_sql, path, guids).Scan(&r)
	return r
}
func SelectGuid(feedId int64) []string {
	var guids []string
	sql := `
 select guid
 from ` + Prefix + `entry
 where true
    and guid <= (select guid
              from ` + Prefix + `entry
              where true
                and id_feed = ?
              order by guid desc
              limit 1)
   and id_feed = ?
   and is_read = 1
   and is_favorite = 0
-- order by guid
 limit 100
`
	db.FreshrssDb.Raw(sql, feedId, feedId).Scan(&guids)
	return guids
}
func DeleteEntryByFeed(feedId int64) {
	if feedId <= 0 {
		return
	}

	db.FreshrssDb.Exec(`
delete from `+Prefix+`entry
where id_feed = ? and is_read = 0 and is_favorite = 0
`, feedId)
}
func DeleteEntry(guids []string) {
	for _, guid := range guids {
		err := db.FreshrssDb.Table(Prefix+"entry").Delete(nil, "guid = ? and is_favorite = 0", guid).Error
		if err != nil {
			slog.Error("delete entry by guid", "err", err)
		}
	}
}

func FeedIds() []int64 {
	// var feeds []map[string]any
	var feeds []int64
	db.FreshrssDb.Table(Prefix+"feed").
		Select("id").
		Find(&feeds, "url like '%/weibo/rss/%'")
	return feeds
}

// 删除空悬的tag和绑定
func DeleteDanglingTag() {
	db.FreshrssDb.Exec(`
delete from ` + Prefix + `entrytag et
where not exists (
	select 1
	from ` + Prefix + `entry e
	where et.id_entry = e.id
);
	`)
	db.FreshrssDb.Exec(`
delete from ` + Prefix + `tag t
where not exists (
	select 1
	from ` + Prefix + `entrytag et
	where et.id_tag = t.id
);
	`)
}

// 删除绑定少的tag
func DeleteLessTag() {
	db.FreshrssDb.Exec(`
delete from ` + Prefix + `tag t
where exists (
	select 1
	from ` + Prefix + `entrytag et
	where et.id_tag = t.id
	group by et.id_tag
	having count(et.id_entry) < 4
);
`)
}

func SelectEntryTag(entryId, tagId int64) a.Rows {
	if entryId <= 0 && tagId <= 0 {
		return nil
	}

	tx := db.FreshrssDb.Table(Prefix + "entrytag")
	if entryId > 0 {
		tx.Where("id_entry = ?", entryId)
	}
	if tagId > 0 {
		tx.Where("id_tag = ?", tagId)
	}
	var rows a.Rows
	tx.Scan(&rows)
	return rows
}
func SaveEntryTag(entryId, tagId int64) error {
	err := db.FreshrssDb.Table(Prefix + "entrytag").
		Create(a.Row{
			"id_entry": entryId,
			"id_tag":   tagId,
		}).Error
	if err != nil {
		slog.Error("绑定tag错误", "entryId", entryId, "tagId", tagId, "error", err)
		return err
	}
	return err
}
func SelectTagByNames(names []string) a.Rows {
	var rows a.Rows
	db.FreshrssDb.Table(Prefix+"tag").
		Select("name, id").
		Where("name in ?", names).
		Scan(&rows)
	return rows
}
func SaveTag(tag string) int64 {
	row := a.Row{
		"name":       tag,
		"attributes": "[]",
	}
	db.FreshrssDb.Table(Prefix + "tag").
		Clauses(clause.Returning{}). // 给map结构回填主键
		Create(row)
	return cast.ToInt64(row["id"])
}

func MaxEntryId() int64 {
	var maxEntryID int64
	db.FreshrssDb.Table(Prefix + "entrytag").
		Select("id_entry").
		Order("id_entry desc, id_tag desc").
		Limit(1).
		Scan(&maxEntryID)
	return maxEntryID
}
