#!/bin/bash

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@v2.5.0
go install github.com/air-verse/air@v1.52.3
