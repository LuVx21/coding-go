#!/bin/bash

#go install -ldflags="-w -s" google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#go install -ldflags="-w -s" github.com/go-micro/generator/cmd/protoc-gen-micro@latest

protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --micro_out=. \
  --micro_opt=paths=source_relative \
  proto/helloworld.proto
