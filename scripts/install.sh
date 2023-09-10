#!/bin/bash

set -eux

TAG=$(git describe --tags)
go install -ldflags "-X github.com/cbguder/books/cmd.Version=${TAG}"
