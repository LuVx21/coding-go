package gorm

import (
	"fmt"
	"testing"
)

func Test_raw_00(t *testing.T) {
	defer beforeAfter("Test_raw_00")()

	var feeds []map[string]any
	gormDB.Debug().Table("freshrss.t_admin_feed").
		Select("url").
		Find(&feeds, "url like '%/weibo/rss/%'")

	fmt.Println(len(feeds))

	var m2 []map[string]any
	gormDB.Debug().Raw(`select url from freshrss.t_admin_feed where url like '%/weibo/rss/%'`).Scan(&m2)
	fmt.Println(len(m2))

	var r []string
	gormDB.Debug().Table("freshrss.t_admin_"+"entry").
		Select("guid").
		Where("guid in ?", []string{"1", "2", "3"}).
		Where("id_feed in (select id from feed where url like ?)", "%/weibo/rss/%").
		Pluck("guid", &r)

	// for rows.Next() {
	// 	fmt.Println("1111111")
	// }

	// gormDB.Raw("").Scan()
	// gormDB.Raw("").ScanRows()
	// gormDB.Raw("").Row()
	// gormDB.Raw("").Rows()
	// gormDB.Raw("").Find()

	// gormDB.Exec()
	// gormDB.Table()
	// gormDB.Model()
	// gormDB.Transaction()
}
