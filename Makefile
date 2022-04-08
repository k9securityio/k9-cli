VERSION := $(shell git describe --abbrev=0 --tags)
REV:= $(shell git rev-parse --short HEAD)
BUILD_NUMBER := "${BUILD_NUMBER}"
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%S%z")
PWD := $(shell pwd)

lint:
	@golangci-lint run \
		-D errcheck -D deadcode -D varcheck -D unused \
		-E gosec -E dupl -E goconst -E misspell -E lll -E unparam -E gochecknoinits
build:
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build \
		-ldflags "-X 'github.com/k9securityio/k9cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-darwinM1
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build \
		-ldflags "-X 'github.com/k9securityio/k9cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-darwin64
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0\
		go build \
		-ldflags "-X 'github.com/k9securityio/k9cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-linux64

package:
	@docker build -t k9securityio/k9:$(VERSION) -t k9securityio/k9:$(REV) -t k9securityio/k9:b$(BUILD_NUMBER) .

push:
	@docker push k9securityio/k9:$(VERSION)
	@docker push k9securityio/k9:$(REV)
	@docker push k9securityio/k9:b$(BUILD_NUMBER)
