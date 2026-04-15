package db

import (
	"fmt"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"gorm.io/gorm"
)

func Test_00(t *testing.T) {
	_sql := `
`

	sql := MySQLClient.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return MySQLClient.Exec(_sql)
	})
	fmt.Println(sql)
}

func Test_mysql(t *testing.T) {
	var results a.SAMS
	MySQLClient.Debug().Raw("select * from user order by id;").
		//Scan(&results)
		Find(&results)
	fmt.Println(results)
	//client.Exec("delete from ")
}
