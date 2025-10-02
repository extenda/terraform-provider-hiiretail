# Research: Multi-API Provider User Experience Enhancement

## Research Objectives
This research phase explores the technical decisions needed to transform the current HiiRetail IAM provider into a multi-API provider that feels familiar to users of major cloud providers like GCP.

## Key Research Areas

### 1. Multi-API Provider Architecture Patterns

**Decision**: Service-based module organization within a single provider
**Rationale**: 
- GCP Terraform provider uses this pattern successfully with 200+ services
- AWS provider organizes resources by service (ec2, s3, iam, etc.)
- Enables consistent authentication across all APIs
- Simpler dependency management and versioning
- Users get familiar experience with unified provider configuration

**Alternatives considered**:
- Separate provider per API: Rejected due to authentication complexity and user confusion
- Flat resource organization: Rejected due to poor discoverability at scale

### 2. Naming Convention Standards

**Decision**: `hiiretail_{api}_{resource}` format
**Rationale**:
- Follows Terraform community conventions (aws_ec2_instance, google_compute_instance)
- Clear service grouping for resource discovery
- Consistent with user's specified requirements
- Enables auto-completion and IDE support

**Examples**:
- `hiiretail_iam_group` (current: hiiretail-iam_iam_group)  
- `hiiretail_iam_custom_role` (current: hiiretail-iam_custom_role)
- `hiiretail_iam_role_binding` (current: hiiretail-iam_iam_role_binding)
- `hiiretail_ccc_kind` (future API example)

**Alternatives considered**:
- Keeping current hyphenated format: Rejected due to inconsistency with major providers
- Using dots as separators: Rejected due to Terraform syntax limitations

### 3. Provider Registration and Discovery

**Decision**: Single provider registration with service-based resource organization
**Rationale**:
- Users install one provider (`terraform { required_providers { hiiretail = {...} } }`)
- Resources are automatically organized by service in documentation
- Consistent with major cloud providers
- Enables centralized configuration and authentication

**Migration Path**:
1. Rename provider from "hiiretail-iam" to "hiiretail"  
2. Update resource names to use underscore convention
3. Maintain backward compatibility aliases during transition
4. Provide migration guides and automation tools

### 4. Documentation Structure and User Experience

**Decision**: Service-grouped documentation with comprehensive examples
**Rationale**:
- Mirrors GCP provider documentation structure
- Service overview pages help users understand API relationships
- Getting started guides reduce time-to-first-success
- Comparison guides help users migrating from other providers

**Structure**:
```
docs/
├── guides/
│   ├── getting-started.md           # Basic setup and first resource
│   ├── authentication.md           # OAuth2 configuration
│   ├── provider-comparison.md      # vs AWS/GCP/Azure patterns
│   └── migration-guides/
│       └── from-hiiretail-iam.md   # Migration from old provider
├── resources/
│   ├── iam/
│   │   ├── overview.md             # IAM service introduction
│   │   ├── group.md                # hiiretail_iam_group
│   │   ├── custom_role.md          # hiiretail_iam_custom_role
│   │   └── role_binding.md         # hiiretail_iam_role_binding
│   └── ccc/ (future)
└── examples/
    ├── basic-iam-setup/
    ├── multi-service-deployment/
    └── enterprise-patterns/
```

### 5. Authentication Architecture

**Decision**: Centralized OAuth2 with service-specific endpoint support
**Rationale**:
- Maintains existing secure OAuth2 implementation
- Enables different APIs to have different endpoints if needed
- Single credential configuration for all services
- Consistent with major cloud providers

**Configuration Pattern**:
```hcl
provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
  
  # Service-specific endpoints (optional overrides)
  iam_endpoint = "https://iam-api.retailsvc.com"
  ccc_endpoint = "https://ccc-api.retailsvc.com"
}
```

### 6. Testing Strategy

**Decision**: Comprehensive test suite with service-based organization
**Rationale**:
- Terraform acceptance tests validate real provider behavior
- Service-specific test suites enable independent API testing
- Integration tests validate cross-service functionality
- Unit tests ensure provider logic correctness

**Test Organization**:
```
tests/
├── acceptance/
│   ├── iam/
│   │   ├── group_test.go
│   │   ├── custom_role_test.go
│   │   └── role_binding_test.go
│   └── provider_test.go             # Provider-level tests
├── integration/
│   ├── oauth2_test.go               # Authentication testing
│   └── multi_service_test.go        # Cross-service scenarios
└── unit/
    ├── provider/
    └── shared/
```

### 7. Migration Strategy

**Decision**: Phased migration with backward compatibility
**Rationale**:
- Minimizes disruption to existing users
- Provides clear migration path
- Enables gradual adoption of new patterns

**Migration Phases**:
1. **Phase 1**: Rename provider, update resource names, maintain aliases
2. **Phase 2**: Enhanced documentation and examples
3. **Phase 3**: Additional API service modules (CCC, etc.)
4. **Phase 4**: Deprecate old aliases, encourage migration

### 8. Code Generation Strategy

**Decision**: Enhanced generator configuration for multi-API support
**Rationale**:
- Maintains consistency with OpenAPI specs
- Enables rapid addition of new APIs
- Ensures compliance with Terraform conventions

**Generator Updates**:
- Update `generator_config.yaml` for new naming conventions
- Service-specific OpenAPI specs
- Enhanced validation and business logic injection

## Research Conclusions

All technical decisions support the primary goal of creating a familiar, scalable multi-API provider experience. The proposed architecture follows established patterns from major cloud providers while maintaining the security and functionality of the existing OAuth2 implementation.

The migration path ensures existing users can upgrade smoothly while new users benefit from the improved organization and documentation structure.

## Next Steps
Proceed to Phase 1 design to create detailed contracts and data models based on these research findings.