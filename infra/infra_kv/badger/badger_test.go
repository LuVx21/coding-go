package badger

import (
	"fmt"
	"path/filepath"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/luvx21/coding-go/coding-common/common_x"
)

var db *badger.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		dbFilePath := filepath.Join(home, "data", "badger", "badger.db")

		db, _ = badger.Open(badger.DefaultOptions(dbFilePath))
	}

	return func() {
		defer db.Close()
		fmt.Println(caseName, "test case end...")
	}
}

func Test_badger_00(t *testing.T) {
	key, value := "foo", "bar"
	defer beforeAfter("Test_badger_00")()

	_ = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})

	var v []byte
	_ = db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte(key))
		item.Value(func(val []byte) error {
			v = val
			return nil
		})
		return nil
	})
	fmt.Println("get 值", string(v))
}
