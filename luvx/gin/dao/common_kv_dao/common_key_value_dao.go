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
	tx := db.MySQLClient
	var kvs []*model.CommonKeyValue
	// tx := client.Debug()
	tx = tx.Where("biz_type = ? and invalid = 0", bizType)
	if len(keys) > 0 {
		tx = tx.Where("common_key in ?", keys)
	}
	_ = tx.Find(&kvs).Error
	return kvs
}

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
