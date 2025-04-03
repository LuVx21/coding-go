package common_kv_dao

import (
	"luvx/gin/db"
	"luvx/gin/model"
)

func JsonArrayAppend(id int64, path string, value any) {
	_sql := `
update common_key_value
set common_value = json_array_append(common_value, ?, ?)
where id = ?;
`
	db.MySQLClient.Exec(_sql, path, value, id)
}

func Get(bizType int32, keys ...string) []*model.CommonKeyValue {
	client := db.MySQLClient
	var kvs []*model.CommonKeyValue
	tx := client.Debug()
	tx = tx.Where("biz_type = ? and invalid = 0", bizType)
	if len(keys) > 0 {
		tx = tx.Where("common_key in ?", keys)
	}
	_ = tx.Find(&kvs).Error
	return kvs
}
