package main

import (
	"fmt"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	// dsn = "%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = "postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Shanghai"

	host, port     = "mini.local", 5432
	user, password = "xxx", "xxx"
	dbname         = "boot"
)

type Querier interface {
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	database := common_x.IfThen(false, mysql.Open, postgres.Open)
	db, _ := gorm.Open(database(fmt.Sprintf(dsn, user, password, host, port, dbname)))

	genRun(db, "user", "common_key_value")
}
func genRun(db *gorm.DB, tableNames ...string) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../gin/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	// g.ApplyBasic(g.GenerateAllTable()...)
	models := slices_x.Transfer(func(n string) any { return g.GenerateModel(n) }, tableNames...)
	g.ApplyBasic(models...)

	//g.ApplyInterface(func(Querier) {}, model.User{}, model.CommonKeyValue{})
	g.Execute()
}
