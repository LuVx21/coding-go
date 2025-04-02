package bolt

import (
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

func OpenDBWithBucket(dbFilePath, bucket string) (*bolt.DB, error) {
	if _, err := os.Stat(dbFilePath); err != nil {
		dir := filepath.Dir(dbFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0766)
		}
		if _, err := os.Create(dbFilePath); err != nil {
			panic(err)
		}
	}

	db, err := bolt.Open(dbFilePath, os.FileMode(0600), &bolt.Options{
		NoFreelistSync: true,
	})
	if err != nil {
		return nil, err
	}

	if len(bucket) > 0 {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("创建分区错误: %s", err)
			}
			return nil
		})
	}

	return db, err
}

func CreateBucket(db *bolt.DB, bucket string) (*bolt.Bucket, error) {
	var b *bolt.Bucket
	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("创建分区错误: %s", err)
		}
		return nil
	})
	return b, err
}

func listBuckets(b *bolt.Bucket, indent string) {
	if b == nil {
		return
	}
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			fmt.Printf("%sBucket: %s\n", indent, string(k))
			listBuckets(b.Bucket(k), indent+"  ")
		}
		return nil
	})
}

func ListBucketDIGUI(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		listBuckets(tx.Cursor().Bucket(), "")
		return nil
	})
}

// ListBucket 列出所有顶级 Bucket
func ListBucket(db *bolt.DB) ([]string, error) {
	r := make([]string, 0)
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			r = append(r, string(name))
			return nil
		})
	})
	return r, err
}

func List(db *bolt.DB, bucket string) (map[string]string, error) {
	m := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				m[string(k)] = string(v)
			}
		}
		return nil
	})
	return m, err
}

func Set(db *bolt.DB, bucket, key, value string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("创建分区错误: %s", err)
			}
		}
		e := b.Put([]byte(key), []byte(value))
		return e
	})
}

func Get(db *bolt.DB, bucket, key string) (string, error) {
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("获取键值对错误,key: %s", key)
	}
	if v == nil {
		return "", fmt.Errorf("键值对不存在,key: %s", key)
	}
	return string(v), nil
}

func Del(db *bolt.DB, bucket, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		e := b.Delete([]byte(key))
		return e
	})
}
