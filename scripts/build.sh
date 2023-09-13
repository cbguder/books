#!/bin/bash

set -eux

TAG=$(git describe --tags)

mkdir -p release
pushd release

build() {
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -ldflags "-X github.com/cbguder/books/cmd.version=${TAG}" ..
  tar -czf "books-${TAG}-$1-$2.tar.gz" books
  rm -f books
}

build darwin amd64
build darwin arm64
build linux amd64
build linux arm64

popd
