

# Go artifacts
GO_ARTIFACTS := index-server

GO_BUILD_PARAMS := -ldflags "-s -w"
GO_TEST_PACKAGES ?= ./...

all: build

# =============================================================================
# Build
# =============================================================================
build: build-go ## Build all artifacts

build-go: $(GO_ARTIFACTS) ## Build Go artifacts

# =============================================================================
# Artifacts
# =============================================================================
$(GO_ARTIFACTS):
	@echo "Building Go artifact $@"
	go build $(GO_BUILD_PARAMS) ./cmd/$@

# =============================================================================
# Tests
# =============================================================================
test: unit-test-go ## Execute all tests

unit-test-go: ## Execute go unit tests
	@echo "Running go unit tests"
	@go clean -testcache
	@go test -p 1 $(GO_BUILD_PARAMS) $(GO_TEST_PACKAGES)