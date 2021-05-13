#!/usr/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -v -ldflags="-s -w -X main.VERSION=1.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" -o nine ../cmd/main.go