GO  := go

# Enable GO111MODULE=off explicitly, enable it with GO111MODULE=on when necessary.
export GO111MODULE := off

ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

ALL_TARGETS := cmd/controller-manager
SRC_PREFIX := github.com/q8s-io/statefulset-pingcap
GIT_VERSION = $(shell ./hack/version.sh | awk -F': ' '/^GIT_VERSION:/ {print $$2}')

# in GOPATH mode, we must use the full path name related to $GOPATH.
# https://github.com/golang/go/issues/19000
ifneq ($(VERSION),)
    LDFLAGS += -X $(SRC_PREFIX)/vendor/k8s.io/component-base/version.gitVersion=${VERSION}
else
    LDFLAGS += -X $(SRC_PREFIX)/vendor/k8s.io/component-base/version.gitVersion=${GIT_VERSION}
endif

all: build
.PHONY: all

build: $(ALL_TARGETS)
.PHONY: all

$(ALL_TARGETS):
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 $(GO) build -ldflags "${LDFLAGS}" -o output/bin/$(OS)/$(ARCH)/$@ $(SRC_PREFIX)/$@
.PHONY: $(ALL_TARGETS)
