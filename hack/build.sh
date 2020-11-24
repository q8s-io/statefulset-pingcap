#!/bin/bash

go mod vendor

OS=`go env GOOS`
ARCH=`go env GOARCH`

LDFLAGS="-X vendor/k8s.io/component-base/version.gitVersion=`git describe --long --dirty --abbrev=14`"

GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o controller-manager github.com/q8s-io/statefulset-pingcap/cmd/controller-manager
