#!/bin/zsh

#go get github.com/tecbot/gorocksdb

LIB_PATH="/opt/homebrew/Cellar"

rocksdb_version=9.10.0
snappy_version=1.2.1
lz4_version=1.10.0
zstd_version=1.5.6

export CGO_CFLAGS="-I${LIB_PATH}/rocksdb/${rocksdb_version}/include"
export CGO_LDFLAGS="-L${LIB_PATH}/rocksdb/${rocksdb_version}/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -L${LIB_PATH}/snappy/${snappy_version}/lib -L${LIB_PATH}/lz4/${lz4_version}/lib -L${LIB_PATH}/zstd/${zstd_version}/lib"

go test -v rocksdb_test.go rocksdb.go
