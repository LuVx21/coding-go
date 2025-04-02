package badger

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/luvx21/coding-go/coding-common/common_x"
)

var db *badger.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		dbFilePath := filepath.Join(home, "data", "kv", "badger", "badger.db")

		db, _ = badger.Open(badger.DefaultOptions(dbFilePath))
	}

	return func() {
		defer db.Close()
		fmt.Println(caseName, "teardown......")
	}
}

func Test_badger_00(t *testing.T) {
	defer beforeAfter("Test_badger_00")()

	key, value := "foo", "bar"

	SetStr(db, key, value)

	v, _ := GetStr(db, key)

	fmt.Println("get 值", string(v))

	m, _ := List(db)
	for k, v := range m {
		fmt.Println(k, "=", string(v))
	}
}

func Test_badger_01(t *testing.T) {
	defer beforeAfter("Test_badger_01")()

	key, value := "foo", "bar"

	for i := range 1000 {
		SetStr(db, key+strconv.Itoa(i), value)
		// DeleteStr(db, key+strconv.Itoa(i))
	}
}
