#!/bin/bash

#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
#go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --micro_out=. \
  --micro_opt=paths=source_relative \
  proto/helloworld.proto
