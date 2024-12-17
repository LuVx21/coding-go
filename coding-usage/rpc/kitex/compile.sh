#!/bin/bash

#go install -ldflags="-w -s" github.com/cloudwego/kitex/tool/cmd/kitex@latest

kitex -module github.com/luvx21/coding-go/coding-usage ./proto/helloworld.proto
