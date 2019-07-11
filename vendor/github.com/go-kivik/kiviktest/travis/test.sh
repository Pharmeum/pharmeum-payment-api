#!/bin/bash
set -e

if [ "$TRAVIS_GO_VERSION" == "1.11.x" ]; then
    go get -u gopkg.in/alecthomas/gometalinter.v1 && gometalinter.v1 --install
    gometalinter.v1 --config .linter.json ./...
fi

go test $(go list ./... | grep -v /vendor/)
