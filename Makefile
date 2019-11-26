BINARY_NAME=cart-api
COMPOSE_FILE=Docker-compose.yml
GOLANGCI_LINT=bin/golangci-lint
GOLANGCI_VER=v1.19.0

RELEASE?=1.0.0
GOOS?=linux
GOARCH?=amd64

HAS_LINTER := $(shell command -v golangci-lint;)

.PHONY: buildStatic
buildStatic:
	GOPRIVATE=gitlab.itechart-group.com GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -ldflags "-linkmode external -extldflags -static" -o $(BINARY_NAME) cmd/telegram-bot/main.go

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

.PHONY: compose-run
compose-run:
	docker-compose -f $(COMPOSE_FILE) up -d --force-recreate
