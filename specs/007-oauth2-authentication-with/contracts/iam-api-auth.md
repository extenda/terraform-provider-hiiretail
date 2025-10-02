# IAM API Authentication Contract
**Endpoints**: 
- Live: `https://iam-api.retailsvc.com/*`
- Test: `https://iam-api.retailsvc-test.com/*`

**Purpose**: Authenticated requests to HiiRetail IAM API

## Request Contract

### Headers
```
Authorization: Bearer {access_token}
Content-Type: application/json
Accept: application/json
X-Tenant-ID: {tenant_id}
```

### Header Schema
```json
{
  "type": "object",
  "properties": {
    "Authorization": {
      "type": "string",
      "pattern": "^Bearer [A-Za-z0-9\\-_=]+\\.[A-Za-z0-9\\-_=]+\\.[A-Za-z0-9\\-_=]+$",
      "description": "OAuth2 Bearer token"
    },
    "Content-Type": {
      "type": "string",
      "enum": ["application/json"],
      "description": "Request content type"
    },
    "Accept": {
      "type": "string",
      "enum": ["application/json"],
      "description": "Accepted response type"
    },
    "X-Tenant-ID": {
      "type": "string",
      "minLength": 1,
      "description": "HiiRetail tenant identifier"
    }
  },
  "required": ["Authorization", "X-Tenant-ID"]
}
```

## Response Contract

### Success Response (200 OK)
```json
{
  "data": {
    // Resource-specific data
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2025-10-01T12:00:00Z"
  }
}
```

### Authentication Error (401 Unauthorized)
```json
{
  "error": {
    "code": "AUTHENTICATION_FAILED",
    "message": "Invalid or expired access token",
    "request_id": "req_123456789"
  }
}
```

### Authorization Error (403 Forbidden)
```json
{
  "error": {
    "code": "INSUFFICIENT_PERMISSIONS",
    "message": "Token does not have required permissions for this operation",
    "request_id": "req_123456789"
  }
}
```

### Tenant Error (404 Not Found)
```json
{
  "error": {
    "code": "TENANT_NOT_FOUND",
    "message": "Specified tenant ID not found or not accessible",
    "request_id": "req_123456789"
  }
}
```

### Error Schema
```json
{
  "type": "object",
  "properties": {
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "string",
          "enum": [
            "AUTHENTICATION_FAILED",
            "INSUFFICIENT_PERMISSIONS", 
            "TENANT_NOT_FOUND",
            "INVALID_REQUEST",
            "INTERNAL_ERROR"
          ],
          "description": "Machine-readable error code"
        },
        "message": {
          "type": "string",
          "description": "Human-readable error message"
        },
        "request_id": {
          "type": "string",
          "description": "Unique request identifier for debugging"
        }
      },
      "required": ["code", "message", "request_id"]
    }
  },
  "required": ["error"]
}
```

## Contract Tests Required

1. **Valid Bearer token with valid tenant** → 200 OK with resource data
2. **Missing Authorization header** → 401 Unauthorized
3. **Invalid Bearer token format** → 401 Unauthorized
4. **Expired Bearer token** → 401 Unauthorized
5. **Valid token with invalid tenant** → 404 Not Found
6. **Valid token with insufficient permissions** → 403 Forbidden
7. **Missing X-Tenant-ID header** → 400 Bad Request

## Endpoint Resolution Tests Required

1. **Live tenant ID pattern** → Routes to iam-api.retailsvc.com
2. **Test tenant ID pattern** → Routes to iam-api.retailsvc-test.com
3. **Mock mode enabled** → Routes to mock server URL
4. **Environment variable override** → Routes to override URL
5. **Invalid tenant ID format** → Error with clear message

## Security Requirements

- **TLS 1.2+** required for all requests
- **Bearer tokens** must not appear in logs or error messages  
- **Token validation** on every request
- **Tenant isolation** enforced at API level
- **Request rate limiting** may apply