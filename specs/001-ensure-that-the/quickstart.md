# Quickstart Guide: HiiRetail IAM Terraform Provider

**Date**: September 28, 2025  
**Feature**: Terraform Provider OIDC Authentication and Testing  
**Status**: ✅ READY FOR USE

## Overview

This quickstart guide walks you through setting up and using the HiiRetail IAM Terraform Provider with OIDC client credentials authentication. You'll learn how to configure the provider, test the authentication, and use it to manage IAM resources.

## Prerequisites

Before starting, ensure you have:

- ✅ **Terraform** 1.0+ installed
- ✅ **Go** 1.21+ (for building from source)
- ✅ **OIDC Credentials** for HiiRetail IAM API:
  - Client ID
  - Client Secret  
  - Tenant ID
- ✅ **Network Access** to HiiRetail IAM API endpoints

## Quick Start (5 minutes)

### Step 1: Install the Provider

Add the provider to your Terraform configuration:

```hcl
# main.tf
terraform {
  required_version = ">= 1.0"
  required_providers {
    hiiretail_iam = {
      source  = "extenda/hiiretail_iam"
      version = "~> 1.0"
    }
  }
}
```

### Step 2: Configure Provider Authentication

Create a provider configuration with your OIDC credentials:

```hcl
# main.tf (continued)
provider "hiiretail_iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-oidc-client-id"
  client_secret = "your-oidc-client-secret"
  # base_url is optional - defaults to test environment
  # base_url = "https://iam-api.retailsvc-prod.com"  # for production
}
```

### Step 3: Initialize and Test

```bash
# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Test provider authentication (plan with no resources)
terraform plan
```

**Expected Output:**
```
Terraform will perform the following actions:
  # (no changes required)

Plan: 0 to add, 0 to change, 0 to destroy.
```

✅ **Success!** Your provider is authenticated and ready to use.

## Environment-Specific Setup

### Development Environment

```hcl
provider "hiiretail_iam" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  # Uses default test environment
}
```

### Production Environment

```hcl
provider "hiiretail_iam" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  base_url      = "https://iam-api.retailsvc-prod.com"
}
```

### Using Environment Variables

Set up environment variables for sensitive values:

```bash
# .env (or export directly)
export TF_VAR_tenant_id="your-tenant-id"
export TF_VAR_client_id="your-client-id"
export TF_VAR_client_secret="your-client-secret"
export TF_VAR_base_url="https://custom-api.example.com"
```

```hcl
# variables.tf
variable "tenant_id" {
  description = "HiiRetail tenant ID"
  type        = string
  sensitive   = false
}

variable "client_id" {
  description = "OIDC client ID"
  type        = string
  sensitive   = false
}

variable "client_secret" {
  description = "OIDC client secret"
  type        = string
  sensitive   = true
}

variable "base_url" {
  description = "IAM API base URL"
  type        = string
  default     = "https://iam-api.retailsvc-test.com"
}
```

```hcl
# main.tf
provider "hiiretail_iam" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  base_url      = var.base_url
}
```

## Testing Your Setup

### 1. Authentication Test

Create a simple configuration to test authentication:

```hcl
# test-auth.tf
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
}

# No resources - just test provider configuration
```

Run the test:

```bash
terraform init && terraform plan
```

### 2. Validate Different Environments

Test with different base URLs:

```bash
# Test environment (default)
terraform plan

# Custom environment
terraform plan -var="base_url=https://custom-api.example.com"
```

### 3. Error Handling Test

Test error scenarios to verify proper error handling:

```hcl
# test-errors.tf
provider "hiiretail_iam" {
  tenant_id     = ""  # Empty tenant_id should cause error
  client_id     = "test-client"
  client_secret = "test-secret"
}
```

Expected error:
```
Error: Missing tenant_id
The tenant_id parameter is required
```

## Common Configuration Patterns

### Multi-Environment Setup

```hcl
# environments/dev/main.tf
provider "hiiretail_iam" {
  tenant_id     = "dev-tenant"
  client_id     = var.dev_client_id
  client_secret = var.dev_client_secret
  # Uses default test environment
}

# environments/prod/main.tf  
provider "hiiretail_iam" {
  tenant_id     = "prod-tenant"
  client_id     = var.prod_client_id
  client_secret = var.prod_client_secret
  base_url      = "https://iam-api.retailsvc-prod.com"
}
```

### Workspace-Based Configuration

```hcl
locals {
  environment_config = {
    dev = {
      base_url = "https://iam-api.retailsvc-test.com"
      tenant_id = "dev-tenant"
    }
    staging = {
      base_url = "https://iam-api.retailsvc-staging.com"
      tenant_id = "staging-tenant"
    }
    prod = {
      base_url = "https://iam-api.retailsvc-prod.com"
      tenant_id = "prod-tenant"
    }
  }
  
  current_env = local.environment_config[terraform.workspace]
}

provider "hiiretail_iam" {
  tenant_id     = local.current_env.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  base_url      = local.current_env.base_url
}
```

## Troubleshooting

### Common Issues and Solutions

#### 1. Authentication Failures

**Problem**: `Error: Missing client_id`

**Solution**: Verify all required parameters are provided:
```hcl
provider "hiiretail_iam" {
  tenant_id     = "required-value"     # ✅ Required
  client_id     = "required-value"     # ✅ Required  
  client_secret = "required-value"     # ✅ Required
  base_url      = "optional-value"     # ❌ Optional
}
```

#### 2. Invalid Base URL

**Problem**: `Error: Invalid base_url`

**Solution**: Ensure base_url has proper format:
```hcl
provider "hiiretail_iam" {
  # ... other config ...
  base_url = "https://api.example.com"  # ✅ Valid
  # base_url = "not-a-url"             # ❌ Invalid
  # base_url = "http://insecure.com"   # ❌ Must use HTTPS
}
```

#### 3. Network/Connectivity Issues

**Problem**: Connection timeouts or network errors

**Solution**: 
- Verify network access to the base_url
- Check firewall rules
- Confirm DNS resolution
- Test with curl: `curl -v https://your-base-url/oauth/token`

#### 4. Token/Authentication Issues

**Problem**: Authentication fails with valid credentials

**Solution**:
- Verify client_id and client_secret are correct
- Check if credentials are expired or revoked
- Confirm tenant_id matches your environment
- Test credentials manually with OAuth2 tools

### Debug Mode

Enable Terraform debug logging for detailed error information:

```bash
export TF_LOG=DEBUG
terraform plan
```

Look for OAuth2 and HTTP client debug information in the logs.

### Validation Checklist

Use this checklist to verify your setup:

- [ ] **Provider Block**: Complete provider configuration with all required fields
- [ ] **Credentials**: Valid client_id, client_secret, and tenant_id
- [ ] **Network**: Access to base_url (test with curl/browser)
- [ ] **Environment**: Correct base_url for your target environment
- [ ] **Terraform**: Version 1.0+ installed and working
- [ ] **Variables**: Sensitive variables properly marked and secured
- [ ] **Initialization**: `terraform init` runs successfully
- [ ] **Validation**: `terraform validate` passes
- [ ] **Planning**: `terraform plan` completes without authentication errors

## Next Steps

After successful setup:

1. **Add Resources**: Start adding IAM resources to your configuration
2. **Version Control**: Commit your configuration to version control
3. **CI/CD Integration**: Set up automated deployments
4. **Monitoring**: Implement monitoring for your Terraform deployments
5. **Documentation**: Document your specific configuration patterns

## Getting Help

If you encounter issues:

1. **Check Logs**: Use `TF_LOG=DEBUG` for detailed error information
2. **Validate Config**: Run `terraform validate` to check syntax
3. **Test Connectivity**: Verify network access to your base_url
4. **Review Documentation**: Check the full provider documentation
5. **Contact Support**: Reach out to the HiiRetail support team

## Complete Example

Here's a complete working example:

```hcl
# Complete example configuration
terraform {
  required_version = ">= 1.0"
  required_providers {
    hiiretail_iam = {
      source  = "extenda/hiiretail_iam"
      version = "~> 1.0"
    }
  }
}

# Variables for sensitive values
variable "tenant_id" {
  description = "HiiRetail tenant ID"
  type        = string
}

variable "client_id" {
  description = "OIDC client ID"
  type        = string
}

variable "client_secret" {
  description = "OIDC client secret"
  type        = string
  sensitive   = true
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "test"
  validation {
    condition     = contains(["test", "staging", "prod"], var.environment)
    error_message = "Environment must be test, staging, or prod."
  }
}

# Environment-specific configuration
locals {
  base_urls = {
    test    = "https://iam-api.retailsvc-test.com"
    staging = "https://iam-api.retailsvc-staging.com"
    prod    = "https://iam-api.retailsvc-prod.com"
  }
}

# Provider configuration
provider "hiiretail_iam" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  base_url      = local.base_urls[var.environment]
}

# Your IAM resources would go here...
```

Save this as `main.tf`, set your variables, and run:

```bash
terraform init
terraform plan -var="tenant_id=your-tenant" -var="client_id=your-client" -var="client_secret=your-secret"
```

🎉 **Congratulations!** You've successfully set up the HiiRetail IAM Terraform Provider with OIDC authentication.
