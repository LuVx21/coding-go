package rosedb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/rosedblabs/rosedb/v2"
)

var db *rosedb.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()

		options := rosedb.DefaultOptions
		options.DirPath = filepath.Join(home, "data", "kv", "rosedb", "rosedb.db")
		db, _ = rosedb.Open(options)
	}

	return func() {
		defer func() {
			_ = db.Close()
		}()
		fmt.Println(caseName, "teardown......")
	}
}

func Test_rosedb_00(t *testing.T) {
	defer beforeAfter("Test_rosedb_00")()
	k, v := "foo", "bar"

	// 设置键值对
	err := db.Put([]byte(k), []byte(v))
	if err != nil {
		panic(err)
	}

	// 获取键值对
	val, _ := db.Get([]byte(k))
	println(string(val))

	// 删除键值对
	_ = db.Delete([]byte(k))

	batch := db.NewBatch(rosedb.DefaultBatchOptions)
	_ = batch.Put([]byte(k), []byte(v+v))
	val, _ = batch.Get([]byte(k))
	println(string(val))

	_ = batch.Delete([]byte(k))
	_ = batch.Commit()
}
