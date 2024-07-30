#!/bin/bash

MODE="linux"
#MODE="windows"

export GOARCH="amd64"
ROOT_DIR=`pwd`

case ${MODE} in
linux)
  export GOOS="linux"
  NAME="tool-server"
  ;;
windows)
  export GOOS="windows"
  NAME="tool-server.exe"
  ;;
esac

GitCommitHash=`git rev-parse HEAD`
Version="0.0.2"
BuildTime=$(date +%s%3N)
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'tool-server/internal/global.GitCommitHash=${GitCommitHash}' \
    -X 'tool-server/internal/global.Version=${Version}' \
    -X 'tool-server/internal/global.tBuildTime=${BuildTime}' \
    -X 'tool-server/internal/global.BuildGoVersion=${BuildGoVersion}' \
"

echo "Build start"

go build -ldflags "$LDFlags" -o ${ROOT_DIR}/${NAME}

echo "Build end"