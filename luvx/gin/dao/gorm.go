package dao

import (
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"gorm.io/gorm"
)

func ToSQL(gdb *gorm.DB, _sql string, values ...any) {
	realSql := gdb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Exec(_sql, values...)
	})
	fmt_x.Infoln("填充参数后SQL: ", realSql)
}

func ToSQLSafe(gdb *gorm.DB, _sql string, values ...any) {
	stmt := gdb.Session(&gorm.Session{DryRun: true}).
		Exec(_sql, values...)
	fmt_x.Infoln("参数:", stmt.Statement.Vars)
	realSql := gdb.Dialector.Explain(stmt.Statement.SQL.String(), stmt.Statement.Vars...)
	fmt_x.Infoln("填充参数后SQL: ", realSql)
}
