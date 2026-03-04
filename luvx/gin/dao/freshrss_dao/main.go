package freshrss_dao

import (
	"luvx/gin/db"
)

const (
	Prefix = "freshrss.t_admin_"
)

var (
	_sql = `
select guid
from entry
where true
  and id_feed in (
    select id
    from feed
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
