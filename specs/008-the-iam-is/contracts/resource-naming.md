# Resource Naming Contract

## Naming Convention Contract

**Pattern**: `hiiretail_{service}_{resource_type}`

**Service Identifiers**:
- `iam`: Identity and Access Management API
- `ccc`: Customer Care Center API (future)
- Additional services follow lowercase, alphanumeric identifiers

**Resource Type Identifiers**:
- Use singular form: `group`, `role`, `binding` (not `groups`, `roles`, `bindings`)
- Use snake_case for multi-word types: `custom_role`, `role_binding`
- Keep names concise but descriptive

## Current Resource Mappings

**IAM Service Resources**:
```hcl
# Old → New
hiiretail-iam_iam_group      → hiiretail_iam_group
hiiretail-iam_custom_role    → hiiretail_iam_custom_role  
hiiretail-iam_iam_role_binding → hiiretail_iam_role_binding
```

**Data Sources**:
```hcl
# Follow same pattern with plural forms where appropriate
hiiretail_iam_groups   # List multiple groups
hiiretail_iam_roles    # List available roles
```

## Future Service Resources

**CCC Service Resources**:
```hcl
hiiretail_ccc_kind           # Example future resource
hiiretail_ccc_workflow       # Another example
```

## Registry and Discovery Contract

**Provider Registry**:
- Provider name: `extenda/hiiretail` (changed from `extenda/hiiretail-iam`)
- Single provider manages all HiiRetail services
- Version compatibility maintained across all services

**Resource Discovery**:
- Resources grouped by service in documentation
- Auto-completion support in IDEs follows service.resource pattern
- Error messages reference correct resource names

**Migration Support**:
- Backward compatibility aliases during transition period
- Clear migration path documentation
- Automated migration tooling where possible