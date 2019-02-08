.DEFAULT_TARGET := build

REVISION := $(shell cat git-revision 2>/dev/null || git rev-parse HEAD)
RELEASE_VERSION := $(shell cat version)

build:
	go build -ldflags "-X github.com/firehydrant/fhcli/cmd.cliRevision=$(REVISION) -X github.com/firehydrant/fhcli/cmd.releaseVersion=$(RELEASE_VERSION)" -o fhcli *.go
