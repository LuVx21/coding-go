package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	infra_bolt "github.com/luvx21/coding-go/infra/infra_kv/bolt"
	"github.com/urfave/cli/v3"
	bolt "go.etcd.io/bbolt"
)

const (
	DBDir, DBName = ".kv", "store.db"
	BUCKET        = "default"
)

func init() {
	mustHaveDBFile()
	mustHaveBucket()
}

func listAllBuckets(db *bolt.DB) ([]string, error) {
	if db == nil {
		db = openDB()
		defer db.Close()
	}

	return infra_bolt.ListBucket(db)
}

func getAll(db *bolt.DB, bucket string) (map[string]string, error) {
	if db == nil {
		db = openDB()
		defer db.Close()
	}

	m, err := infra_bolt.List(db, common_x.IfThen(len(bucket) != 0, bucket, BUCKET))
	if err != nil {
		return nil, errors.New("获取所有键值对出现错误")
	}
	if len(m) == 0 {
		return nil, errors.New("空")
	}
	return m, err
}

func get(bucket, key string) (string, error) {
	db := openDB()
	defer db.Close()
	v, _ := infra_bolt.Get(db, common_x.IfThen(len(bucket) != 0, bucket, BUCKET), key)
	return v, nil
}

func set(bucket, key, value string) error {
	db := openDB()
	defer db.Close()

	if err := infra_bolt.Set(db, common_x.IfThen(len(bucket) != 0, bucket, BUCKET), key, value); err != nil {
		return fmt.Errorf("设置键值对错误-> %s:%s", key, value)
	}
	return nil
}

func del(bucket, key string) error {
	db := openDB()
	defer db.Close()

	err := infra_bolt.Del(db, common_x.IfThen(len(bucket) != 0, bucket, BUCKET), key)
	if err != nil {
		return fmt.Errorf("删除键值对错误.key: %s", key)
	}
	return nil
}
func backup(ctx context.Context, c *cli.Command) error {
	db := openDB()
	defer db.Close()

	home, _ := common_x.Dir()
	bs, err := listAllBuckets(db)
	if err != nil {
		return err
	}
	fi, err := os.OpenFile(filepath.Join(home, DBDir, "kv.sh"), os.O_TRUNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()
	for _, b := range bs {
		kvs, err := getAll(db, b)
		if err != nil {
			fmt_x.Errorf("获取bucket内容错误:%s\n", "err", err)
			break
		}
		for k, v := range kvs {
			fmt.Fprintf(fi, "kv set %s@%s %s\n", k, b, v)
		}
	}

	return nil
}

func openDB() *bolt.DB {
	home, _ := common_x.Dir()
	dbFilePath := filepath.Join(home, DBDir, DBName)
	db, err := bolt.Open(dbFilePath, os.FileMode(0600), nil)
	checkError(err)
	return db
}

func mustHaveDBFile() {
	home, _ := common_x.Dir()
	dbFilePath := filepath.Join(home, DBDir, DBName)
	if _, err := os.Stat(dbFilePath); err != nil {
		dir := filepath.Dir(dbFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0766)
		}
		if _, err := os.Create(dbFilePath); err != nil {
			panic(err)
		}
	}
}

func mustHaveBucket() {
	db := openDB()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BUCKET))
		if err != nil {
			return fmt.Errorf("创建分区错误: %s", err)
		}
		return nil
	})
}
