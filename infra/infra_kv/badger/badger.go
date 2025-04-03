package badger

import (
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func ListByPrefix(db *badger.DB) (map[string][]byte, error) {
	r := make(map[string][]byte)
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("1234")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				r[string(k)] = v
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return r, err
}

func List(db *badger.DB) (map[string][]byte, error) {
	r := make(map[string][]byte)
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				r[string(k)] = v
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return r, err
}

func SetStr(db *badger.DB, key, value string, exp time.Duration) error {
	return Set(db, []byte(key), []byte(value), exp)
}

func Set(db *badger.DB, key, value []byte, exp time.Duration) error {
	err := db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)
		if exp > 0 {
			entry.WithTTL(exp)
		}
		return txn.SetEntry(entry)
	})
	return err
}

func GetStr(db *badger.DB, key string) ([]byte, bool) {
	return Get(db, []byte(key))
}

func Get(db *badger.DB, key []byte) ([]byte, bool) {
	var v []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			v = val
			return nil
		})
		return err
	})
	if err != nil || len(v) == 0 {
		return nil, false
	}
	return v, true
}

func BatchSet(db *badger.DB, m map[string][]byte) error {
	txn := db.NewTransaction(true)
	for k, v := range m {
		kk := []byte(k)
		if err := txn.Set(kk, v); err == badger.ErrTxnTooBig {
			_ = txn.Commit()
			txn = db.NewTransaction(true)
			_ = txn.Set(kk, v)
		}
	}
	return txn.Commit()
}

func DeleteStr(db *badger.DB, keys ...string) ([]string, error) {
	bs := slices_x.Transfer(func(s string) []byte { return []byte(s) }, keys...)
	var err error
	bs, err = Delete(db, bs...)
	return slices_x.Transfer(func(b []byte) string { return string(b) }, bs...), err
}

func Delete(db *badger.DB, keys ...[]byte) ([][]byte, error) {
	deleted := make([][]byte, 0)
	err := db.Update(func(txn *badger.Txn) error {
		for _, k := range keys {
			err := txn.Delete(k)
			if err != nil {
				return err
			}
			deleted = append(deleted, k)
		}
		return nil
	})
	return deleted, err
}
