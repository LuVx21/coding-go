#!/bin/bash

#brew install protobuf
#go install -ldflags="-w -s" google.golang.org/protobuf/cmd/protoc-gen-go@latest
#go install -ldflags="-w -s" google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/helloworld.proto
