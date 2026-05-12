package freshrss_dao

import (
	"luvx/gin/db"
)

const (
	mysql_prefix = "freshrss.t_admin_"
	Prefix       = ""
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
