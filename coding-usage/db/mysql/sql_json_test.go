package mysql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luvx21/coding-go/infra/infra_sql/infra_mysql"
	"gorm.io/gorm"
)

func Test_json_map_select(t *testing.T) {
	defer beforeAfter("Test_json_map")()
	var ff int = 0
	err := gormDB.Debug().
		Table("common_key_value").
		Where("id = ?", 91).
		Select(`common_value ->> '$._43510.expireAt' as common_value`).
		Find(&ff).
		Error

	fmt.Println(err, ff)
}

func Test_json_map(t *testing.T) {
	defer beforeAfter("Test_json_map")()
	tx := gormDB.Debug().
		Table("common_key_value").
		Where("id = ?", 4).
		Update("common_value", gorm.Expr(infra_mysql.JSON_SET+"(common_value, ?, ?, ?, ?)", "$.id", 4, "$.userName", "bar44"))
	fmt.Println("异常", tx.Error)
}

func Test_json_map_set(t *testing.T) {
	defer beforeAfter("Test_json_map_set")()
	tx := gormDB.Debug().
		Table("common_key_value").
		Where("id = ?", 4).
		Update("common_value", gorm.Expr("JSON_SET(common_value, ?, ?, ?, ?, ?, JSON_ARRAY(?))", "$.id", 4, "$.userName", "bar4", "$.nums", []int{22, 11}))
	fmt.Println("异常", tx.Error)
}

func Test_json_map_insert(t *testing.T) {
	defer beforeAfter("Test_json_map_insert")()
	tx := gormDB.Debug().
		Table("common_key_value").
		Where("id = ?", 4).
		Update("common_value", gorm.Expr("JSON_INSERT(common_value, ?, ?, ?, ?)", "$.id", 44, "$.userName4", "barbar"))
	fmt.Println("异常", tx.Error)
}

func Test_json_map_replace(t *testing.T) {
	defer beforeAfter("Test_json_map_replace")()
	gormDB.Debug().Where("id = ?", 4)
	tx := gormDB.Debug().
		Table("common_key_value").
		Where("id = ?", 4).
		Update("common_value", gorm.Expr("JSON_REPLACE(common_value, ?, ?, ?, ?)", "$.id", 44, "$.userName5", "barbar"))
	fmt.Println("异常", tx.Error)
}
