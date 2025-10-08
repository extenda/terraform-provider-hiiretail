# HiiRetail IAM Resource Examples

This directory contains comprehensive examples and guides for using the `hiiretail_iam_resource` resource.

## Quick Start

The simplest way to create a resource:

```hcl
terraform {
  required_providers {
    hiiretail = {
      source = "extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  #Optional configuration
}

resource "hiiretail_iam_resource" "my_first_resource" {
  id   = "bu:001"
  name = "My First Store"
}
```

## Complete Examples

### [Basic Resource](./basic/main.tf)
Simple resource creation with minimal configuration.

### [Resource with Properties](./with-properties/main.tf)
Resource with complex JSON properties for metadata and configuration.

### [Multi-Store Setup](./multi-store/main.tf)
Managing multiple store resources with consistent naming patterns.

### [Department Hierarchy](./department-hierarchy/main.tf)  
Creating hierarchical department resources with relationships.

### [Application Resources](./application-resources/main.tf)
Resources representing applications and services in your infrastructure.

### [POS Terminal Resources](./pos-terminals/main.tf)
Point-of-sale terminal resources with location and configuration data.

### [Integration with Roles](./with-roles/main.tf)
Complete RBAC setup using resources with custom roles and bindings.

### [Import Existing Resources](./import-existing/import.sh)
Scripts and examples for importing existing resources into Terraform.

## Usage Patterns

### Resource Naming Conventions
- **Stores**: `store:location-id` or `store:region:location`
- **Departments**: `dept:name` or `dept:store:name`  
- **Applications**: `app:service-name` or `app:environment:service`
- **POS Systems**: `pos:store:terminal` or `pos:location:device-id`
- **User Resources**: `user:category:identifier`

See [naming-conventions.md](./naming-conventions.md) for detailed guidelines.

## Common Scenarios

### Scenario 1: Retail Chain Management
```hcl
# Main store resource
resource "hiiretail_iam_resource" "flagship_store" {
  id   = "store:flagship:nyc"
  name = "Flagship Store - New York City"
  props = jsonencode({
    address = "123 Broadway, New York, NY 10001"
    manager = "sarah.johnson@company.com"
    square_footage = 15000
    departments = ["electronics", "clothing", "home"]
  })
}

# Department resources within the store
resource "hiiretail_iam_resource" "electronics_dept" {
  id   = "dept:flagship:electronics"
  name = "Electronics Department - Flagship"
  props = jsonencode({
    store_id = hiiretail_iam_resource.flagship_store.id
    manager = "mike.chen@company.com"
    budget = 500000
    categories = ["mobile", "computers", "gaming"]
  })
}
```

### Scenario 2: Multi-Environment Application Management
```hcl
# Production inventory app
resource "hiiretail_iam_resource" "inventory_prod" {
  id   = "app:prod:inventory"
  name = "Inventory Management - Production"
  props = jsonencode({
    environment = "production"
    version = "2.1.4"
    endpoints = ["https://api.company.com/inventory"]
    database = "inventory-prod-db"
    replicas = 3
  })
}

# Staging inventory app
resource "hiiretail_iam_resource" "inventory_staging" {
  id   = "app:staging:inventory"
  name = "Inventory Management - Staging"
  props = jsonencode({
    environment = "staging"
    version = "2.2.0-beta"
    endpoints = ["https://staging-api.company.com/inventory"]
    database = "inventory-staging-db"
    replicas = 1
  })
}
```

### Scenario 3: POS Terminal Fleet Management
```hcl
# Generate POS terminals for multiple stores
locals {
  stores = {
    "store-001" = { name = "Downtown Store", terminals = 3 }
    "store-002" = { name = "Mall Location", terminals = 5 }
    "store-003" = { name = "Airport Store", terminals = 2 }
  }
}

# Create store resources
resource "hiiretail_iam_resource" "stores" {
  for_each = local.stores
  
  id   = "store:${each.key}"
  name = each.value.name
  props = jsonencode({
    terminal_count = each.value.terminals
    status = "active"
  })
}

# Create POS terminal resources
resource "hiiretail_iam_resource" "pos_terminals" {
  for_each = {
    for combo in flatten([
      for store_id, store in local.stores : [
        for i in range(store.terminals) : {
          store_id = store_id
          terminal_id = i + 1
          name = "${store.name} - Terminal ${i + 1}"
        }
      ]
    ]) : "${combo.store_id}-terminal-${combo.terminal_id}" => combo
  }
  
  id   = "pos:${each.value.store_id}:terminal-${each.value.terminal_id}"
  name = each.value.name
  props = jsonencode({
    store_id = each.value.store_id
    terminal_number = each.value.terminal_id
    model = "VeriFone MX925"
    status = "active"
  })
}
```

## Error Handling Examples

### Validation Errors
```hcl
# This will fail - invalid ID format
resource "hiiretail_iam_resource" "invalid_id" {
  id   = "store/001"  # Forward slash not allowed
  name = "Invalid Store"
}

# This will fail - consecutive underscores
resource "hiiretail_iam_resource" "invalid_id_2" {
  id   = "store__001"  # Double underscores not allowed  
  name = "Invalid Store"
}

# This will fail - reserved pattern
resource "hiiretail_iam_resource" "invalid_id_3" {
  id   = "."  # Reserved pattern
  name = "Invalid Store"
}
```

### JSON Validation Errors
```hcl
# This will fail - invalid JSON
resource "hiiretail_iam_resource" "invalid_json" {
  id   = "store:001"
  name = "Store with Invalid JSON"
  props = "{invalid: json}"  # Missing quotes around key
}

# Correct version
resource "hiiretail_iam_resource" "valid_json" {
  id   = "store:001"
  name = "Store with Valid JSON"
  props = jsonencode({
    status = "active"
    location = "downtown"
  })
}
```

## Best Practices

### 1. Use Consistent Naming Patterns
```hcl
# Good: Consistent hierarchy
resource "hiiretail_iam_resource" "store_main" {
  id = "store:region-east:location-001"
  name = "East Region - Store 001"
}

resource "hiiretail_iam_resource" "dept_electronics" {
  id = "dept:region-east:location-001:electronics"  
  name = "Electronics - East Region Store 001"
}

# Avoid: Inconsistent patterns
resource "hiiretail_iam_resource" "random_store" {
  id = "some-store-here"  # No clear hierarchy
  name = "Random Store"
}
```

### 2. Structure Properties Consistently
```hcl
resource "hiiretail_iam_resource" "well_structured" {
  id   = "store:001"
  name = "Well Structured Store"
  props = jsonencode({
    # Metadata - always first
    created_at = "2024-01-15T10:30:00Z"
    created_by = "admin@company.com"
    version    = "1.0"
    
    # Core business properties
    location = {
      address = "123 Main St"
      city    = "New York"
      state   = "NY"
      zip     = "10001"
    }
    
    # Operational data
    manager = "jane.doe@company.com"
    phone   = "+1-555-0123"
    
    # Configuration
    settings = {
      auto_backup = true
      sync_enabled = true
      notifications = true
    }
    
    # References to other resources
    parent_resource = "region:east"
    child_resources = ["dept:electronics", "dept:clothing"]
  })
}
```

### 3. Use Variables for Reusability
```hcl
variable "stores" {
  description = "Map of store configurations"
  type = map(object({
    name     = string
    address  = string
    manager  = string
    departments = list(string)
  }))
  default = {
    "store-001" = {
      name = "Downtown Store"
      address = "123 Main St, New York, NY"
      manager = "john.doe@company.com"
      departments = ["electronics", "clothing"]
    }
    "store-002" = {
      name = "Mall Store"
      address = "456 Mall Ave, New York, NY"
      manager = "jane.smith@company.com"
      departments = ["electronics", "home", "beauty"]
    }
  }
}

resource "hiiretail_iam_resource" "stores" {
  for_each = var.stores
  
  id   = "store:${each.key}"
  name = each.value.name
  props = jsonencode({
    address = each.value.address
    manager = each.value.manager
    departments = each.value.departments
    created_at = timestamp()
  })
}
```

### 4. Use Outputs for Integration
```hcl
output "store_resources" {
  description = "Map of created store resources"
  value = {
    for k, v in hiiretail_iam_resource.stores : k => {
      id = v.id
      name = v.name
      tenant_id = v.tenant_id
    }
  }
}

output "store_ids" {
  description = "List of store resource IDs for use in roles"
  value = [for store in hiiretail_iam_resource.stores : store.id]
}

# Use in other configurations
data "terraform_remote_state" "stores" {
  backend = "local"
  config = {
    path = "../stores/terraform.tfstate"
  }
}

resource "hiiretail_iam_custom_role" "store_manager" {
  for_each = data.terraform_remote_state.stores.outputs.store_resources
  
  name = "Manager - ${each.value.name}"
  permissions = [
    "inventory.${each.value.id}.read",
    "inventory.${each.value.id}.write",
    "sales.${each.value.id}.read"
  ]
}
```

## Testing Your Configuration

### 1. Validation
```bash
# Validate syntax and configuration
terraform validate

# Format code consistently
terraform fmt -recursive

# Check for potential issues
terraform plan
```

### 2. Dry Run Testing
```bash
# Create a test environment
cp main.tf test.tf

# Modify resource IDs for testing
sed -i 's/store:/test-store:/g' test.tf

# Apply to test environment
terraform apply -var-file="test.tfvars"

# Clean up
terraform destroy
rm test.tf
```

### 3. Resource Verification
```hcl
# Add data source to verify resource exists
data "external" "verify_resource" {
  program = ["bash", "-c", <<-EOT
    # This would make an API call to verify the resource
    # For example purposes, we'll return success
    echo '{"status": "verified", "id": "${hiiretail_iam_resource.my_resource.id}"}'
  EOT
  ]
  
  depends_on = [hiiretail_iam_resource.my_resource]
}

output "verification_result" {
  value = data.external.verify_resource.result
}
```

## Troubleshooting Guide

### Common Issues and Solutions

#### Issue: "Resource ID format validation failed"
```
Error: Invalid Resource ID
Resource ID 'store/001' is invalid: resource ID cannot contain forward slashes
```

**Solution**: Use colons (`:`) or hyphens (`-`) instead of slashes:
```hcl
# Wrong
id = "store/001"

# Correct
id = "store:001"
# or
id = "store-001"
```

#### Issue: "JSON validation failed in props field"
```
Error: Invalid JSON in props field
```

**Solution**: Always use `jsonencode()` function:
```hcl
# Wrong
props = '{"key": "value"}'

# Correct
props = jsonencode({
  key = "value"
})
```

#### Issue: "Authentication failed"
```
Error: Authentication Failed
OAuth2 credentials are invalid or expired
```

**Solutions**:
1. Check your client credentials:
```hcl
provider "hiiretail" {
  client_id     = var.client_id      # Verify this is correct
  client_secret = var.client_secret  # Verify this is correct
  tenant_id     = var.tenant_id      # Verify this is correct
  token_url     = "https://auth.hiiretail.com/oauth/token"
}
```

2. Verify environment variables:
```bash
export HIIRETAIL_CLIENT_ID="your-client-id"
export HIIRETAIL_CLIENT_SECRET="your-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
```

3. Test authentication separately:
```bash
curl -X POST "https://auth.hiiretail.com/oauth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=your-client-id&client_secret=your-client-secret"
```

#### Issue: "Permission denied"
```
Error: Permission Denied
You don't have permission to create resource 'store:001'
```

**Solutions**:
1. Check OAuth2 scopes:
```hcl
provider "hiiretail" {
  # ... other config
  scopes = ["iam:read", "iam:write"]  # Ensure both scopes are included
}
```

2. Verify account permissions with your administrator.

3. Test with a simpler resource first to isolate permission issues.

#### Issue: "Resource already exists"
```
Error: Resource Conflict
A resource with ID 'store:001' already exists
```

**Solutions**:
1. Import the existing resource:
```bash
terraform import hiiretail_iam_resource.store_001 store:001
```

2. Use a different resource ID:
```hcl
resource "hiiretail_iam_resource" "store_001" {
  id = "store:001-new"  # Different ID
  name = "Store 001"
}
```

3. Delete the existing resource first (if safe to do so).

#### Issue: "Rate limit exceeded"
```
Error: Rate Limit Exceeded
Too many requests to the API
```

**Solutions**:
1. Add delays between operations:
```bash
terraform apply
sleep 10
terraform apply  # For additional resources
```

2. Reduce parallelism:
```bash
terraform apply -parallelism=1
```

3. Contact support to increase rate limits if needed.

### Debugging Techniques

#### 1. Enable Debug Logging
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log
terraform apply
```

#### 2. Isolate Issues
```hcl
# Comment out all but one resource to isolate issues
resource "hiiretail_iam_resource" "test_only" {
  id   = "test:resource"
  name = "Test Resource"
}

# Uncomment others one by one after testing
# resource "hiiretail_iam_resource" "second_resource" {
#   id   = "test:resource2"  
#   name = "Second Resource"
# }
```

#### 3. Use terraform console for Testing
```bash
terraform console

# Test jsonencode function
> jsonencode({key = "value"})
"{"key":"value"}"

# Test variable values
> var.tenant_id
"your-tenant-id"
```

#### 4. Check State File
```bash
# View current state
terraform show

# List resources in state
terraform state list

# Show specific resource
terraform state show hiiretail_iam_resource.store_001
```

### Getting Help

If you continue to experience issues:

1. **Check the logs**: Enable debug logging and review the output
2. **Validate configuration**: Use `terraform validate` and `terraform plan`
3. **Review documentation**: Check the [resource documentation](../../docs/resources/hiiretail_iam_resource.md)
4. **Contact support**: Provide your tenant ID, resource configuration, and full error messages

### Support Information Template

When contacting support, include:

```
**Tenant ID**: your-tenant-id
**Resource ID**: store:001
**Terraform Version**: 1.6.0
**Provider Version**: 1.0.0

**Configuration**:
```hcl
resource "hiiretail_iam_resource" "problematic_resource" {
  id   = "store:001"
  name = "Problem Store"
  props = jsonencode({
    location = "test"
  })
}
```

**Error Message**:
```
Error: Authentication Failed
The create operation failed...
```

**Steps to Reproduce**:
1. Configure provider with credentials
2. Run terraform plan (succeeds)
3. Run terraform apply (fails)

**Debug Logs**: [Attach terraform.log file]
```