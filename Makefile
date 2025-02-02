BUILD_DIR = builds
BINARY_NAME = alert-forge
MODULE = github.com/soerenschneider/$(BINARY_NAME)
CHECKSUM_FILE = checksum.sha256
SIGNATURE_KEYFILE = ~/.signify/github.sec
DOCKER_PREFIX = ghcr.io/soerenschneider

version-info:
	$(eval VERSION := $(shell git describe --tags --abbrev=0 || echo "dev"))
	$(eval COMMIT_HASH := $(shell git rev-parse HEAD))


generate:
	go generate  ./...

build: version-info
	CGO_ENABLED=1 go build -ldflags="-w -X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BINARY_NAME) ./cmd

tests:
	go test ./...
