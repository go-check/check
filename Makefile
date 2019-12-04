# Include go binaries into path

VERSION=`git rev-parse --short HEAD`

GODEBUG:=GODEBUG=gocacheverify=1
GOBIN := $(PWD)/bin/
ENV:=GOBIN=$(GOBIN)

# Defaults...
all: mod tests build

# Run tests
tests: test

test:
	@echo "======================================================================"
	@go clean -cache
	@go clean -testcache
	go test --check.format=teamcity ./

deps:
	@echo "======================================================================"
	@echo 'MAKE: deps...'
	@mkdir -p $(GOBIN)
	@$(ENV) go get golang.org/x/lint/golint

mod:
	@echo "======================================================================"
	@echo "Run MOD..."
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod tidy
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod vendor
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod download
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify
