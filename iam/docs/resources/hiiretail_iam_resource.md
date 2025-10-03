# hiiretail_iam_resource

Provides a HiiRetail IAM resource for managing granular access control within your tenant. Resources represent logical entities in your system (stores, departments, applications, etc.) that can be referenced in role-based access control (RBAC) policies.

## Example Usage

### Basic Resource

```hcl
resource "hiiretail_iam_resource" "store_001" {
  id   = "store:001"
  name = "Main Store - Downtown"
}
```

### Resource with Properties

```hcl
resource "hiiretail_iam_resource" "pos_terminal" {
  id   = "pos:store-001:terminal-01"
  name = "POS Terminal 01 - Main Store"
  props = jsonencode({
    location    = "checkout-1"
    department  = "electronics"
    active      = true
    metadata = {
      install_date = "2024-01-15"
      model       = "VeriFone MX925"
    }
  })
}
```

### Department Resource

```hcl
resource "hiiretail_iam_resource" "electronics_dept" {
  id   = "dept:electronics"
  name = "Electronics Department"
  props = jsonencode({
    manager     = "john.doe@company.com"
    budget      = 50000
    categories  = ["mobile", "computers", "accessories"]
    permissions = {
      inventory_read  = true
      inventory_write = true
      reports_access  = true
    }
  })
}
```

### Application Resource

```hcl
resource "hiiretail_iam_resource" "inventory_app" {
  id   = "app:inventory-management"
  name = "Inventory Management System"
  props = jsonencode({
    version     = "2.1.4"
    environment = "production"
    endpoints = [
      "https://api.company.com/inventory",
      "https://api.company.com/reports"
    ]
    features = {
      real_time_sync = true
      batch_import   = true
      audit_trail    = true
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) Unique identifier for the resource within the tenant. Must follow the pattern rules:
  - 1-1500 characters long
  - Cannot contain forward slashes (`/`)
  - Cannot be `.` or `..`
  - Cannot contain consecutive underscores (`__`)
  - Examples: `store:001`, `dept:electronics`, `pos:store-001:terminal-01`

* `name` - (Required) Human-readable display name for the resource. Must be a non-empty string.

* `props` - (Optional) JSON-encoded string containing additional properties and metadata for the resource. Can include any valid JSON data types (string, number, boolean, array, object). Use `jsonencode()` function for complex objects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tenant_id` - The tenant identifier, inherited from the provider configuration.

## Import

HiiRetail IAM resources can be imported using their resource ID:

```shell
terraform import hiiretail_iam_resource.example store:001
```

The resource ID must follow the same validation rules as when creating a resource. After import, you can run `terraform plan` to see what changes (if any) need to be applied to match your configuration.

## Resource ID Patterns

Valid resource ID patterns:

* **Store Resources**: `store:001`, `store:main-branch`, `store:location-nyc`
* **Department Resources**: `dept:electronics`, `dept:clothing`, `dept:home-garden`
* **POS Resources**: `pos:store-001:terminal-01`, `pos:mobile:tablet-05`
* **Application Resources**: `app:inventory-management`, `app:customer-portal`
* **User Profile Resources**: `user.profile.settings`, `user.preferences.theme`
* **Generic Resources**: `resource-name`, `my-custom-resource`, `data:customer:12345`

Invalid patterns (will cause validation errors):

* `store/001` - Contains forward slash
* `store__001` - Contains consecutive underscores
* `.` or `..` - Reserved patterns
* Empty string - Must have at least 1 character
* Strings over 1500 characters - Exceeds maximum length

## JSON Properties Guide

The `props` field accepts any valid JSON structure. Here are common patterns:

### Simple Key-Value Properties
```hcl
props = jsonencode({
  location = "main-floor"
  active   = true
  priority = 1
})
```

### Nested Objects
```hcl
props = jsonencode({
  metadata = {
    created_by = "admin@company.com"
    created_at = "2024-01-15T10:30:00Z"
    version    = "1.0"
  }
  settings = {
    auto_sync    = true
    notifications = false
  }
})
```

### Arrays and Lists
```hcl
props = jsonencode({
  tags        = ["retail", "pos", "production"]
  permissions = ["read", "write", "admin"]
  endpoints   = [
    "https://api.company.com/endpoint1",
    "https://api.company.com/endpoint2"
  ]
})
```

### Mixed Data Types
```hcl
props = jsonencode({
  name          = "Store Manager Dashboard"
  employee_count = 25
  is_flagship   = true
  departments   = ["electronics", "clothing", "home"]
  location = {
    address = "123 Main St"
    city    = "New York"
    coordinates = {
      lat = 40.7128
      lng = -74.0060
    }
  }
})
```

## Error Handling

The provider includes comprehensive error handling for common scenarios:

### Authentication Errors
```
Error: Authentication Failed
The create operation failed for resource 'store:001'. Please check that:
• Your OAuth2 credentials are valid and not expired
• Your client has the necessary permissions
• The tenant ID is correct
```

### Permission Errors
```
Error: Permission Denied
You don't have permission to create resource 'store:001'. Please check that:
• Your OAuth2 token includes the required scopes (iam:read, iam:write)
• Your account has the necessary IAM permissions
• You're accessing the correct tenant
```

### Validation Errors
```
Error: Invalid Request
The create request for resource 'store/001' was invalid. This usually means:
• Resource ID follows the pattern (1-1500 chars, no slashes, no '.', '..', or '__')
• Resource name is not empty
• Props field contains valid JSON if provided
```

### Rate Limiting
```
Error: Rate Limit Exceeded
Rate limit exceeded for create operation on resource 'store:001'. Please:
• Wait before retrying the operation
• Reduce the frequency of API calls
• Contact support if the problem persists
```

### Service Errors
```
Error: Server Error
An internal server error occurred during create operation on resource 'store:001'. Please:
• Retry the operation after a short delay
• Check the HiiRetail service status
• Contact support if the problem persists
```

## Integration with IAM Roles

Resources are typically used in conjunction with custom IAM roles to create fine-grained access control:

```hcl
# Define the resource
resource "hiiretail_iam_resource" "store_001" {
  id   = "store:001"
  name = "Main Store - Downtown"
  props = jsonencode({
    location = "downtown"
    manager  = "jane.smith@company.com"
  })
}

# Create a custom role that references the resource
resource "hiiretail_iam_custom_role" "store_manager" {
  name        = "Store Manager - Store 001"
  description = "Manager permissions for specific store"
  permissions = [
    "inventory.${hiiretail_iam_resource.store_001.id}.read",
    "inventory.${hiiretail_iam_resource.store_001.id}.write",
    "sales.${hiiretail_iam_resource.store_001.id}.read",
    "reports.${hiiretail_iam_resource.store_001.id}.generate"
  ]
}

# Bind the role to a user or group
resource "hiiretail_iam_role_binding" "store_manager_binding" {
  role_name = hiiretail_iam_custom_role.store_manager.name
  subjects = [{
    type = "user"
    name = "jane.smith@company.com"
  }]
}
```

## Best Practices

### Resource Naming Convention
Use consistent, hierarchical naming patterns:
- **Stores**: `store:location-identifier` (e.g., `store:nyc-001`, `store:la-main`)
- **Departments**: `dept:department-name` (e.g., `dept:electronics`, `dept:home-garden`)
- **Applications**: `app:application-name` (e.g., `app:inventory-mgmt`, `app:pos-system`)
- **POS Systems**: `pos:store:terminal` (e.g., `pos:store-001:term-01`)

### Properties Organization
Structure your properties consistently:
```hcl
props = jsonencode({
  # Core metadata
  created_at = "2024-01-15T10:30:00Z"
  created_by = "admin@company.com"
  version    = "1.0"
  
  # Business properties
  location    = "main-floor"
  department  = "electronics"
  manager     = "john.doe@company.com"
  
  # Configuration
  settings = {
    auto_sync     = true
    notifications = true
    backup_enabled = true
  }
  
  # Access control hints (for documentation)
  permissions_hint = [
    "This resource requires 'iam.resource.read' permission to view",
    "This resource requires 'iam.resource.write' permission to modify"
  ]
})
```

### Resource Lifecycle
1. **Plan**: Use `terraform plan` to preview changes before applying
2. **Create**: Resources are created via `PUT /api/v1/tenants/{tenantId}/resources/{id}`
3. **Update**: Updates use the same `PUT` endpoint with modified data
4. **Delete**: Resources are deleted via `DELETE /api/v1/tenants/{tenantId}/resources/{id}`
5. **Import**: Existing resources can be imported using `terraform import`

### Security Considerations
- Use descriptive resource names that help with auditing
- Store sensitive data in properties only if necessary
- Regularly review resource permissions and usage
- Use consistent resource hierarchies for easier management

## Troubleshooting

### Common Issues

**Q: Resource creation fails with "Invalid Request" error**
A: Check your resource ID format. Ensure it doesn't contain `/`, `__`, or reserved patterns like `.` or `..`.

**Q: Import fails with "Resource Not Found" error**
A: Verify the resource ID exists in your tenant and you have read permissions.

**Q: JSON validation errors in props field**
A: Use `jsonencode()` function and ensure your JSON structure is valid. Test with online JSON validators.

**Q: Permission denied errors**
A: Ensure your OAuth2 token includes `iam:read` and `iam:write` scopes and your account has appropriate permissions.

**Q: Updates aren't reflected immediately**
A: Resources use eventual consistency. Allow a few seconds for changes to propagate across the system.

### Debugging

Enable detailed logging:
```hcl
provider "hiiretail" {
  # ... other configuration
  
  # Enable debug logging (in development only)
  log_level = "DEBUG"
}
```

Check resource state:
```shell
# Show current resource state
terraform show

# Validate configuration
terraform validate

# Plan changes
terraform plan

# Apply with detailed output
terraform apply -auto-approve
```

For additional support, contact the HiiRetail support team with:
- Your tenant ID
- Resource ID experiencing issues  
- Full error messages
- Terraform configuration (sanitized)