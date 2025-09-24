# Multi-stage Dockerfile for ImageManager Go Application
# Supports production, development, testing, CI, and QA environments

# ============================================================================
# Base stage - Common Go environment
# ============================================================================
FROM golang:1.24-alpine AS base

# Install essential tools and dependencies
RUN apk add --no-cache \
    git \
    make \
    bash \
    curl \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Set working directory
WORKDIR /app


# Copy go workspace configuration
COPY go.work go.work.sum ./

# ============================================================================
# Dependencies stage - Download and cache Go modules
# ============================================================================
FROM base AS dependencies

# Copy all go.mod and go.sum files for efficient caching
COPY cli/go.mod cli/go.sum ./cli/
COPY core/go.mod core/go.sum ./core/
COPY integration_test/go.mod integration_test/go.sum ./integration_test/

# Download dependencies for all modules
RUN cd cli && go mod download
RUN cd core && go mod download
RUN cd integration_test && go mod download

# ============================================================================
# Development stage - Full development environment
# ============================================================================
FROM dependencies AS development

# Install development tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.3.0 && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install golang.org/x/tools/cmd/deadcode@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Copy source code
COPY . .

# Create necessary directories
RUN mkdir -p /app/logs /app/output /app/test_results

# Set development environment
ENV GO_ENV=development
ENV CGO_ENABLED=1

# Default command for development
CMD ["tail", "-f", "/dev/null"]

# ============================================================================
# Testing stage - Optimized for running tests
# ============================================================================
FROM development AS testing

# Install additional test tools
RUN go install github.com/axw/gocov/gocov@latest && \
    go install github.com/matm/gocov-html@latest

# Set test environment
ENV GO_ENV=test
ENV CGO_ENABLED=1

# Create test results directory
RUN mkdir -p /app/test_results

# Default command for testing  
CMD ["sh", "-c", "echo '🧪 Running unit tests...' && go test -race -v -timeout=2m -short ./cli/... ./core/..."]

# ============================================================================
# CI stage - CI/CD optimized environment
# ============================================================================
FROM dependencies AS ci

# Install CI-specific tools (minimal set for performance)
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.3.0 && \
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Copy source code (read-only for CI)
COPY . .

# Create CI results directory
RUN mkdir -p /app/ci_results

# Set CI environment
ENV GO_ENV=ci
ENV CGO_ENABLED=0

# Default command for CI
CMD ["make", "ci-full"]

# ============================================================================
# QA stage - Quality assurance with all analysis tools
# ============================================================================
FROM development AS qa

# Install additional QA tools
RUN apk add --no-cache shellcheck && \
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest && \
    go install github.com/client9/misspell/cmd/misspell@latest

# Create QA results directory
RUN mkdir -p /app/qa_results

# Set QA environment
ENV GO_ENV=qa
ENV CGO_ENABLED=1

# Default command for QA
CMD ["make", "quality-audit"]

# ============================================================================
# Builder stage - Build optimized binary
# ============================================================================
FROM dependencies AS builder

# Copy source code
COPY . .

# Build the application with optimizations
RUN cd cli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'docker') -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o ImageManager .

# ============================================================================
# Production stage - Minimal production image
# ============================================================================
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S imagemanager && \
    adduser -u 1001 -S imagemanager -G imagemanager

# Set working directory
WORKDIR /app

# Create necessary directories and set permissions
RUN mkdir -p /app/input /app/output /app/config /app/logs && \
    chown -R imagemanager:imagemanager /app

# Copy binary from builder stage
COPY --from=builder --chown=imagemanager:imagemanager /app/cli/ImageManager /usr/local/bin/imagemanager

# Copy configuration examples
COPY --from=builder --chown=imagemanager:imagemanager /app/cli/example/ /app/config/

# Switch to non-root user
USER imagemanager

# Set production environment
ENV GO_ENV=production

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD imagemanager --help > /dev/null || exit 1

# Default command
CMD ["imagemanager", "--help"]
