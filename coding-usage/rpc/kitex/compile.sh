#!/bin/bash

if [[ -z $(which kitex) ]]; then
    go install -ldflags="-w -s" github.com/cloudwego/kitex/tool/cmd/kitex@latest
fi


kitex -module github.com/luvx21/coding-go/coding-usage ./proto/helloworld.proto
