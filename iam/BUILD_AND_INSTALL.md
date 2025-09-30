# Building and Installing the HiiRetail IAM Terraform Provider Locally

This guide explains how to build and install the HiiRetail IAM Terraform provider locally for development and testing purposes.

## Prerequisites

- **Go 1.21+**: Required for building the provider
- **Terraform 1.0+**: Required for using the provider
- **Git**: For cloning and version control

## Building the Provider

### 1. Navigate to the Provider Directory

```bash
cd /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam
```

### 2. Download Dependencies

```bash
go mod download
go mod tidy
```

### 3. Build the Provider Binary

```bash
# Build for your current platform
go build -o terraform-provider-hiiretail-iam .

# Or build with version information
go build -ldflags "-X main.version=dev-local" -o terraform-provider-hiiretail-iam .
```

This creates the binary `terraform-provider-hiiretail-iam` in the current directory.

### 4. Verify the Build

```bash
./terraform-provider-hiiretail-iam --help
```

You should see the provider's help output.

## Installing the Provider Locally

### Method 1: Development Override (Recommended for Development)

This method allows you to use the local provider without installing it globally.

#### 1. Create a Development Override Configuration

Create or edit `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "extenda/hiiretail-iam" = "/Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam"
  }

  # For all other providers, use the normal registry
  direct {}
}
```

#### 2. Use in Terraform Configuration

```hcl
terraform {
  required_providers {
    hiiretail_iam = {
      source = "extenda/hiiretail-iam"
    }
  }
}

provider "hiiretail_iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  base_url      = "https://iam-api.retailsvc-test.com"
}
```

#### 3. Test the Provider

```bash
terraform init
terraform plan
```

**Note**: With dev overrides, you'll see a warning that the provider is using a development override.

### Method 2: Local Plugin Directory Installation

This method installs the provider in Terraform's local plugin directory.

#### 1. Create the Plugin Directory Structure

```bash
# Determine your OS and architecture
export OS_ARCH="$(go env GOOS)_$(go env GOARCH)"
export PLUGIN_DIR="$HOME/.terraform.d/plugins/extenda/hiiretail-iam/dev-local/$OS_ARCH"

# Create the directory
mkdir -p "$PLUGIN_DIR"
```

#### 2. Copy the Provider Binary

```bash
cp terraform-provider-hiiretail-iam "$PLUGIN_DIR/"
```

#### 3. Use in Terraform Configuration

```hcl
terraform {
  required_providers {
    hiiretail_iam = {
      source  = "extenda/hiiretail-iam"
      version = "dev-local"
    }
  }
}

provider "hiiretail_iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Advanced Build Options

### Cross-Platform Building

Build for different platforms:

```bash
# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o terraform-provider-hiiretail-iam-linux-amd64 .

# Build for Windows AMD64
GOOS=windows GOARCH=amd64 go build -o terraform-provider-hiiretail-iam-windows-amd64.exe .

# Build for macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o terraform-provider-hiiretail-iam-darwin-arm64 .

# Build for macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o terraform-provider-hiiretail-iam-darwin-amd64 .
```

### Build with Debug Information

```bash
go build -gcflags="all=-N -l" -o terraform-provider-hiiretail-iam .
```

### Optimized Release Build

```bash
go build -ldflags "-s -w -X main.version=1.0.0" -o terraform-provider-hiiretail-iam .
```

## Testing the Local Installation

### 1. Create a Test Configuration

Create a file `test-config.tf`:

```hcl
terraform {
  required_providers {
    hiiretail_iam = {
      source = "extenda/hiiretail_iam"
    }
  }
}

provider "hiiretail_iam" {
  tenant_id     = "test-tenant"
  client_id     = "test-client"
  client_secret = "test-secret"
  base_url      = "https://iam-api.retailsvc-test.com"
}

# Test the role binding resource
resource "hiiretail_iam_role_binding" "test" {
  role_id = "test-role"
  bindings = [
    {
      type = "user"
      id   = "test-user@example.com"
    }
  ]
  description = "Test role binding"
}

# Test the custom role resource
resource "hiiretail_iam_custom_role" "test" {
  name        = "test-role"
  description = "Test custom role"
  permissions = ["read", "write"]
}

# Test the group resource
resource "hiiretail_iam_group" "test" {
  name        = "test-group"
  description = "Test group"
}
```

### 2. Initialize and Validate

```bash
terraform init
terraform validate
terraform plan
```

## Environment Variables for Testing

Set up environment variables for testing:

```bash
export HIIRETAIL_TENANT_ID="your-tenant-id"
export HIIRETAIL_CLIENT_ID="your-client-id"
export HIIRETAIL_CLIENT_SECRET="your-client-secret"
export HIIRETAIL_BASE_URL="https://iam-api.retailsvc-test.com"
```

Then use a simplified provider configuration:

```hcl
provider "hiiretail_iam" {
  # Configuration will be read from environment variables
}
```

## Debugging the Provider

### Debug Mode

Run the provider in debug mode:

```bash
./terraform-provider-hiiretail-iam -debug
```

### Terraform Debug Logging

Enable Terraform debug logging:

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH="./terraform-debug.log"
terraform apply
```

### Provider-Specific Logging

Enable provider-specific logging:

```bash
export TF_LOG_PROVIDER=DEBUG
terraform apply
```

## Makefile for Automation

Create a `Makefile` to automate common tasks:

```makefile
.PHONY: build install test clean

BINARY_NAME=terraform-provider-hiiretail-iam
VERSION=dev-local
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/extenda/hiiretail_iam/$(VERSION)/$(OS_ARCH)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY_NAME) $(PLUGIN_DIR)/

test:
	go test ./... -v

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(PLUGIN_DIR)

dev-override:
	@echo "Add this to your ~/.terraformrc:"
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"extenda/hiiretail_iam\" = \"$(PWD)\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"

.DEFAULT_GOAL := build
```

Save this as `Makefile` in the provider directory, then use:

```bash
# Build the provider
make build

# Install to local plugin directory
make install

# Run tests
make test

# Show dev override configuration
make dev-override

# Clean up
make clean
```

## Troubleshooting

### Common Issues

1. **"Provider not found"**
   - Ensure the binary is in the correct plugin directory
   - Check that the provider source matches your configuration
   - Verify the binary has execute permissions: `chmod +x terraform-provider-hiiretail-iam`

2. **"Version constraints not met"**
   - Make sure the version in your Terraform configuration matches the installed version
   - Use `dev-local` or omit version for development

3. **"Binary not executable"**
   ```bash
   chmod +x terraform-provider-hiiretail-iam
   ```

4. **Build errors**
   ```bash
   go mod tidy
   go mod download
   ```

### Verification Commands

```bash
# Check Go version
go version

# Check Terraform version
terraform version

# Verify provider binary
./terraform-provider-hiiretail-iam --help

# List installed providers
terraform providers

# Check plugin directory
ls -la ~/.terraform.d/plugins/extenda/hiiretail-iam/
```

## Next Steps

After successful installation:

1. **Configure Authentication**: Set up your OIDC client credentials
2. **Test Resources**: Try creating IAM groups, custom roles, and role bindings
3. **Review Documentation**: Check the resource documentation in `docs/`
4. **Run Tests**: Execute the test suite to ensure everything works
5. **Development**: Make changes and rebuild as needed

The provider is now ready for local development and testing!