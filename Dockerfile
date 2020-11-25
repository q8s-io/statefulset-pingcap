# build builder
FROM golang:1.14 as builder

WORKDIR /go/src/github.com/q8s-io/statefulset-pingcap

COPY . .

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOPROXY=https://mirrors.aliyun.com/goproxy/  go build ./tmp/statefulset-pingcap-controller-manager github.com/q8s-io/statefulset-pingcap/cmd/controller-manager

# build server
FROM alpine:3.8

WORKDIR /

COPY --from=builder /go/src/github.com/q8s-io/statefulset-pingcap/tmp/statefulset-pingcap-controller-manager /usr/local/bin/statefulset-pingcap-controller-manager

ENTRYPOINT ["/usr/local/bin/statefulset-pingcap-controller-manager"]
