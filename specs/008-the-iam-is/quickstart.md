# Quickstart Guide: Multi-API HiiRetail Provider

## Prerequisites
- HiiRetail OAuth2 credentials (client_id and client_secret)
- Terraform 1.0+ installed
- Access to HiiRetail APIs

**Estimated Time**: 15 minutes

## Step 1: Provider Configuration

Create a new directory and `main.tf` file:

```hcl
terraform {
  required_providers {
    hiiretail = {
      source = "extenda/hiiretail"
      version = ">= 1.0.0"
    }
  }
}

provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
}

variable "client_id" {
  description = "HiiRetail OAuth2 Client ID"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "HiiRetail OAuth2 Client Secret" 
  type        = string
  sensitive   = true
}
```

## Step 2: Create First Resource

Add an IAM group resource to your `main.tf`:

```hcl
resource "hiiretail_iam_group" "quickstart_group" {
  name        = "quickstart-demo-group"
  description = "Demo group created during quickstart"
}

output "group_id" {
  description = "ID of the created group"
  value       = hiiretail_iam_group.quickstart_group.id
}

output "group_name" {
  description = "Name of the created group"
  value       = hiiretail_iam_group.quickstart_group.name
}
```

## Step 3: Configure Credentials

Create `terraform.tfvars` file with your credentials:

```hcl
client_id     = "your-client-id-here"
client_secret = "your-client-secret-here"
```

**Security Note**: Never commit `terraform.tfvars` to version control. Add it to your `.gitignore`.

## Step 4: Initialize and Apply

```bash
# Initialize Terraform and download the provider
terraform init

# Review the planned changes
terraform plan

# Apply the configuration
terraform apply
```

Expected output:
```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

group_id = "generated-group-id"
group_name = "quickstart-demo-group"
```

## Step 5: Verify and Explore

### Verify Resource Creation
Check your HiiRetail IAM console to confirm the group was created.

### Explore Provider Resources
```bash
# List all available resources
terraform providers schema -json | jq '.provider_schemas["registry.terraform.io/extenda/hiiretail"].resource_schemas | keys[]'
```

Expected resource types:
- `hiiretail_iam_group`
- `hiiretail_iam_custom_role`
- `hiiretail_iam_role_binding`

## Step 6: Create Advanced Configuration

Expand your configuration to include multiple resources:

```hcl
# Custom role with specific permissions
resource "hiiretail_iam_custom_role" "quickstart_role" {
  id   = "quickstart-demo-role"
  name = "Quickstart Demo Role"
  
  permissions = [
    {
      id = "pos.products.read"
      attributes = {}
    },
    {
      id = "pos.products.write"  
      attributes = {}
    }
  ]
}

# Role binding connecting group and role
resource "hiiretail_iam_role_binding" "quickstart_binding" {
  role_id = hiiretail_iam_custom_role.quickstart_role.id
  bindings = [
    hiiretail_iam_group.quickstart_group.id
  ]
  is_custom = true
}
```

Apply the expanded configuration:
```bash
terraform plan
terraform apply
```

## Step 7: Data Sources

Use data sources to reference existing resources:

```hcl
# Reference existing groups
data "hiiretail_iam_groups" "all_groups" {}

output "all_group_names" {
  description = "Names of all IAM groups"
  value       = data.hiiretail_iam_groups.all_groups.groups[*].name
}
```

## Step 8: Clean Up

When finished with the quickstart:

```bash
terraform destroy
```

## Next Steps

### Explore Documentation
- [Authentication Guide](../guides/authentication.md) - Advanced OAuth2 configuration
- [IAM Resources](../resources/iam/overview.md) - Complete IAM resource documentation
- [Examples](../examples/) - More complex usage patterns

### Migration from Old Provider
If you're migrating from the `hiiretail-iam` provider:
- [Migration Guide](../guides/migration-guides/from-hiiretail-iam.md)

### Advanced Usage
- [Multi-Service Examples](../examples/multi-service-deployment/) - Using multiple HiiRetail APIs
- [Enterprise Patterns](../examples/enterprise-patterns/) - Large-scale deployments

## Validation Checklist

✅ Provider initializes without errors  
✅ Resources can be planned and applied  
✅ Resources appear in HiiRetail console  
✅ Data sources return expected information  
✅ Resources can be cleanly destroyed  

## Troubleshooting

**Common Issues:**

1. **Authentication Error**: Verify client_id and client_secret are correct
2. **Network Error**: Check internet connectivity and firewall settings
3. **Permission Error**: Ensure OAuth2 credentials have necessary API permissions
4. **Resource Not Found**: Verify resource exists and you have read permissions

**Getting Help:**
- Check provider documentation for detailed resource information
- Review error messages for specific troubleshooting steps  
- Consult examples directory for working configurations

## What You Learned

- How to configure the multi-API HiiRetail provider
- Basic resource creation and management patterns
- Using data sources to reference existing resources
- Provider naming conventions and resource organization
- Security best practices for credential management

This quickstart demonstrates the improved user experience of the unified HiiRetail provider compared to the previous service-specific approach.