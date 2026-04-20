#!/bin/bash

if [[ -z $(which wire) ]]; then
    go install -ldflags="-w -s" github.com/google/wire/cmd/wire@latest
fi

wire

go run main.go provider.go wire_gen.go
