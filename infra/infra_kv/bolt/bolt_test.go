package bolt

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
	bolt "go.etcd.io/bbolt"
)

const bucket = "test"

var db *bolt.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		dbFilePath := filepath.Join(home, "data", "kv", "bolt", "bolt.db")

		db, _ = OpenDBWithBucket(dbFilePath, bucket)
	}

	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func Test_bolt_00(t *testing.T) {
	defer beforeAfter("Test_bolt_00")()

	key, value := "foo", "bar"

	defer db.Close()

	Set(db, bucket, key, value)

	v, _ := Get(db, bucket, key)
	fmt.Println(value, v)

	fmt.Println(ListBucket(db))
}
