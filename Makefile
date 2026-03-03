.DEFAULT_TARGET := build

REVISION := $(shell cat git-revision 2>/dev/null || git rev-parse HEAD)
RELEASE_VERSION := $(shell cat version)

build:
	go build -ldflags "-X main.commit=$(REVISION) -X main.version=$(RELEASE_VERSION)" -o fhcli cmd/main.go
