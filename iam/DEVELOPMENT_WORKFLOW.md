# Terraform Provider Development Workflow

This document explains how to develop and test the HiiRetail IAM Terraform provider locally.

## Prerequisites

1. **Go 1.21+** installed
2. **Terraform 1.0+** installed  
3. **Development overrides** configured in `~/.terraformrc`

## Setup Development Environment

### 1. Configure Terraform Development Overrides

Add this to your `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "extenda/hiiretail-iam" = "/Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam"
  }
  
  # Direct installation for other providers
  direct {}
}
```

**Important**: Replace the path with your actual workspace path.

### 2. Build the Provider

```bash
# In the provider workspace
cd /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam
go build .
```

This creates the `terraform-provider-hiiretail-iam` binary that Terraform will use.

## Development Workflow

### 1. **Make Code Changes**
- Edit provider code in `internal/provider/`
- Edit auth package in `internal/provider/auth/` 
- Edit resource implementations
- Add/modify tests

### 2. **Build and Test**
```bash
# Build the provider
go build .

# Run unit tests
go test ./...

# Run specific auth tests
go test ./internal/provider/auth/...

# Run provider tests
go test ./internal/provider/...
```

### 3. **Test with Terraform**

Create a test Terraform configuration:

```hcl
# test/main.tf
terraform {
  required_providers {
    hiiretail-iam = {
      source = "extenda/hiiretail-iam"
    }
  }
}

provider "hiiretail-iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  
  # Optional: use test environment
  base_url = "https://auth.retailsvc-test.com"
  scopes   = ["iam:read", "iam:write"]
}

resource "hiiretail-iam_custom_role" "example" {
  id   = "example-role"
  name = "Example Role"
  
  permissions = [
    {
      id = "example:read"
      attributes = {
        resource = "example"
      }
    }
  ]
}
```

### 4. **Run Terraform Commands**
```bash
cd test/
terraform init
terraform plan
terraform apply
```

**Note**: Terraform will show a warning about development overrides:
```
Warning: Provider development overrides are in effect
```
This is normal and expected during development.

## How Development Overrides Work

### What Happens When You Run Terraform

1. **Provider Resolution**: Terraform sees `provider "hiiretail-iam"` in your config
2. **Check Overrides**: Terraform checks `~/.terraformrc` for dev_overrides
3. **Use Local Binary**: Instead of downloading from registry, Terraform uses:
   ```
   /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam/terraform-provider-hiiretail-iam
   ```
4. **Execute Provider**: Terraform launches your local binary for all provider operations

### Key Benefits

- ✅ **Instant Testing**: No need to publish/install provider packages  
- ✅ **Debug Friendly**: Easy to add debug flags and logging
- ✅ **Fast Iteration**: Make changes → build → test immediately
- ✅ **Version Control**: No version conflicts with published providers

## Testing OAuth2 Authentication

### Environment Variables

Set up your OAuth2 credentials for testing:

```bash
export HIIRETAIL_TENANT_ID="your-tenant-id"
export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"  
export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"

# Optional: use test environment
export HIIRETAIL_ENVIRONMENT="test"
```

### Run the Demo Program

Test OAuth2 authentication directly:

```bash
cd demo/
go build -o oauth2-demo .
./oauth2-demo
```

### Test in Terraform

```hcl
provider "hiiretail-iam" {
  # Credentials from environment variables
  # OR specify directly (not recommended for production)
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id" 
  client_secret = "your-client-secret"
}
```

## Debugging

### Enable Debug Mode

Build and run the provider in debug mode:

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" .

# Run Terraform with debug logging
export TF_LOG=DEBUG
terraform plan
```

### Provider Debug Mode

Start the provider in debug mode:

```bash
./terraform-provider-hiiretail-iam -debug
```

This will output a `TF_REATTACH_PROVIDERS` environment variable that you can use with Terraform.

### OAuth2 Debug Logging

The auth package includes structured logging. Set log level:

```bash
export HIIRETAIL_LOG_LEVEL=debug
terraform plan
```

## Troubleshooting

### "Provider not found" Error

If Terraform can't find your provider:

1. **Check `.terraformrc`**: Ensure dev_overrides path is correct
2. **Build Provider**: Run `go build .` to create the binary
3. **Check Binary**: Verify `terraform-provider-hiiretail-iam` exists
4. **Path Issues**: Use absolute paths in dev_overrides

### OAuth2 Authentication Errors

1. **Verify Credentials**: Check environment variables are set correctly
2. **Network Access**: Ensure you can reach HiiRetail endpoints
3. **Test Demo**: Run `demo/oauth2-demo` to isolate auth issues
4. **Check Logs**: Enable debug logging to see OAuth2 flow

### Development Override Warnings

This warning is expected during development:
```
Warning: Provider development overrides are in effect
```

To disable dev_overrides (use published providers):
```bash
# Comment out or remove dev_overrides from ~/.terraformrc
provider_installation {
  # dev_overrides {
  #   "extenda/hiiretail-iam" = "/path/to/provider"
  # }
  direct {}
}
```

## Production Deployment

When ready for production:

1. **Remove Dev Overrides**: Comment out dev_overrides in `.terraformrc`
2. **Build Release**: Build provider for target platforms
3. **Publish**: Publish to Terraform Registry or internal registry
4. **Version**: Use proper semantic versioning

## Best Practices

### Development
- ✅ Always build before testing: `go build .`
- ✅ Run tests after changes: `go test ./...`
- ✅ Use test environment for OAuth2 testing
- ✅ Keep credentials in environment variables, not code

### Testing  
- ✅ Test with minimal Terraform configurations first
- ✅ Gradually add complexity to test configurations
- ✅ Test both success and error scenarios
- ✅ Verify OAuth2 token refresh and retry logic

### Security
- ✅ Never commit OAuth2 credentials to version control
- ✅ Use test credentials for development
- ✅ Rotate credentials regularly
- ✅ Review logs for credential exposure

## Quick Reference

```bash
# Build provider
go build .

# Run all tests
go test ./...

# Test OAuth2 auth
cd demo && ./oauth2-demo

# Test with Terraform
cd test && terraform plan

# Debug mode
export TF_LOG=DEBUG
terraform plan
```