# Research Report: IAM Role Binding Resource Implementation

**Date**: September 30, 2025  
**Feature**: IAM Role Binding Resource Implementation and Testing  
**Phase**: 0 - Technical Research

## Research Questions Resolved

### 1. Terraform Plugin Framework Resource Patterns

**Decision**: Use terraform-plugin-framework v1.16.0 with existing provider patterns  
**Rationale**: 
- Provider already established with custom role and group resources following framework conventions
- Consistent schema definitions using framework's tfsdk types
- Proven OAuth2 integration and state management patterns exist
- Mock server infrastructure already implemented for testing

**Alternatives Considered**:
- terraform-plugin-sdk v2: Rejected due to project migration to newer framework
- Custom implementation: Rejected due to existing framework investment

### 2. Resource Schema Design Patterns

**Decision**: Follow generated schema from provider_code_spec.json with business logic enhancement  
**Rationale**:
- Generated schema provides proper terraform-plugin-framework structure
- Existing pattern from iam_custom_role resource demonstrates enhancement approach
- Schema already includes required fields: role_id, bindings, is_custom flag
- Max 10 bindings validation aligns with business requirements

**Technical Details**:
```go
// Schema structure (from generated code)
- role_id: schema.StringAttribute{Required: true}
- bindings: schema.ListAttribute{Required: true, MaxItems: 10}
- is_custom: schema.BoolAttribute{Optional: true, Default: false}
- tenant_id: schema.StringAttribute{Computed: true}  
- id: schema.StringAttribute{Computed: true}
```

### 3. CRUD Operation Implementation Strategy

**Decision**: Implement full CRUD with OAuth2 client integration following iam_custom_role patterns  
**Rationale**:
- Create: POST to role binding endpoint with tenant context
- Read: GET with ID and tenant isolation  
- Update: PUT with atomic binding updates
- Delete: DELETE with proper cleanup
- Import: Standard Terraform import with ID format validation

**Authentication Flow**:
- Reuse existing OAuth2 client from internal/client/
- Token refresh handling already implemented
- Tenant context passed in all API calls

### 4. Testing Strategy and Framework

**Decision**: Three-tier testing approach following Terraform provider best practices  
**Rationale**:
- Unit Tests: Test resource logic, validation, state management
- Integration Tests: Test with mock HTTP server (existing infrastructure)
- Acceptance Tests: Full Terraform lifecycle testing with terraform-plugin-testing

**Test Infrastructure**:
- Mock server: Reuse existing httptest.Server setup from acceptance_tests/
- Test fixtures: JSON response templates for role binding operations
- Provider factory: Use existing testProviderFactory from provider tests

### 5. Provider Registration and Integration

**Decision**: Register resource in main provider configuration following existing pattern  
**Rationale**:
- Add to Resources map in hiiretail_iam_provider_gen.go
- Follow naming convention: "hiiretail_iam_role_binding"
- Maintain consistency with iam_custom_role and iam_group resources

**Integration Points**:
```go
// In provider registration
"hiiretail_iam_role_binding": resourceIamRoleBinding.NewIamRoleBindingResource(),
```

### 6. Error Handling and Validation Patterns

**Decision**: Implement comprehensive error handling with Terraform diagnostics  
**Rationale**:
- HTTP error mapping to appropriate Terraform errors
- Validation errors for max bindings, role_id format
- Clear diagnostic messages for troubleshooting
- Retry logic for transient failures

**Error Categories**:
- Validation: Schema validation, business rule violations
- Authentication: OAuth2 token issues, permission errors  
- API: HTTP errors, timeout handling, retry logic
- State: Inconsistent state detection and recovery

### 7. Mock Server Test Data Strategy

**Decision**: Create comprehensive test fixtures covering all CRUD scenarios  
**Rationale**:
- JSON fixtures for create/read/update/delete responses
- Error scenario fixtures for negative testing
- Multiple tenant isolation test cases
- Boundary condition testing (max bindings, invalid IDs)

**Test Data Structure**:
```json
{
  "role_bindings": {
    "valid_create": {...},
    "max_bindings": {...},
    "invalid_role_id": {...},
    "tenant_isolation": {...}
  }
}
```

## Implementation Dependencies

### External Dependencies (Already Available)
- terraform-plugin-framework v1.16.0: Core framework
- terraform-plugin-testing: Acceptance test framework  
- golang.org/x/oauth2: OAuth2 client implementation
- httptest: Mock server infrastructure

### Internal Dependencies (Existing)
- internal/client/oauth2_client.go: Authentication client
- acceptance_tests/mock_server_test.go: Test infrastructure
- Generated schema from resource_iam_role_binding/iam_role_binding_resource_gen.go

### Missing Components (To Implement)
- Resource implementation file: iam_role_binding_resource.go
- Unit test file: iam_role_binding_resource_test.go  
- Provider registration in hiiretail_iam_provider_gen.go
- Acceptance test file: iam_role_binding_resource_test.go
- Mock server fixtures for role binding endpoints

## Risk Assessment

### Low Risk
- Schema generation and framework integration (established patterns)
- OAuth2 authentication (existing implementation)
- Mock server setup (infrastructure exists)

### Medium Risk  
- API endpoint behavior assumptions (mitigated by mock server testing)
- Terraform state management edge cases (addressed by comprehensive testing)

### High Risk
- None identified - following established provider patterns significantly reduces implementation risk

## Next Phase Requirements

**Phase 1 Prerequisites Met**:
- All technical unknowns resolved
- Implementation approach defined
- Testing strategy established
- Integration points identified

**Ready for Phase 1**: Design & Contracts generation with concrete technical decisions documented above.