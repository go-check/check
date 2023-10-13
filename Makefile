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
	go test -cover -coverprofile=./check_coverage.out ./
	go test -cover -coverprofile=./check_test_coverage.out -coverpkg  check_test ./
	go tool cover -html=check_coverage.out -o check_coverage.html
	go tool cover -html=check_test_coverage.out -o check_test_coverage.html
	rm -f ./check_coverage.out
	rm -f ./check_test_coverage.out

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


clean-cache: ## Clean golang cache
	@echo "clean-cache started..."
	go clean -cache
	go clean -testcache
	@echo "clean-cache complete!"

clean-vendor: ## Remove vendor folder
	@echo "clean-vendor started..."
	rm -fr ./vendor
	@echo "clean-vendor complete!"

# full cleaning. Don't use it: it removes outside golang packages for all projects
clean: clean clean-cache clean-vendor ## Remove all build artifacts
	@echo "======================================================================"
	@echo "Run clean"
	rm -f ./application
	rm -f ./application_consumer
	go clean -i -r -x -cache -testcache -modcache -fuzzcache