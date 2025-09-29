# API Contract: Custom Role CRUD Operations

**Date**: September 28, 2025  
**API Base**: `{base_url}/api/v1/tenants/{tenant_id}`  
**Authentication**: OAuth2 Bearer Token

## Create Custom Role

**Endpoint**: `POST /custom-roles`

**Request Schema**:
```json
{
  "type": "object",
  "required": ["id", "permissions"],
  "properties": {
    "id": {
      "type": "string",
      "description": "Unique identifier for the custom role"
    },
    "name": {
      "type": "string",
      "minLength": 3,
      "maxLength": 256,
      "description": "Human-readable name for the role"
    },
    "permissions": {
      "type": "array",
      "minItems": 1,
      "maxItems": 500,
      "items": {
        "$ref": "#/definitions/Permission"
      }
    }
  }
}
```

**Response Schema** (201 Created):
```json
{
  "type": "object",
  "properties": {
    "id": {"type": "string"},
    "name": {"type": "string"},
    "tenant_id": {"type": "string"},
    "permissions": {
      "type": "array",
      "items": {"$ref": "#/definitions/PermissionResponse"}
    },
    "created_at": {"type": "string", "format": "date-time"},
    "updated_at": {"type": "string", "format": "date-time"}
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid permission format, constraint violations
- `401 Unauthorized`: Invalid or expired OAuth token
- `403 Forbidden`: Insufficient permissions for tenant
- `409 Conflict`: Role ID already exists
- `422 Unprocessable Entity`: Permission limit exceeded

## Read Custom Role

**Endpoint**: `GET /custom-roles/{role_id}`

**Response Schema** (200 OK):
```json
{
  "type": "object",
  "properties": {
    "id": {"type": "string"},
    "name": {"type": "string"},
    "tenant_id": {"type": "string"},
    "permissions": {
      "type": "array",
      "items": {"$ref": "#/definitions/PermissionResponse"}
    },
    "created_at": {"type": "string", "format": "date-time"},
    "updated_at": {"type": "string", "format": "date-time"}
  }
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or expired OAuth token
- `403 Forbidden`: Insufficient permissions for tenant
- `404 Not Found`: Role does not exist

## Update Custom Role

**Endpoint**: `PUT /custom-roles/{role_id}`

**Request Schema**:
```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "minLength": 3,
      "maxLength": 256
    },
    "permissions": {
      "type": "array",
      "minItems": 1,
      "maxItems": 500,
      "items": {"$ref": "#/definitions/Permission"}
    }
  }
}
```

**Response Schema** (200 OK):
```json
{
  "type": "object",
  "properties": {
    "id": {"type": "string"},
    "name": {"type": "string"},
    "tenant_id": {"type": "string"},
    "permissions": {
      "type": "array",
      "items": {"$ref": "#/definitions/PermissionResponse"}
    },
    "updated_at": {"type": "string", "format": "date-time"}
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid permission format, constraint violations
- `401 Unauthorized`: Invalid or expired OAuth token
- `403 Forbidden`: Insufficient permissions for tenant
- `404 Not Found`: Role does not exist
- `422 Unprocessable Entity`: Permission limit exceeded

## Delete Custom Role

**Endpoint**: `DELETE /custom-roles/{role_id}`

**Response**: `204 No Content`

**Error Responses**:
- `401 Unauthorized`: Invalid or expired OAuth token
- `403 Forbidden`: Insufficient permissions for tenant
- `404 Not Found`: Role does not exist
- `409 Conflict`: Role still referenced by other resources

## Schema Definitions

### Permission (Request)
```json
{
  "type": "object",
  "required": ["id"],
  "properties": {
    "id": {
      "type": "string",
      "pattern": "^[a-z][-a-z]{2}\\.[a-z][-a-z]{1,15}\\.[a-z][-a-z]{1,15}$",
      "description": "Permission ID following {systemPrefix}.{resource}.{action} pattern"
    },
    "attributes": {
      "type": "object",
      "maxProperties": 10,
      "additionalProperties": {
        "type": "string",
        "maxLength": 256
      },
      "description": "Optional key-value attributes, max 10 props, keys ≤40 chars, values ≤256 chars"
    }
  }
}
```

### PermissionResponse (Response)
```json
{
  "type": "object",
  "properties": {
    "id": {"type": "string"},
    "alias": {
      "type": "string",
      "description": "Server-computed alias for the permission"
    },
    "attributes": {
      "type": "object",
      "maxProperties": 10,
      "additionalProperties": {"type": "string"}
    }
  }
}
```

## Business Rules

### Permission Limits
- General permissions: maximum 100 per role
- POS permissions (id starts with "pos."): maximum 500 per role
- Mixed roles: POS limit applies to pos.* permissions, general limit to others

### Validation Rules
- Permission ID pattern must be strictly enforced
- Attribute key length validation (≤ 40 characters)  
- Attribute value length validation (≤ 256 characters)
- Maximum 10 attribute properties per permission

### Error Handling
- Detailed validation error messages for each constraint violation
- Specific error codes for different failure scenarios
- Proper HTTP status codes following REST conventions
- OAuth2 token refresh handling for expired tokens