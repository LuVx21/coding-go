package kv

import (
	"testing"

	"github.com/rosedblabs/rosedb/v2"
)

func Test_rosedb_00(t *testing.T) {
	options := rosedb.DefaultOptions
	options.DirPath = "/tmp/rosedb_basic"

	// 打开数据库
	db, err := rosedb.Open(options)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	k, v := "foo", "bar"

	// 设置键值对
	err = db.Put([]byte(k), []byte(v))
	if err != nil {
		panic(err)
	}

	// 获取键值对
	val, err := db.Get([]byte(k))
	if err != nil {
		panic(err)
	}
	println(string(val))

	// 删除键值对
	err = db.Delete([]byte(k))
	if err != nil {
		panic(err)
	}

	batch := db.NewBatch(rosedb.DefaultBatchOptions)
	_ = batch.Put([]byte(k), []byte(v+v))
	val, _ = batch.Get([]byte(k))
	println(string(val))

	_ = batch.Delete([]byte(k))
	_ = batch.Commit()
}
