#!/bin/bash

MODE="linux"
#MODE="windows"

export GOARCH="amd64"
ROOT_DIR=`pwd`

case ${MODE} in
linux)
  export GOOS="linux"
  NAME="ai-software-copyright-server"
  ;;
windows)
  export GOOS="windows"
  NAME="ai-software-copyright-server.exe"
  ;;
esac

GitCommitHash=`git rev-parse HEAD`
Version="0.0.2"
BuildTime=$(date +%s%3N)
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'ai-software-copyright-server/internal/global.GitCommitHash=${GitCommitHash}' \
    -X 'ai-software-copyright-server/internal/global.Version=${Version}' \
    -X 'ai-software-copyright-server/internal/global.tBuildTime=${BuildTime}' \
    -X 'ai-software-copyright-server/internal/global.BuildGoVersion=${BuildGoVersion}' \
    -X 'ai-software-copyright-server/internal/global.Host=https://rz.nineya.com' \
"

echo "Build start"

rm -f ${NAME}
rm -f ${NAME}-bak

go build -ldflags "$LDFlags" -o ${ROOT_DIR}/${NAME}

cp ${NAME} ${NAME}-bak

echo "Build end"