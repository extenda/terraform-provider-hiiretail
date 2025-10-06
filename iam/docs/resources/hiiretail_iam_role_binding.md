---
subcategory: "IAM Management"
---

# hiiretail_iam_role_binding Resource

The `hiiretail_iam_role_binding` resource manages role bindings in the HiiRetail IAM system, allowing you to assign roles to users, groups, or service accounts.

## Example Usage

### Basic Role Binding

```hcl
resource "hiiretail_iam_role_binding" "example" {
  role_id = "custom-role-123"
  bindings = [
    {
      type = "user"
      id   = "user-456"
    }
  ]
}
```

### Multiple Bindings

```hcl
resource "hiiretail_iam_role_binding" "team_access" {
  role_id = "team-lead-role"
  bindings = [
    {
      type = "user"
      id   = "john.doe@company.com"
    },
    {
      type = "group"
      id   = "team-leads-group"
    },
    {
      type = "service_account"
      id   = "automation-sa-789"
    }
  ]
  description = "Team leads access permissions"
}
```

### System Role Binding

```hcl
resource "hiiretail_iam_role_binding" "admin_access" {
  role_id = "system-admin"
  bindings = [
    {
      type = "user"
      id   = "admin@company.com"
    }
  ]
  description = "System administrator access"
}
```

## Argument Reference

The following arguments are supported:

- `role_id` - (Required) The ID of the role to bind. Can be either a custom role ID or a system role name.
- `bindings` - (Required) List of binding objects. Maximum of 10 bindings per role binding. Each binding contains:
  - `type` - (Required) The type of entity being bound. Must be one of: `user`, `group`, or `service_account`.
  - `id` - (Required) The ID or identifier of the entity being bound.
- `description` - (Optional) A description for this role binding.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for this role binding.
- `tenant_id` - The tenant ID associated with this role binding.
- `status` - The current status of the role binding.
- `created_at` - The timestamp when the role binding was created.
- `updated_at` - The timestamp when the role binding was last updated.

## Validation Rules

The resource enforces several validation rules:

1. **Maximum Bindings**: A role binding can have at most 10 bindings.
2. **Tenant Isolation**: All bindings are automatically isolated to the provider's configured tenant.
3. **Binding Format**: Each binding must have a valid `type` and non-empty `id`.
4. **Role Requirements**: The `role_id` must be specified and non-empty.

## Import

Role bindings can be imported using their ID:

```shell
terraform import hiiretail_iam_role_binding.example rb-12345678-1234-1234-1234-123456789012
```

## Notes

- Role bindings are tenant-isolated and will only work within the configured tenant context.
- Changes to the `role_id` will force recreation of the resource.
- The resource supports both custom roles (created via `hiiretail_iam_custom_role`) and system-defined roles.
- Binding order is preserved but not semantically significant.

## Error Handling

The resource provides detailed error messages for common issues:

- Invalid binding formats
- Exceeding maximum binding limits
- Non-existent roles or entities
- Permission denied scenarios
- Network connectivity issues

All errors include context to help with troubleshooting and resolution.