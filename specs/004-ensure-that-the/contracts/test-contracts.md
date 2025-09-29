# Test Contracts: IAM Custom Role Resource

**Date**: September 28, 2025  
**Purpose**: Contract tests that must fail initially (no implementation yet)

## Contract Test Categories

### 1. Schema Validation Contracts

**Test**: `TestCustomRoleSchema_RequiredFields`
**Contract**: Resource schema must enforce required fields (id, permissions)
**Expected**: FAIL (schema validation not implemented)

**Test**: `TestCustomRoleSchema_OptionalFields`  
**Contract**: Resource schema must handle optional fields (name, tenant_id) properly
**Expected**: FAIL (optional field handling not implemented)

**Test**: `TestCustomRoleSchema_PermissionValidation`
**Contract**: Permission objects must validate ID pattern and attributes constraints
**Expected**: FAIL (permission validation not implemented)

### 2. CRUD Operation Contracts

**Test**: `TestCustomRoleCreate_ValidRole`
**Contract**: Valid role configuration must create resource successfully
**Expected**: FAIL (Create operation not implemented)

**Test**: `TestCustomRoleRead_ExistingRole`
**Contract**: Existing role must be readable with all attributes populated
**Expected**: FAIL (Read operation not implemented)

**Test**: `TestCustomRoleUpdate_PermissionChanges`
**Contract**: Permission list updates must modify resource state correctly
**Expected**: FAIL (Update operation not implemented)

**Test**: `TestCustomRoleDelete_RemoveRole`
**Contract**: Role deletion must remove resource and clean up state
**Expected**: FAIL (Delete operation not implemented)

### 3. Validation Contracts

**Test**: `TestPermissionPattern_ValidFormats`
**Contract**: Valid permission patterns must pass validation
**Test Cases**:
- `pos.payment.create` → PASS
- `sys.user.manage` → PASS  
- `abc.resource.action` → PASS

**Test**: `TestPermissionPattern_InvalidFormats`
**Contract**: Invalid permission patterns must fail validation
**Test Cases**:
- `invalid-format` → FAIL
- `too.short.x` → FAIL
- `ab.toolongresourcename12345.action` → FAIL
- `123.numeric.start` → FAIL

**Test**: `TestPermissionLimits_GeneralPermissions`
**Contract**: General permissions must enforce 100-item limit
**Expected**: FAIL (limit validation not implemented)

**Test**: `TestPermissionLimits_POSPermissions`
**Contract**: POS permissions must allow up to 500 items
**Expected**: FAIL (POS-specific validation not implemented)

**Test**: `TestAttributeConstraints_SizeLimits`
**Contract**: Attribute objects must enforce size constraints
**Test Cases**:
- 10 properties max → FAIL if exceeded
- 40 char key limit → FAIL if exceeded  
- 256 char value limit → FAIL if exceeded

### 4. Provider Integration Contracts

**Test**: `TestProviderRegistration_CustomRoleResource`
**Contract**: Custom role resource must be registered with provider
**Expected**: FAIL (resource not registered in provider.go)

**Test**: `TestOAuth2Integration_Authentication`
**Contract**: Custom role operations must use OAuth2 client from provider
**Expected**: FAIL (OAuth2 integration not implemented)

**Test**: `TestTenantContext_Inheritance`
**Contract**: Custom roles must inherit tenant context from provider
**Expected**: FAIL (tenant context not implemented)

### 5. Error Handling Contracts

**Test**: `TestAPIErrors_HTTPStatusMapping`
**Contract**: API errors must map to appropriate Terraform diagnostics
**Test Cases**:
- 400 Bad Request → ValidationError
- 401 Unauthorized → AuthenticationError
- 403 Forbidden → AuthorizationError
- 404 Not Found → NotFoundError
- 409 Conflict → ConflictError
- 422 Unprocessable → ValidationError

**Test**: `TestRetryLogic_TransientFailures`
**Contract**: Transient failures must trigger retry with exponential backoff
**Expected**: FAIL (retry logic not implemented)

**Test**: `TestConcurrency_StateConsistency`
**Contract**: Concurrent operations must maintain state consistency
**Expected**: FAIL (concurrency handling not implemented)

### 6. Performance Contracts

**Test**: `TestPerformance_MaximumPermissions`
**Contract**: Operations with 500 permissions must complete within reasonable time
**Expected**: FAIL (performance optimization not implemented)

**Test**: `TestMemoryUsage_LargeRoles`
**Contract**: Large roles must not exceed memory usage limits
**Expected**: FAIL (memory optimization not implemented)

## Mock Server Contracts

### API Endpoint Simulation

**Contract**: Mock server must simulate all CRUD endpoints
**Endpoints Required**:
- `POST /custom-roles` → 201 Created
- `GET /custom-roles/{id}` → 200 OK
- `PUT /custom-roles/{id}` → 200 OK  
- `DELETE /custom-roles/{id}` → 204 No Content

**Contract**: Mock server must simulate error scenarios
**Error Scenarios**:
- Invalid permission format → 400 Bad Request
- Permission limit exceeded → 422 Unprocessable Entity
- Role not found → 404 Not Found
- Authentication failure → 401 Unauthorized

**Contract**: Mock server must validate request schemas
**Validations**:
- Required field presence
- Permission pattern matching
- Attribute constraint enforcement
- Permission count limits

## Contract Test Execution Strategy

### Test-Driven Development Flow
1. **Write failing contract tests** (this phase)
2. **Implement minimal resource structure** to make schema tests pass
3. **Implement CRUD operations** to make operation tests pass
4. **Add validation logic** to make validation tests pass
5. **Integrate with provider** to make integration tests pass
6. **Add error handling** to make error tests pass
7. **Optimize performance** to make performance tests pass

### Contract Verification
- All contract tests must initially fail (proving no implementation)
- Each implementation phase should make specific contract test groups pass
- Final implementation should make all contract tests pass
- Contract tests serve as acceptance criteria for feature completion

### CI/CD Integration
- Contract tests run in every build
- Contract failures block deployment
- Contract test coverage metrics tracked
- Contract test performance monitored

This contract suite ensures comprehensive validation of the custom role resource implementation while following test-driven development principles.