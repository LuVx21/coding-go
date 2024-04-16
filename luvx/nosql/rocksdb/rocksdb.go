package rocksdb

import (
    "errors"
    "github.com/linxGnu/grocksdb"
    "luvx/api/logs"
)

const (
    DB_PATH = "/Users/renxie/data/RocksDB"
)
const (
    KB int = 1024
    MB     = KB * KB
)

var (
    RocksdbClient *grocksdb.DB
)

func put(k string, v string, db *grocksdb.DB) {
    writeOptions := grocksdb.NewDefaultWriteOptions()
    writeOptions.SetSync(true)

    var key = []byte(k)
    var value = []byte(v)

    db.Put(writeOptions, key, value)
}

func delete(k string, db *grocksdb.DB) {
    writeOptions := grocksdb.NewDefaultWriteOptions()
    writeOptions.SetSync(true)

    db.Delete(writeOptions, []byte(k))
}

func get(k string, db *grocksdb.DB) {
    readOptions := grocksdb.NewDefaultReadOptions()
    readOptions.SetFillCache(true)

    slice, err := db.Get(readOptions, []byte(k))
    defer slice.Free()
    if err != nil {
        logs.Log.Println("读取数据异常:", k, err)
    }
    if slice.Exists() {
        logs.Log.Printf("获取数据-> %s=%s 值大小:%d", k, slice.Data(), slice.Size())
    }
}

func all(db *grocksdb.DB) {
    readOptions := grocksdb.NewDefaultReadOptions()
    readOptions.SetFillCache(true)

    it := db.NewIterator(readOptions)
    defer it.Close()
    for it.SeekToFirst(); it.Valid(); it.Next() {
        key := it.Key()
        value := it.Value()
        logs.Log.Printf("获取数据-> %s=%s", key.Data(), value.Data())
        key.Free()
        value.Free()
    }
    if err := it.Err(); err != nil {
    }
}

func GetDBClient() (*grocksdb.DB, error) {
    if RocksdbClient != nil {
        return RocksdbClient, nil
    }
    bloomFilter := grocksdb.NewBloomFilter(10)

    bbTableOptions := grocksdb.NewDefaultBlockBasedTableOptions()
    bbTableOptions.SetBlockCache(grocksdb.NewLRUCache(uint64(64 * KB)))
    bbTableOptions.SetFilterPolicy(bloomFilter)
    bbTableOptions.SetBlockSizeDeviation(5)
    bbTableOptions.SetBlockRestartInterval(10)
    // bbTableOptions.SetBlockCacheCompressed(grocksdb.NewLRUCache(uint64(64 * KB)))
    bbTableOptions.SetCacheIndexAndFilterBlocks(true)
    bbTableOptions.SetIndexType(grocksdb.KHashSearchIndexType)

    options := grocksdb.NewDefaultOptions()
    options.SetCreateIfMissing(true)

    readOptions := grocksdb.NewDefaultReadOptions()
    readOptions.SetFillCache(false)

    rateLimiter := grocksdb.NewRateLimiter(10000000, 10000, 10)
    options.EnableStatistics()
    options.SetBlockBasedTableFactory(bbTableOptions)
    options.SetCompactionStyle(grocksdb.UniversalCompactionStyle)
    options.SetCompression(grocksdb.SnappyCompression)
    options.SetCreateIfMissing(true)
    options.SetHashSkipListRep(2000000, 4, 4)
    options.SetMaxBackgroundCompactions(10)
    options.SetMaxWriteBufferNumber(3)
    options.SetPrefixExtractor(grocksdb.NewFixedPrefixTransform(3))
    options.SetRateLimiter(rateLimiter)
    options.SetWriteBufferSize(uint64(8 * KB))
    options.SetAllowConcurrentMemtableWrites(false)

    db, err := grocksdb.OpenDb(options, DB_PATH)

    if err != nil {
        logs.Log.Fatalln("OPEN DB error", db, err)
        db.Close()
        return nil, errors.New("fail to open db")
    }
    RocksdbClient = db
    return db, nil
}
