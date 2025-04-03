#!/bin/bash

if [[ -z $(which protoc) ]]; then
  brew install protobuf
fi

if [[ -z $(which protoc-gen-go) ]]; then
  go install -ldflags="-w -s" google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if [[ -z $(which protoc-gen-go-triple) ]]; then
  go install -ldflags="-w -s" github.com/dubbogo/protoc-gen-go-triple/v3@latest
fi


protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-triple_out=. \
  --go-triple_opt=paths=source_relative \
  proto/helloworld.proto
