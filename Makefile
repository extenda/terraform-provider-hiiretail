.PHONY: build install test clean dev-override fmt lint

BINARY_NAME=terraform-provider-hiiretail
VERSION=dev-local
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/extenda/hiiretail/$(VERSION)/$(OS_ARCH)

# Default target
all: build

# Build the provider binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "✅ Build complete: $(BINARY_NAME)"

# Install to local Terraform plugin directory
install: build
	@echo "Installing to $(PLUGIN_DIR)..."
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY_NAME) $(PLUGIN_DIR)/
	chmod +x $(PLUGIN_DIR)/$(BINARY_NAME)
	@echo "✅ Provider installed locally"

# Run all tests
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "✅ Tests complete"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test ./internal/provider/resource_iam_role_binding -v
	go test ./internal/provider/resource_iam_custom_role -v
	go test ./internal/provider/iam/resources -v
	@echo "✅ Unit tests complete"

# Run acceptance tests
test-acceptance:
	@echo "Running acceptance tests..."
	TF_ACC=1 go test ./internal/provider -v -timeout 30m
	@echo "✅ Acceptance tests complete"

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint Go code with security focus
lint:
	@echo "Linting Go code with security checks..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, falling back to basic linting"; \
		go vet ./...; \
	fi
	@echo "✅ Code linted"

# Security audit for OAuth2 credential handling
security-audit:
	@echo "Running security audit for OAuth2 implementation..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --enable gosec --disable-all; \
	else \
		echo "⚠️  golangci-lint not found, install it for security auditing"; \
		echo "Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
	fi
	@echo "✅ Security audit complete"

# Install golangci-lint for enhanced linting
install-lint:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	@echo "✅ golangci-lint installed"

# Setup OAuth2 security git hooks
setup-security-hooks:
	@echo "Setting up OAuth2 security git hooks..."
	git config core.hooksPath .githooks
	@echo "✅ Security git hooks enabled"
	@echo "   Pre-commit hook will now check for credential exposure"

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(PLUGIN_DIR)
	@echo "✅ Clean complete"

# Show development override configuration
dev-override:
	@echo ""
	@echo "🔧 Development Override Configuration"
	@echo "Add this to your ~/.terraformrc file:"
	@echo ""
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"extenda/hiiretail\" = \"$(PWD)\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"
	@echo ""
	@echo "Then run: terraform init"

# Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Multi-platform build complete"

# Quick development setup
dev-setup: deps build dev-override
	@echo ""
	@echo "🚀 Development setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Add the dev override configuration to ~/.terraformrc (shown above)"
	@echo "2. Create a test Terraform configuration"
	@echo "3. Run: terraform init && terraform plan"

# Release validation and testing targets
validate-release-config:
	@echo "Validating release configuration..."
	./test_goreleaser_config.sh
	./test_github_actions.sh
	@echo "✅ Release configuration validation complete"

test-build-local:
	@echo "Testing local build with GoReleaser..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser build --snapshot --clean --timeout 10m; \
		./test_build_artifacts.sh; \
	else \
		echo "❌ GoReleaser not installed. Install with: brew install goreleaser"; \
		exit 1; \
	fi
	@echo "✅ Local build test complete"

test-release-process:
	@echo "Testing complete release process..."
	./test_release_process.sh
	@echo "✅ Release process test complete"

validate-all-release: validate-release-config test-build-local test-release-process
	@echo "🎉 All release validation tests passed!"

# Install release tooling
install-release-tools:
	@echo "Installing release tooling..."
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "Installing GoReleaser..."; \
		brew install goreleaser; \
	else \
		echo "✅ GoReleaser already installed"; \
	fi
	@if ! command -v act >/dev/null 2>&1; then \
		echo "Installing act (GitHub Actions local runner)..."; \
		brew install act; \
	else \
		echo "✅ act already installed"; \
	fi
	@echo "✅ Release tooling installation complete"

# Show help
help:
	@echo ""
	@echo "🏗️  HiiRetail Unified Terraform Provider Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build                     Build the provider binary"
	@echo "  install                   Build and install to local plugin directory"
	@echo "  test                      Run all tests"
	@echo "  test-unit                 Run unit tests only"
	@echo "  test-acceptance           Run acceptance tests only"
	@echo "  fmt                       Format Go code"
	@echo "  lint                      Lint Go code with security checks"
	@echo "  security-audit            Run security audit for OAuth2 credential handling"
	@echo "  install-lint              Install golangci-lint for enhanced linting"
	@echo "  setup-security-hooks      Setup git hooks for credential exposure detection"
	@echo "  clean                     Clean up build artifacts"
	@echo "  dev-override              Show development override configuration"
	@echo "  deps                      Download and tidy dependencies"
	@echo "  build-all                 Build for multiple platforms"
	@echo "  dev-setup                 Complete development setup"
	@echo ""
	@echo "Release targets:"
	@echo "  validate-release-config   Validate GoReleaser and GitHub Actions configuration"
	@echo "  test-build-local          Test local build with GoReleaser"
	@echo "  test-release-process      Test complete release process"
	@echo "  validate-all-release      Run all release validation tests"
	@echo "  install-release-tools     Install GoReleaser and act"
	@echo ""
	@echo "  help                      Show this help message"
	@echo ""

.DEFAULT_GOAL := build