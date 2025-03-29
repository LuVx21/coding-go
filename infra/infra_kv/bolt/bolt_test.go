package bolt

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

func Test_bolt_00(t *testing.T) {
	bucket, key, value := "default", "foo", "bar"
	home, _ := common_x.Dir()
	dbFilePath := filepath.Join(home, "data", "bolt", "bolt.db")

	db, _ := OpenDBWithBucket(dbFilePath, bucket)
	defer db.Close()

	Set(db, bucket, key, value)

	v, _ := Get(db, bucket, key)
	fmt.Println(value, v)

	fmt.Println(ListBucket(db))
}
