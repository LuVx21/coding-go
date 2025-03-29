#!/bin/bash

if [[ -z $(which protoc) ]]; then
  brew install protobuf
fi

if [[ -z $(which protoc-gen-go) ]]; then
  go install -ldflags="-w -s" google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if [[ -z $(which protoc-gen-go-grpc) ]]; then
  go install -ldflags="-w -s" google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/helloworld.proto
