#!/bin/zsh

#go get github.com/tecbot/gorocksdb

version=8.11.3
export CGO_CFLAGS="-I/opt/homebrew/Cellar/rocksdb/$version/include"
export CGO_LDFLAGS="-L/opt/homebrew/Cellar/rocksdb/$version/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -L/opt/homebrew/Cellar/snappy/1.1.10/lib -L/opt/homebrew/Cellar/lz4/1.9.4/lib -L/opt/homebrew/Cellar/zstd/1.5.6/lib"

go test -v rocksdb_test.go rocksdb.go
