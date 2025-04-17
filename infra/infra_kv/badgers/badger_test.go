package badgers

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/test"
)

var (
	db          *badger.DB
	beforeAfter = func(caseName string) func() {
		return test.BeforeAfterTest(caseName, func() {
			if db != nil {
				return
			}
			dbFilePath := filepath.Join(common_x.Home(), "data", "kv", "badger", "badger.db")
			db, _ = badger.Open(badger.DefaultOptions(dbFilePath))
		}, func() { db.Close() })
	}
)

func Test_badger_00(t *testing.T) {
	defer beforeAfter("Test_badger_00")()

	key, value := "foo", "bar"

	SetStr(db, key, value, 0)

	v, _ := GetStr(db, key)

	fmt.Println("get 值", string(v))
}

func Test_badger_01(t *testing.T) {
	defer beforeAfter("Test_badger_01")()

	key, value := "foo", "bar"

	for i := range 1000 {
		SetStr(db, key+strconv.Itoa(i), value, 0)
		// DeleteStr(db, key+strconv.Itoa(i))
	}
}

func Test_badger_02(t *testing.T) {
	defer beforeAfter("Test_badger_02")()

	key, value := "foo", "bar"
	for i := range 20 {
		SetStr(db, key+":"+strconv.Itoa(i), value, 0)
	}

	m, _ := List(db)
	for k, v := range m {
		fmt.Println(k, "=", string(v))
	}

	m1, _ := ListByPrefix(db, []byte(key+":"))
	for k, v := range m1 {
		fmt.Println(k, "=", string(v))
	}
}

func Test_badger_03(t *testing.T) {
	defer beforeAfter("Test_badger_03")()

	stream := db.NewStream()
	stream.NumGo = 16
	stream.Prefix = []byte("foo:")
	stream.LogPrefix = "Badger.Streaming"

	// stream.ChooseKey = func(item *badger.Item) bool {
	// 	return bytes.HasPrefix(item.Key(), []byte("foo:"))
	// }

	stream.Send = func(buf *z.Buffer) error {
		list, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range list.Kv {
			if kv.StreamDone == true {
				return nil
			}
			fmt.Println(string(kv.Key), "=", string(kv.Value))
			// cp := proto.Clone(kv).(*pb.KV)
		}
		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		slog.Error("stream orchestrate error", "err", err)
	}
}
