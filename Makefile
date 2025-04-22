COVERAGE_FILE ?= coverage.out

.PHONY: build
build: build_zoo

.PHONY: build_zoo
build_zoo:
	@echo "Running build"
	@mkdir -p .bin
	@go build -o ./bin/ddd_zoo ./cmd/ddd_zoo

## test: run all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/maklybae/ddd-zoo/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'

.PHONY: lint
lint: lint-golang

.PHONY: lint-golang
lint-golang:
	@if ! command -v 'golangci-lint' &> /dev/null; then \
  		echo "Please install golangci-lint!"; exit 1; \
  	fi;
	@golangci-lint -v run --fix ./...

.PHONY: generate
generate: generate_openapi

.PHONY: generate_openapi
generate_openapi:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;
	@mkdir -p internal/types/openapi/v1
	@oapi-codegen -package v1 \
		-generate types,gin \
		api/openapi/v1/ddd_zoo.yaml > internal/types/openapi/v1/ddd_zoo.gen.go

.PHONY: clean
clean:
	@rm -rf./bin