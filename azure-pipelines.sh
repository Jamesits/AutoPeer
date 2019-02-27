#!/bin/bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive
export GOPATH=/tmp/go
export GOBIN=/tmp/go/bin
export EXE=autopeer

function install_deps() {
    echo "Installing dependencies"

    sudo apt-get update
    sudo apt-get install -y upx
}

function build() {
    echo "Building for OS=$1 ARCH=$2"

    env GOOS="$1" GOARCH="$2" go build -ldflags="-s -w" -o ${BUILD_ARTIFACTSTAGINGDIRECTORY}/${EXE}-"$3"
    ! upx --ultra-brute ${BUILD_ARTIFACTSTAGINGDIRECTORY}/${EXE}-"$3" || true
}

function test_binary() {
    echo "Testing binary"

    BINARY=${BUILD_ARTIFACTSTAGINGDIRECTORY}/${EXE}-linux-amd64
    ${BINARY}
}

install_deps
go get ./...
build linux amd64 linux-amd64
# build windows amd64 windows-amd64.exe
#test_binary
