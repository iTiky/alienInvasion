# Version
GIT_TAG ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date +%FT%T%z)

# Go build ldflags
LDFLAG_TAG := -X 'github.com/itiky/alieninvasion/cmd/alieninvasion.GitTag=$(GIT_TAG)'
LDFLAG_COMMIT := -X 'github.com/itiky/alieninvasion/cmd/alieninvasion.GitCommit=$(GIT_COMMIT)'
LDFLAG_BUILD := -X 'github.com/itiky/alieninvasion/cmd/alieninvasion.BuildDate=$(BUILD_DATE)'

# Output names
BINARY_NAME := ai

.PHONY: deps
deps:
	@echo "Downloading go.mod dependencies"
	go mod download

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build-binary: deps
	@echo "Building binary ($(BINARY_NAME)): $(GIT_TAG).$(GIT_COMMIT).$(BUILD_DATE)"
	@go build -ldflags "$(LDFLAG_TAG) $(LDFLAG_COMMIT) $(LDFLAG_BUILD)" -o $(BINARY_NAME) ./cmd/main.go
