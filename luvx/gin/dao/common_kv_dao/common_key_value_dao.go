package common_kv_dao

import (
	"log/slog"
	"luvx/gin/db"
	"luvx/gin/model"

	"gorm.io/gorm"
)

func JsonArrayAppend(id int64, path string, value any) {
	_sql := `
update common_key_value
set common_value = json_array_append(common_value, ?, ?)
where id = ?;
`
	db.MySQLClient.Exec(_sql, path, value, id)
}

// UpdateJsonMap 操作json字段
// JSON_SET 有则覆盖, 无则添加
// JSON_INSERT 有则忽略, 无则添加
// JSON_REPLACE 有则替换, 无则忽略
func UpdateJsonMap(bizType int32, key string, expr string, args ...any) {
	err := db.MySQLClient.
		Debug().
		Model(&model.CommonKeyValue{}).
		Where("biz_type = ? and invalid = 0", bizType).
		Where("common_key = ?", key).
		Update("common_value", gorm.Expr(expr, args...)).
		Error
	if err != nil {
		slog.Error("异常结束", "Error", err)
	}
}

func Get(bizType int32, keys ...string) []*model.CommonKeyValue {
	tx := db.MySQLClient
	var kvs []*model.CommonKeyValue
	// tx := client.Debug()
	tx = tx.Where("biz_type = ? and invalid = 0", bizType)
	if len(keys) > 0 {
		if len(keys) == 1 {
			tx = tx.Where("common_key = ?", keys[0])
		} else {
			tx = tx.Where("common_key in ?", keys)
		}
	}
	_ = tx.Find(&kvs).Error
	return kvs
}

// GetByCursor id < cursorID and biz_type = bizType and common_key in keys order by id desc limit limit
func GetByCursor(cursorID int, limit int, bizType int32, keys ...string) ([]*model.CommonKeyValue, int, error) {
	if cursorID < 0 {
		return nil, 0, nil
	}

	tx := db.MySQLClient
	// tx := client.Debug()
	if cursorID > 0 {
		tx = tx.Where("id < ?", cursorID)
	}
	if bizType > 0 {
		tx = tx.Where("biz_type = ?", bizType)
	}
	// tx = tx.Where("invalid = 0")
	if len(keys) > 0 {
		tx = tx.Where("common_key in ?", keys)
	}

	tx = tx.Order("id desc").Limit(limit)

	var kvs []*model.CommonKeyValue
	err := tx.Find(&kvs).Error
	if err != nil {
		return nil, 0, err
	}

	nextCursorID := 0
	if len(kvs) > 0 {
		nextCursorID = int(kvs[len(kvs)-1].ID)
	}

	return kvs, nextCursorID, nil
}

func Create(kv *model.CommonKeyValue) error {
	return db.MySQLClient.Create(kv).Error
}

func Delete(ids []int) error {
	return db.MySQLClient.Where("id in ?", ids).Delete(&model.CommonKeyValue{}).Error
}

func Update(kv *model.CommonKeyValue) error {
	return db.MySQLClient.Save(kv).Error
}
