package service

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	kv "luvx_service_sdk/proto_gen/proto_kv"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"github.com/klauspost/compress/zstd"
	"github.com/luvx21/coding-go/coding-common/retry"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/infra_kv"
	badgers "github.com/luvx21/coding-go/infra/infra_kv/badgers"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	once       sync.Once
	db         *badger.DB
	encoder, _ = zstd.NewWriter(nil)
	decoder, _ = zstd.NewReader(nil)
)

type KVServiceImpl struct {
	kv.UnimplementedKVServer
}

func (s *KVServiceImpl) Put(ctx context.Context, req *kv.PutRequest) (resp *kv.PutResponse, err error) {
	_ = initDb()

	key, bytes, exp := req.Entry.Key, req.Entry.Value, req.Expire
	bytes = encoder.EncodeAll(bytes, nil)

	msg, err := retry.SupplyWithRetry("rpc-server kv put", func() string {
		err = badgers.Set(db, []byte(key), bytes, time.Duration(exp)*time.Second)
		if err != nil {
			panic("fast-fail retry:" + err.Error())
		}
		return "ok"
	}, 5, 3*time.Second)

	return &kv.PutResponse{Message: msg}, err
}
func (s *KVServiceImpl) Get(ctx context.Context, req *kv.Key) (resp *kv.Entry, err error) {
	_ = initDb()

	key := req.Key
	bytes, exist := badgers.GetStr(db, key)
	if !exist {
		return &kv.Entry{Key: key}, badger.ErrKeyNotFound
	}
	if len(bytes) > 0 {
		bytes, err = decoder.DecodeAll(bytes, nil)
	}
	return &kv.Entry{Key: key, Value: bytes}, err
}
func (s *KVServiceImpl) List(context.Context, *emptypb.Empty) (*kv.Entries, error) {
	_ = initDb()
	m, err := badgers.List(db)
	for k, v := range m {
		bytes, _ := decoder.DecodeAll(v, nil)
		m[k] = bytes
	}
	return &kv.Entries{EntryMap: m}, err
}
func (s *KVServiceImpl) ListByPrefix(ctx context.Context, req *kv.Key) (*kv.Entries, error) {
	_ = initDb()
	m, err := badgers.ListByPrefix(db, []byte(req.Key))
	for k, v := range m {
		bytes, _ := decoder.DecodeAll(v, nil)
		m[k] = bytes
	}
	return &kv.Entries{EntryMap: m}, err
}
func (s *KVServiceImpl) ListKeyByPrefix(_ context.Context, in *kv.Key) (*kv.Keys, error) {
	_ = initDb()
	array, err := badgers.ListKeyByPrefixStream(db, []byte(in.Key))
	keys := slices_x.Transfer(func(bytes []byte) string { return string(bytes) }, array...)
	return &kv.Keys{Keys: keys}, err
}
func (s *KVServiceImpl) Delete(ctx context.Context, in *kv.Keys) (*kv.Keys, error) {
	_ = initDb()
	keys, err := badgers.DeleteStr(db, in.Keys...)
	return &kv.Keys{Keys: keys}, err
}

func initDb() (err error) {
	if db == nil {
		once.Do(func() {
			opts := badger.DefaultOptions(filepath.Join(infra_kv.KV_MAIN_DIR, "badger", "kv_service.db"))
			opts.Compression = options.ZSTD
			db, err = badger.Open(opts)
		})
		if err != nil {
			slog.Error("初始化badger db 异常", "err", err)
		}
		return err
	}
	return nil
}

func CleanUpKv() {
	defer func(db *badger.DB) {
		if db != nil {
			_ = db.Close()
		}
	}(db)
}
