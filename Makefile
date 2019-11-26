BINARY_NAME=cart-api
GOLANGCI_LINT=bin/golangci-lint
GOLANGCI_VER=v1.19.0

HAS_LINTER := $(shell command -v golangci-lint;)

.PHONY: build
build:
	GO111MODULE=on go build -o $(BINARY_NAME) main.go

.PHONY: bootstrap
bootstrap:
ifndef HAS_LINTER
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

.PHONY: test
test:
	go test -race ./... -v -coverprofile .testCoverage.test

.PHONY: lint
lint:
	if [ ! -f $(GOLANGCI_LINT) ]; \
	then \
		echo "golangci-lint not found, installing"; \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s $(GOLANGCI_VER); \
	else \
		echo "found golangci-lint"; \
	fi
	./bin/golangci-lint -v run

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

.PHONY: generate
generate:
	GO111MODULE=on go generate ./...
