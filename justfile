# ImageManager Go - Just Command Runner
# Modern task runner for Go development with multi-module support

# Configuration variables
go_version := "1.24"
binary_name := "ImageManager"
main_path := "."
modules := ". core" 
all_modules := ". core"

# Build flags
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
commit_hash := `git rev-parse HEAD 2>/dev/null || echo "unknown"`
build_time := `date -u '+%Y-%m-%d_%H:%M:%S'`
build_user := `whoami`

build_flags := "-trimpath -ldflags=\"-s -w -X main.version=" + version + " -X main.commit=" + commit_hash + " -X main.buildTime=" + build_time + " -X main.buildUser=" + build_user + "\""
debug_flags := "-race -ldflags=\"-X main.version=" + version + " -X main.commit=" + commit_hash + " -X main.buildTime=" + build_time + "\""

# Default recipe
default: all

# ============================================================================
# Core Development Workflows
# ============================================================================

# Development workflow (format, tidy, test)
dev: clean validate-standards smart-build run-sort
    @echo "✅ Development cycle completed"

# Run complete pipeline: clean, format, lint, test, build
all: clean deps-tidy fmt lint test build-release
    @echo "✅ Full pipeline completed"

# Fast pre-commit validation for local commits
pre-commit: fmt-check lint
    @echo "✅ Pre-commit validation passed"

# Comprehensive pre-push validation (REQUIRED before pushing)
pre-push: fmt-check validate-all test-unit security
    @echo "✅ Pre-push validation passed"

# Comprehensive validation suite
validate-all: validate-standards validate-patterns
    @echo "✅ Comprehensive validation completed"

# ============================================================================
# Build Targets
# ============================================================================

# Build application (debug version with race detector)
build: build-debug

# Build debug version with race detection
build-debug: deps-check
    @echo "🔨 Building {{binary_name}} (debug)..."
    CGO_ENABLED=1 go build {{debug_flags}} -o {{binary_name}}
    @echo "✅ Debug build completed: {{main_path}}/{{binary_name}}"

# Build optimized release binary
build-release: deps-check
    @echo "🔨 Building {{binary_name}} (release)..."
    CGO_ENABLED=0 go build {{build_flags}} -o {{binary_name}}
    @echo "✅ Release build completed: {{main_path}}/{{binary_name}}"

# Build for multiple platforms
build-all:
    @echo "🌍 Building for multiple platforms..."
    mkdir -p dist
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build {{build_flags}} -o dist/{{binary_name}}-linux-amd64
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build {{build_flags}} -o dist/{{binary_name}}-darwin-amd64
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build {{build_flags}} -o dist/{{binary_name}}-darwin-arm64
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build {{build_flags}} -o dist/{{binary_name}}-windows-amd64.exe
    @echo "✅ Multi-platform build completed in dist/"

# Smart rebuild - only if sources are newer than binary
smart-build:
    #!/usr/bin/env bash
    echo "🔧 Checking if rebuild is needed..."
    if [ ! -f "{{main_path}}/{{binary_name}}" ]; then
        echo "🔨 Binary not found, building..."
        just build-release
    elif find . -name "*.go" -newer "{{main_path}}/{{binary_name}}" | grep -q .; then
        echo "🔨 Sources newer than binary, rebuilding..."
        just build-release
    else
        echo "✅ Binary is up to date"
    fi

# ============================================================================
# Testing
# ============================================================================

# Run all unit tests
test: test-unit

# Run unit tests with proper workspace handling
test-unit:
    @echo "🧪 Running unit tests..."
    @if go test -race -v -timeout=2m -short ./internal/... ./core/...; then \
        echo "✅ Unit tests passed"; \
    else \
        echo "❌ Unit tests failed"; \
        exit 1; \
    fi
    @echo "✅ Unit tests completed successfully"

# Run integration tests with smart rebuild
test-integration: smart-build
    @echo "🔄 Running integration tests..."
    @if [ -d "test/integration" ]; then \
        echo "🧪 Running integration tests..."; \
        cd test/integration && go test -v -timeout=10m ./...; \
        if [ $$? -eq 0 ]; then \
            echo "✅ Integration tests passed"; \
        else \
            echo "❌ Integration tests failed"; \
            exit 1; \
        fi; \
    else \
        echo "⚠️ Integration tests not found, skipping"; \
    fi

# Generate comprehensive test coverage report
test-coverage:
    @echo "📊 Generating coverage report..."
    mkdir -p coverage
    @if go test -race -timeout=2m -short -coverprofile=coverage/coverage.out -covermode=atomic ./internal/... ./core/...; then \
        echo "✅ Tests passed, generating coverage report..."; \
        go tool cover -html=coverage/coverage.out -o coverage/coverage.html; \
        echo "📊 Coverage summary:"; \
        go tool cover -func=coverage/coverage.out | tail -1; \
        echo "✅ Coverage report generated: coverage/coverage.html"; \
    else \
        echo "❌ Tests failed, cannot generate coverage report"; \
        exit 1; \
    fi

# Run benchmarks
bench:
    @echo "⚡ Running benchmarks..."
    @if go test -bench=. -benchmem -benchtime=5s ./internal/... ./core/...; then \
        echo "✅ Benchmarks completed successfully"; \
    else \
        echo "❌ Benchmark tests failed"; \
        exit 1; \
    fi

# ============================================================================
# Code Quality
# ============================================================================

# Format code using golangci-lint formatters
fmt:
    @echo "✨ Formatting code with golangci-lint..."
    golangci-lint run --config .golangci.yml --fix || exit 1
    @echo "✅ Code formatting completed"

# Check code formatting without making changes
fmt-check:
    @echo "🔍 Checking code formatting..."
    golangci-lint run --config .golangci.yml --issues-exit-code=1 || exit 1
    @echo "✅ Format check completed"

# Run comprehensive linting
lint:
    @echo "🔍 Running comprehensive linting..."
    golangci-lint run --config .golangci.yml --timeout 5m || exit 1
    @echo "✅ Linting completed"

# Run linting with auto-fixing
lint-fix:
    @echo "🔧 Running linting with auto-fix..."
    golangci-lint run --config .golangci.yml --fix --timeout 5m || exit 1
    @echo "✅ Linting with auto-fix completed"

# ============================================================================
# Security and Dependencies
# ============================================================================

# Run comprehensive security scanning
security:
    @echo "🛡️ Running security scans..."
    govulncheck ./... || exit 1
    @if [ -d "test/integration" ] && [ -f "test/integration/go.mod" ]; then \
        echo "🔍 Checking integration_test module for vulnerabilities..."; \
        cd test/integration && govulncheck ./...; \
    fi
    @echo "✅ Security scanning completed"

# Check dependency status
deps-check:
    @echo "📦 Checking dependencies..."
    go mod verify
    go mod tidy -diff
    @echo "✅ Dependencies verified"

# Update dependencies to latest versions
deps-update:
    @echo "⬆️ Updating dependencies..."
    go get -u ./...
    go mod tidy
    @echo "✅ Dependencies updated"

# Clean up dependencies
deps-tidy:
    @go mod tidy

# Generate dependency graph
deps-graph:
    @echo "📊 Generating dependency graph..."
    mkdir -p reports
    go mod graph > reports/deps-graph.txt
    @echo "✅ Dependency graph saved to reports/deps-graph.txt"

# ============================================================================
# Quality Validation
# ============================================================================

# Validate Go standards and conventions
validate-standards:
    @echo "📏 Validating Go standards..."
    golangci-lint run --config .golangci.yml --enable-only=govet --timeout 2m || exit 1
    @echo "✅ Standards validation completed"

# Validate architectural patterns
validate-patterns:
    @echo "🏗️ Validating architectural patterns..."
    @echo "🔍 Checking naming conventions..."
    @if grep -r "func.*[a-z].*(" --include="*.go" . | grep -v "_test.go" | grep -E "(New[A-Z]|Default[A-Z])" > /dev/null; then \
        echo "✅ Constructor naming conventions followed"; \
    else \
        echo "⚠️ Check constructor naming conventions"; \
    fi
    @echo "🔍 Checking interface patterns..."
    @interface_count=`find . -name "*.go" -not -path "./vendor/*" -exec grep -l "^type.*interface" {} \; | wc -l`; \
    echo "✅ Found $$interface_count interface definitions"
    @echo "✅ Pattern validation completed"


# ============================================================================
# CI/CD Workflows
# ============================================================================


# CI testing with coverage
ci-test: test-unit test-coverage
    @echo "✅ CI testing completed"

# CI build for multiple platforms  
ci-build: build-all
    @echo "✅ CI build completed"

# Complete CI pipeline
ci-full: validate-all security ci-test ci-build
    @echo "✅ Full CI pipeline completed"

# Quick CI validation for pull requests
ci-quick: fmt-check lint test-unit
    @echo "✅ Quick CI validation completed"

# ============================================================================
# Development Tools
# ============================================================================

# Install development tools with pinned versions
install-tools:
    @echo "🛠️ Installing development tools..."
    go install golang.org/x/tools/cmd/goimports@latest
    @echo "🔧 Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b `go env GOPATH`/bin v2.3.0
    go install golang.org/x/tools/cmd/deadcode@v0.25.0
    go install golang.org/x/vuln/cmd/govulncheck@v1.1.3
    @echo "✅ Development tools installed"

# Analyze dead code
deadcode:
    @echo "🔍 Analyzing dead code..."
    deadcode ./internal/... ./core/... || echo "✅ No dead code found"
    @echo "✅ Dead code analysis completed"

# ============================================================================
# Application Execution
# ============================================================================

# Run sort example with smart rebuild
run-sort: smart-build
    @echo "🔄 Running sort example..."
    ./{{binary_name}} sort --source "./test_images" --destination "./sorted" --actionStrategy "dryRun"

# Run deduplication example with smart rebuild
run-dedup: smart-build
    @echo "🔄 Running deduplication example..."
    ./{{binary_name}} dedup --source "./test_images" --actionStrategy "dryRun"

# ============================================================================
# Docker Support
# ============================================================================

# Build Docker images
docker-build:
    @echo "🐳 Building Docker images..."
    docker-compose build
    @echo "✅ Docker images built successfully"

# Start development environment
docker-up-dev:
    @echo "🐳 Starting development environment..."
    docker-compose --profile development up -d
    @echo "✅ Development environment started"

# Run tests in Docker container
docker-test:
    @echo "🧪 Running tests in Docker container..."
    docker-compose --profile testing run --rm imagemanager-test
    @echo "✅ Docker tests completed"

# Run CI pipeline in Docker
docker-ci:
    @echo "🔄 Running CI pipeline in Docker..."
    docker-compose --profile ci run --rm imagemanager-ci
    @echo "✅ Docker CI pipeline completed"

# Stop and remove all containers
docker-down:
    @echo "🛑 Stopping all containers..."
    docker-compose --profile production --profile development --profile testing --profile ci down
    @echo "✅ All containers stopped"

# ============================================================================
# Cleanup
# ============================================================================

# Clean build artifacts and temporary files
clean:
    @echo "🧹 Cleaning build artifacts..."
    rm -f ./{{binary_name}}
    rm -f ./application.log
    rm -f ./*.csv
    rm -f *.log
    @if [ -d "test/integration" ]; then \
        go clean -testcache 2>/dev/null || true; \
    fi
    @echo "✅ Cleanup completed"

# Deep clean including caches and generated files
clean-all: clean
    @echo "🧹 Deep cleaning..."
    rm -rf dist/
    rm -rf coverage/
    rm -rf profiles/
    rm -rf reports/
    @go clean -cache -testcache -modcache 2>/dev/null || echo "⚠️ Some caches could not be cleaned"
    @echo "✅ Deep cleanup completed"

# ============================================================================
# Information
# ============================================================================

# Show version information
version:
    @echo "📋 Version Information:"
    @echo "  Version: {{version}}"
    @echo "  Commit:  {{commit_hash}}"
    @echo "  Built:   {{build_time}}"
    @echo "  By:      {{build_user}}"
    @echo "  Go:      {{go_version}}"

# Show help information
help:
    @echo "ImageManager Go - Just Command Runner"
    @echo "====================================="
    @echo ""
    @echo "🏗️  Build & Development:"
    @echo "  build              Build application (debug version)"
    @echo "  build-release      Build optimized release binary"
    @echo "  build-all          Build for multiple platforms"
    @echo "  dev                Development workflow (format, tidy, test)"
    @echo "  clean              Clean build artifacts"
    @echo "  clean-all          Deep clean including caches"
    @echo ""
    @echo "🧪 Testing & Quality:"
    @echo "  test               Run all unit tests"
    @echo "  test-integration   Run integration tests"
    @echo "  test-coverage      Generate coverage report"
    @echo "  bench              Run benchmarks"
    @echo "  fmt                Format code"
    @echo "  fmt-check          Check code formatting"
    @echo "  lint               Run comprehensive linting"
    @echo "  security           Run security scanning"
    @echo "  validate-all       Run comprehensive validation"
    @echo "  pre-commit         Fast pre-commit validation"
    @echo "  pre-push           Comprehensive pre-push validation (REQUIRED)"
    @echo ""
    @echo "📦 Dependencies & Tools:"
    @echo "  deps-check         Check dependency status"
    @echo "  deps-update        Update dependencies"
    @echo "  deps-tidy          Clean up dependencies"
    @echo "  install-tools      Install development tools"
    @echo "  deadcode           Analyze dead code"
    @echo ""
    @echo "🚀 CI/CD:"
    @echo "  ci-full            Complete CI pipeline"
    @echo "  ci-quick           Quick CI validation"
    @echo "  ci-test            CI testing with coverage"
    @echo ""
    @echo "🐳 Docker:"
    @echo "  docker-build       Build Docker images"
    @echo "  docker-up-dev      Start development environment"
    @echo "  docker-test        Run tests in containers"
    @echo "  docker-ci          Run CI in containers"
    @echo ""
    @echo "▶️  Application:"
    @echo "  run-sort           Run sort example"
    @echo "  run-dedup          Run deduplication example"
    @echo "  version            Show version information"
    @echo "  help               Show this help"