# API Issue Report: V2 Group Role Endpoints Inconsistency

**Date**: October 3, 2025  
**Reporter**: Terraform Provider Development Team  
**Priority**: High  
**API Version**: V2  

## Issue Summary

There is a critical inconsistency in the V2 Group Role endpoints that prevents proper state management in infrastructure-as-code tools like Terraform. The GET endpoint does not return role bindings that were successfully created via POST operations.

## Affected Endpoint

**Base Endpoint**: `/api/v2/tenants/{tenantId}/groups/{groupId}/roles`

- **POST**: ✅ Works correctly - creates role bindings successfully
- **GET**: ❌ Broken - returns empty array despite successful creations

## Detailed Problem Description

### Current Behavior (Incorrect)

1. **Create Role Binding** (POST):
   ```http
   POST /api/v2/tenants/your-tenant-id/groups/EYNaCiYX6WFmoPxXCGMf/roles
   Content-Type: application/json
   
   {
     "roleId": "custom-roles/custom.TerraformTestShayne",
     "isCustom": true,
     "bindings": ["EYNaCiYX6WFmoPxXCGMf"]
   }
   ```
   **Response**: HTTP 200/201 ✅ Success

2. **Retrieve Role Bindings** (GET):
   ```http
   GET /api/v2/tenants/your-tenant-id/groups/EYNaCiYX6WFmoPxXCGMf/roles
   ```
   **Response**: HTTP 200 with empty array `[]` ❌ **Should contain the created role binding**

### Expected Behavior (Correct)

After a successful POST operation, the GET endpoint should return:

```json
[
  {
    "isCustom": true,
    "roleId": "custom.TerraformTestShayne",
    "bindings": ["EYNaCiYX6WFmoPxXCGMf"],
    "fixedBindings": []
  }
]
```

## Evidence & Testing

### Test Case
- **Tenant ID**: `your-tenant-id`
- **Group ID**: `EYNaCiYX6WFmoPxXCGMf`
- **Role ID**: `custom-roles/custom.TerraformTestShayne`

### API Call Logs
```
DEBUG: POST request successful - HTTP 200
DEBUG: Role binding created with payload:
  roleId: custom-roles/custom.TerraformTestShayne
  isCustom: true
  bindings: [EYNaCiYX6WFmoPxXCGMf]

DEBUG: GET request - HTTP 200, Body length: 2
DEBUG: API response for group EYNaCiYX6WFmoPxXCGMf (found 0 role bindings): []
```

## Business Impact

### Infrastructure Management
- **Terraform Provider**: Cannot achieve idempotency
- **State Drift**: Every `terraform plan` shows phantom changes
- **CI/CD Pipelines**: Broken automation workflows
- **User Experience**: Confusing behavior, loss of confidence

### Technical Impact
- **CRUD Operations**: Incomplete - Create works, Read fails
- **Data Consistency**: POST and GET operations are out of sync
- **API Contract**: Violates RESTful principles

## Current Workaround

We have implemented a temporary workaround in our Terraform provider that constructs expected responses when the API returns empty results. However, this is:

- ❌ **Not sustainable** long-term
- ❌ **Error-prone** - based on assumptions
- ❌ **Masks the real issue** - doesn't fix the API

## Requests for API Development Team

### Immediate Actions Required

1. **🔍 Investigate Root Cause**
   - Why don't GET operations return role bindings created via POST?
   - Is there a data persistence issue?
   - Are there database consistency problems?

2. **🛠️ Fix Data Consistency**
   - Ensure POST operations properly persist data
   - Verify GET operations query the correct data source
   - Test the complete CRUD cycle

3. **✅ Verify Response Format**
   - Confirm `RoleBindingDto` matches OpenAPI specification
   - Validate field names and data types
   - Check array structure consistency

### Quality Assurance

4. **🧪 Add Integration Tests**
   - Test POST → GET cycle
   - Verify data persistence across operations
   - Prevent regression of this issue

5. **📋 Test Complete CRUD Operations**
   - **C**reate: POST role binding
   - **R**ead: GET role binding (currently broken)
   - **U**pdate: PUT/PATCH role binding
   - **D**elete: DELETE role binding

## Schema References

### CreateRoleBindingDto (POST Request)
```json
{
  "roleId": "string",
  "isCustom": "boolean",
  "bindings": ["string"]
}
```

### RoleBindingDto (GET Response)
```json
{
  "isCustom": "boolean",
  "roleId": "string", 
  "bindings": ["string"],
  "fixedBindings": ["string"]
}
```

## Next Steps

1. **Acknowledgment**: Please confirm receipt and priority assignment
2. **Timeline**: Provide estimated fix timeline
3. **Communication**: Keep us updated on investigation progress
4. **Testing**: We can assist with testing once a fix is deployed

## Contact Information

**Team**: HiiRetail Terraform Provider Development  
**Context**: Infrastructure as Code implementation  
**Urgency**: Blocking production automation workflows  

---

**Note**: This issue has been thoroughly investigated from the client side. The problem is confirmed to be API-side data consistency between POST creation and GET retrieval operations.