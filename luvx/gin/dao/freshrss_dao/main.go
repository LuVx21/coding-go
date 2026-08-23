package freshrss_dao

import (
	"log/slog"
	"luvx/gin/config"
	"luvx/gin/db"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

const (
	mysql_prefix = "freshrss.t_admin_"
)

var (
	Prefix = common_x.IfThen(config.GetSwitch("rss.freshrssSqlite"), "", mysql_prefix)
)

var (
	_sql = `
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
)

func ExistedGuids(path string, guids []string) []string {
	if len(guids) == 0 || path == "" {
		return []string{}
	}
	var r []string
	db.FreshrssDb.Raw(_sql, path, guids).Scan(&r)
	return r
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
func DeleteUntag() {
	db.FreshrssDb.Exec(`
delete
from ` + Prefix + `entrytag
where not exists (
    select 1
    from ` + Prefix + `entry
    where id_entry=id
);
	`)
	db.FreshrssDb.Exec(`
delete
from ` + Prefix + `tag
where not exists (
    select 1
    from ` + Prefix + `entrytag
    where id_tag=id
);
	`)
}
