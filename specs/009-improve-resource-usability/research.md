# Research: Improve Resource Usability

**Phase**: 0 | **Date**: October 2, 2025  
**Feature**: Improve Resource Usability

## Research Questions Resolved

### 1. Specific Usability Pain Points
**Decision**: Focus on validation, error messages, and documentation based on Terraform provider best practices

**Rationale**: 
- Analysis of simple_test.tf reveals complex resource configurations (role bindings with references, permission arrays)
- Terraform providers commonly struggle with unclear validation errors and poor cross-resource reference validation
- HashiCorp's provider development documentation emphasizes clear error messages and comprehensive validation

**Alternatives considered**:
- UI/UX improvements (N/A for Terraform providers)
- Performance optimizations (not the primary usability concern)
- API design changes (out of scope - provider layer improvements only)

### 2. Resource-Specific Usability Challenges
**Decision**: Prioritize role binding complexity, permission validation, and resource reference integrity

**Rationale**:
- **hiiretail_iam_role_binding**: Most complex resource with role references, member arrays, and optional conditions
- **hiiretail_iam_custom_role**: Permission strings require validation against known patterns
- **Cross-resource references**: Role bindings reference groups and roles that may not exist
- **Data source integration**: Users need clear guidance on querying existing resources

**Alternatives considered**:
- Focusing only on simple resources (insufficient - complex resources drive most usability issues)
- API-level improvements (beyond provider scope)
- Complete resource redesign (breaking change, violates backward compatibility)

### 3. Enterprise Setup Patterns
**Decision**: Support common IAM patterns based on industry best practices and simple_test.tf examples

**Rationale**:
- **Basic setup**: Groups → Custom roles → Role bindings (as shown in test file)
- **Enterprise patterns**: Service accounts, conditional access, permission hierarchies
- **Multi-tenant scenarios**: Tenant-specific configurations with proper isolation
- **Compliance requirements**: Audit trails, principle of least privilege

**Alternatives considered**:
- Domain-specific templates (too narrow)
- Generic examples (insufficient for real-world usage)
- API-generated templates (requires additional tooling development)

## Technology Research

### Terraform Plugin Framework Validation
**Best Practices Found**:
- Use `schema.StringValidator` and custom validators for complex fields
- Implement `plan_modifier.RequiresReplace()` for immutable fields
- Use `path.Expressions` for detailed error location reporting
- Implement cross-field validation in resource-level `Validate` methods

**Error Message Standards**:
- Include field path, current value, expected format, and example
- Use consistent terminology across all resources
- Provide actionable guidance, not just error descriptions
- Support internationalization where applicable

### Resource Reference Validation
**Patterns Identified**:
- Plan-time validation for resource existence checking
- State-based validation for dependency ordering
- API-level validation as final safety net
- Graceful degradation when validation services unavailable

**Implementation Approaches**:
- Use Terraform's built-in dependency resolution where possible
- Implement custom validators for complex reference patterns
- Provide suggestions for common typos and naming patterns
- Support both resource references and direct string values

### Documentation Standards
**HashiCorp Requirements**:
- Complete resource schema documentation
- Working examples for common use cases
- Import documentation with specific syntax
- Troubleshooting guides for common error scenarios

**Content Structure**:
- Resource overview with business context
- Argument reference with validation rules
- Attribute reference with output descriptions
- Example configurations for basic and advanced scenarios
- Migration guides for breaking changes

## Implementation Patterns

### Validation Architecture
```go
// Multi-level validation approach
1. Schema-level validation (format, type, required fields)
2. Resource-level validation (cross-field logic, business rules)
3. Plan-time validation (resource existence, references)
4. API-level validation (server-side constraints)
```

### Error Message Template
```
Field '{field_path}' validation failed: {specific_issue}
Current value: '{current_value}'
Expected format: {format_description}
Example: {working_example}
Additional context: {helpful_guidance}
```

### Reference Resolution Strategy
```go
// For role bindings referencing groups and roles
1. Parse reference format (group:name, role:name, user:email)
2. Validate reference syntax before API calls
3. Resolve references during plan phase where possible
4. Provide clear error messages for unresolvable references
```

## Success Metrics

### Measurable Improvements
- **Error Clarity**: Error messages include specific field names, expected formats, and examples
- **Validation Coverage**: 100% of resource fields have appropriate validation rules
- **Documentation Completeness**: All resources have working examples and troubleshooting guides
- **Reference Integrity**: Cross-resource references validated at plan time where possible

### User Experience Indicators
- **Reduced Trial-and-Error**: Users can configure resources correctly on first attempt with good examples
- **Faster Troubleshooting**: Error messages provide actionable guidance for resolution
- **Self-Service Capability**: Documentation enables users to solve common issues independently
- **Configuration Confidence**: Users understand the impact of their configurations before applying

---

**Research Status**: COMPLETE ✅  
**All NEEDS CLARIFICATION resolved**: YES ✅  
**Ready for Phase 1**: YES ✅