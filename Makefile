NAME ?= riff
OUTPUT = ./bin/$(NAME)
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin
VERSION ?= $(shell cat VERSION)
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")
LDFLAGS_VERSION = -X github.com/projectriff/cli/pkg/cli.cli_name=$(NAME) \
				  -X github.com/projectriff/cli/pkg/cli.cli_version=$(VERSION) \
				  -X github.com/projectriff/cli/pkg/cli.cli_gitsha=$(GITSHA) \
				  -X github.com/projectriff/cli/pkg/cli.cli_gitdirty=$(GITDIRTY)

MOCKERY ?= go run -modfile hack/go.mod github.com/vektra/mockery/cmd/mockery
GOIMPORTS ?= go run -modfile hack/go.mod golang.org/x/tools/cmd/goimports

.PHONY: all
all: build test verify-goimports docs ## Build, test, verify source formatting and regenerate docs

.PHONY: clean
clean: ## Delete build output
	rm -rf bin/
	rm -rf dist/

.PHONY: build
build: $(OUTPUT) ## Build riff

.PHONY: test
test: ## Run the tests
	go test ./...

.PHONY: install
install: build ## Copy build to GOPATH/bin
	cp $(OUTPUT) $(GOBIN)

.PHONY: coverage
coverage: ## Run the tests with coverage and race detection
	go test -v --race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: goimports
goimports: ## Runs goimports on the project
	@$(GOIMPORTS) -w pkg cmd

.PHONY: verify-goimports
verify-goimports: ## Verifies if all source files are formatted correctly
	@$(GOIMPORTS) -l pkg cmd | (! grep .) || (echo above files are not formatted correctly. please run \"make goimports\" && false)

$(OUTPUT): $(GO_SOURCES) VERSION
	go build -o $(OUTPUT) -ldflags "$(LDFLAGS_VERSION)" ./cmd/riff

.PHONY: release
release: $(GO_SOURCES) VERSION ## Cross-compile riff for various operating systems
	@mkdir -p dist
	GOOS=darwin   GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/riff && tar -czf dist/$(NAME)-darwin-amd64.tgz  -C bin . && rm -f $(OUTPUT)
	GOOS=linux    GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/riff && tar -czf dist/$(NAME)-linux-amd64.tgz   -C bin . && rm -f $(OUTPUT)
	GOOS=windows  GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT).exe ./cmd/riff && zip -rj  dist/$(NAME)-windows-amd64.zip    bin   && rm -f $(OUTPUT).exe

docs: $(OUTPUT) clean-docs ## Generate documentation
	$(OUTPUT) docs

.PHONY: verify-docs
verify-docs: docs ## Verify the generated docs are up to date
	git diff --exit-code docs

.PHONY: clean-docs
clean-docs: ## Delete the generated docs
	rm -fR docs

.PHONY: gen-mocks
gen-mocks: clean-mocks ## Generate mocks
	@$(MOCKERY) -output ./pkg/testing/pack -outpkg pack -dir ./pkg/pack -name Client
	@$(MOCKERY) -output ./pkg/testing/kail -outpkg kail -dir ./pkg/kail -name Logger
	@make goimports

.PHONY: clean-mocks
clean-mocks: ## Delete mocks
	@rm -fR pkg/testing/pack
	@rm -fR pkg/testing/kail

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
