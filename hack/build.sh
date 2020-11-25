#!/bin/bash

OS=`go env GOOS`
ARCH=`go env GOARCH`

GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 go build -o ./tmp/statefulset-pingcap-controller-manager github.com/q8s-io/statefulset-pingcap/cmd/controller-manager
