# HiiRetail Terraform Provider - Registry Installation Example

This example demonstrates how to use the HiiRetail provider installed from the Terraform Registry.

## Prerequisites

1. **Terraform** version 1.0 or later installed
2. **HiiRetail OAuth2 credentials**:
   - Client ID
   - Client Secret  
   - Tenant ID

## Quick Start

### 1. Create Configuration

Create a new directory and add the following files:

**main.tf**:
```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 1.0"
    }
  }
}

# Configure the HiiRetail Provider
provider "hiiretail" {
  # Credentials are provided via environment variables:
  # - HIIRETAIL_CLIENT_ID
  # - HIIRETAIL_CLIENT_SECRET
  # - HIIRETAIL_TENANT_ID
}

# Example: Create an IAM group
resource "hiiretail_iam_group" "example" {
  name        = "terraform-registry-example"
  description = "Example group created using provider from Terraform Registry"
}

# Example: Output the group ID
output "group_id" {
  description = "The ID of the created IAM group"
  value       = hiiretail_iam_group.example.id
}
```

**terraform.tfvars.example**:
```hcl
# Copy this file to terraform.tfvars and fill in your values
tenant_id     = "your-tenant-id"
client_id     = "your-client-id"
client_secret = "your-client-secret"
```

### 2. Initialize Terraform

The provider will be automatically downloaded from the Terraform Registry:

```bash
terraform init
```

You should see output similar to:
```
Initializing the backend...

Initializing provider plugins...
- Finding extenda/hiiretail versions matching "~> 1.0"...
- Installing extenda/hiiretail v1.0.0...
- Installed extenda/hiiretail v1.0.0 (signed by a HashiCorp partner, key ID ...)

Terraform has been successfully initialized!
```

### 3. Set Environment Variables

Instead of using terraform.tfvars (more secure):

```bash
export HIIRETAIL_CLIENT_ID="your-client-id"
export HIIRETAIL_CLIENT_SECRET="your-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
```

### 4. Plan and Apply

```bash
terraform plan
terraform apply
```

## Version Constraints

Different version constraint examples:

```hcl
# Allow patch updates in 1.x series
version = "~> 1.0"

# Allow minor updates in 1.x series
version = "~> 1.0.0"

# Explicit range
version = ">= 1.0.0, < 2.0.0"

# Pin to exact version
version = "1.0.0"

# Latest version (not recommended for production)
# version = ">= 1.0.0"
```

## Provider Authentication Methods

The provider supports multiple authentication methods with precedence:

### 1. Environment Variables (Recommended)
```bash
export HIIRETAIL_CLIENT_ID="your-client-id"
export HIIRETAIL_CLIENT_SECRET="your-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
```

### 2. Terraform Variables
```hcl
variable "client_id" {
  description = "HiiRetail OAuth2 client ID"
  type        = string
  sensitive   = true
}

provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}
```

### 3. Direct Configuration (Not Recommended)
```hcl
provider "hiiretail" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  tenant_id     = "your-tenant-id"
}
```

## Troubleshooting

### Provider Download Issues

If the provider fails to download:

1. **Check internet connectivity** and proxy settings
2. **Verify Terraform version** is 1.0 or later
3. **Check provider version constraints** match available versions
4. **Clear Terraform cache**:
   ```bash
   rm -rf .terraform .terraform.lock.hcl
   terraform init
   ```

### Authentication Issues

1. **Verify credentials** are correct and have necessary permissions
2. **Check environment variables** are exported in current shell
3. **Validate OAuth2 scopes** include required permissions
4. **Test credentials** with HiiRetail API directly

### Registry Access Issues

1. **Corporate firewall** may block registry.terraform.io
2. **Private registries** may need additional configuration
3. **Mirror configuration** for air-gapped environments

## Next Steps

- Review the [complete documentation](../../docs/)
- Explore more [resource examples](../resources/)
- Check the [authentication guide](../../docs/guides/authentication.md)
- View [GitHub releases](https://github.com/extenda/terraform-provider-hiiretail/releases) for changelog