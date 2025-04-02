#!/bin/zsh

#go get github.com/tecbot/gorocksdb

LIB_PATH="$(brew --prefix)/Cellar"
versions=$(brew ls --versions rocksdb snappy lz4 zstd)

rocksdb_home=${LIB_PATH}/rocksdb/$(echo $versions | grep rocksdb | cut -d' ' -f2)
snappy_home=${LIB_PATH}/snappy/$(echo $versions | grep snappy | cut -d' ' -f2)
lz4_home=${LIB_PATH}/lz4/$(echo $versions | grep lz4 | cut -d' ' -f2)
zstd_home=${LIB_PATH}/zstd/$(echo $versions | grep zstd | cut -d' ' -f2)

export CGO_CFLAGS="-I${rocksdb_home}/include"
export CGO_LDFLAGS="-L${rocksdb_home}/lib -L${snappy_home}/lib -L${lz4_home}/lib -L${zstd_home}/lib"
# -lrocksdb -lstdc++ -lm -lz -lbz2

go test -v rocksdb_test.go rocksdb.go
