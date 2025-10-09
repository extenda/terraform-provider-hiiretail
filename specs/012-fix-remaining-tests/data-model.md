# Data Model: Fix Remaining Tests

## Entities

### Resource Models
- **RoleBindingResourceModel**
  - Fields: id, tenant_id, roles, bindings, is_custom
- **IamCustomRoleModel**
  - Fields: id, tenant_id, permissions, description
- **IamResourceModel**
  - Fields: id, tenant_id, name, type, attributes

### Provider Model
- **HiiRetailProviderModel**
  - Fields: base_url, ccc_endpoint, client_id, client_secret, iam_endpoint, max_retries, scopes, timeout_seconds, token_url, tenant_id

## Relationships
- RoleBindingResourceModel references RoleModel
- IamCustomRoleModel references permissions
- All resource models reference tenant_id for multi-tenancy

## Validation Rules
- All models must have tenant_id field present and correctly mapped
- All provider configuration tests must validate required fields and error messages
- All resource contract tests must initialize IAMClient and other dependencies

## State Transitions
- Resource models must support Create, Read, Update, Delete operations
- Provider model must support configuration validation and error handling

---
