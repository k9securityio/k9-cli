VERSION := $(shell git describe --abbrev=0 --tags)
REV:= $(shell git rev-parse --short HEAD)
BUILD_NUMBER := "${BUILD_NUMBER}"
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%S%z")
PWD := $(shell pwd)

setup:
	# Add installation of developer tools here
	# * golangci-lint
	# * upx
	# * goimports

lint:
	@golangci-lint run \
		-D errcheck -D deadcode -D varcheck -D unused \
		-E gosec -E dupl -E goconst -E misspell -E lll -E unparam -E gochecknoinits
build:
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
		go build \
		-ldflags "-s -w -X 'github.com/k9securityio/k9-cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9-cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9-cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-darwinM1
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build \
		-ldflags "-s -w -X 'github.com/k9securityio/k9-cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9-cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9-cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-darwin64
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0\
		go build \
		-ldflags "-s -w -X 'github.com/k9securityio/k9-cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9-cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9-cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-linux64
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0\
		go build \
		-ldflags "-X 'github.com/k9securityio/k9-cli/cmd.version=$(VERSION)' -X 'github.com/k9securityio/k9-cli/cmd.revision=$(REV)' -X 'github.com/k9securityio/k9-cli/cmd.buildtime=$(BUILD_TIME)'" \
		-o ./bin/k9-windows64.exe

package:
	@docker build -t k9securityio/k9:$(VERSION) -t k9securityio/k9:$(REV) -t k9securityio/k9:b$(BUILD_NUMBER) .

push:
	@docker push k9securityio/k9:$(VERSION)
	@docker push k9securityio/k9:$(REV)
	@docker push k9securityio/k9:b$(BUILD_NUMBER)
