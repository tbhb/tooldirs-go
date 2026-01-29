set dotenv-load := true
set positional-arguments := true

export PATH := `go env GOPATH` + "/bin:" + env("PATH")

# Default recipe
default: lint test

# Install all development dependencies
install: install-tools
    pnpm install

# Install Go tools
install-tools:
    go install github.com/daixiang0/gci@latest
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
    go install github.com/rhysd/actionlint/cmd/actionlint@latest
    go install golang.org/x/pkgsite/cmd/pkgsite@latest
    go install golang.org/x/tools/cmd/goimports@latest
    go install mvdan.cc/gofumpt@latest

# Run all linters
lint: lint-go lint-markdown lint-prose lint-web lint-actions

# Run Go linters (golangci-lint)
lint-go:
    golangci-lint run

lint-markdown:
    pnpm exec markdownlint-cli2 **/*.md

# Run prose linter (vale)
lint-prose:
    vale .

# Run web linters (biome)
lint-web:
    pnpm exec biome check .

# Run GitHub Actions linter
lint-actions:
    actionlint

# Alias for lint-actions
actionlint: lint-actions

# Run all formatters
fmt: fmt-go fmt-web

# Format Go code (uses golangci-lint formatters)
fmt-go:
    golangci-lint fmt

# Format web files (biome)
fmt-web:
    pnpm exec biome check --write --files-ignore-unknown=true --no-errors-on-unmatched .

# Fix all linting issues (where possible)
fix: fix-go fix-web

# Fix Go linting issues
fix-go:
    golangci-lint fmt
    golangci-lint run --fix

# Fix web linting issues
fix-web:
    pnpm exec biome check --write --files-ignore-unknown=true --no-errors-on-unmatched .

# Run tests
test *args:
    go test ./... "$@"

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -coverprofile=coverage.out -covermode=atomic ./...

# Run tests with coverage and generate HTML report
test-coverage-html: test-coverage
    go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
test-race:
    go test -race ./...

# Run benchmarks
bench *args:
    go test -bench=. -benchmem ./... "$@"

# Generate code (if any)
generate:
    go generate ./...

# Start local documentation server
docs:
    pkgsite -http localhost:6060

# Tidy go.mod
tidy:
    go mod tidy

# Verify dependencies
verify:
    go mod verify

# Update dependencies
update:
    go get -u ./...
    go mod tidy

# Check that everything is ready for commit
check: tidy verify lint test

# Install pre-commit hooks (using prek)
hooks-install:
    prek install

# Run pre-commit hooks on all files
hooks-run:
    prek run --all-files

# Run pre-commit hooks on staged files only
hooks-staged:
    prek run

# Update pre-commit hook versions
hooks-update:
    prek autoupdate
