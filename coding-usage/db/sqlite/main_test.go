package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	"gorm.io/gorm"

	//_ "modernc.org/sqlite"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	gorm_sqlite "gorm.io/driver/sqlite"
)

const (
	//driverName = "sqlite"
	driverName = "sqlite3"

	path  = "/data/sqlite/main.db"
	path1 = "/docker/freshrss/data/users/admin/db1.sqlite"
)

var (
	db     *sql.DB
	gormDB *gorm.DB
)

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		_url := home + path1
		db, _ = sql.Open(driverName, _url)
	}

	if gormDB == nil {
		gormDB, _ = gorm.Open(gorm_sqlite.New(gorm_sqlite.Config{Conn: db}), &gorm.Config{})
	}

	return func() {
		fmt.Println(caseName, "end...")
	}
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	rows, _ := db.Query("SELECT * FROM user where id >= $1", 1)
	defer rows.Close()
	dbs.PrintRows(rows)
}

func Test_01(t *testing.T) {
	defer beforeAfter("Test_01")()
	// var feeds []map[string]any
	var feeds []int64

	gormDB.Table("feed").
		Select("id").
		Find(&feeds, "url like '%/weibo/rss/%'")
	fmt.Println(feeds)
}
