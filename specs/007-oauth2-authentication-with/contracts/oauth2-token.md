# OAuth2 Token Endpoint Contract
**Endpoint**: `POST https://auth.retailsvc.com/oauth/token`  
**Purpose**: OAuth2 client credentials token acquisition

## Request Contract

### Headers
```
Content-Type: application/x-www-form-urlencoded
Accept: application/json
```

### Request Body (form-encoded)
```
grant_type=client_credentials
client_id={client_id}
client_secret={client_secret}
scope=hiiretail:iam
```

### Request Schema
```json
{
  "type": "object",
  "properties": {
    "grant_type": {
      "type": "string",
      "enum": ["client_credentials"],
      "description": "OAuth2 grant type"
    },
    "client_id": {
      "type": "string",
      "minLength": 1,
      "description": "OAuth2 client identifier"
    },
    "client_secret": {
      "type": "string",
      "minLength": 1,
      "description": "OAuth2 client secret"
    },
    "scope": {
      "type": "string",
      "enum": ["hiiretail:iam"],
      "description": "Requested access scope"
    }
  },
  "required": ["grant_type", "client_id", "client_secret", "scope"]
}
```

## Response Contract

### Success Response (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "hiiretail:iam"
}
```

### Success Schema
```json
{
  "type": "object",
  "properties": {
    "access_token": {
      "type": "string",
      "minLength": 1,
      "description": "OAuth2 access token"
    },
    "token_type": {
      "type": "string",
      "enum": ["Bearer"],
      "description": "Token type"
    },
    "expires_in": {
      "type": "integer",
      "minimum": 1,
      "description": "Token lifetime in seconds"
    },
    "scope": {
      "type": "string",
      "description": "Granted access scope"
    }
  },
  "required": ["access_token", "token_type", "expires_in"]
}
```

### Error Response (400 Bad Request)
```json
{
  "error": "invalid_request",
  "error_description": "Missing required parameter: client_id"
}
```

### Error Response (401 Unauthorized)
```json
{
  "error": "invalid_client",
  "error_description": "Invalid client credentials"
}
```

### Error Schema
```json
{
  "type": "object",
  "properties": {
    "error": {
      "type": "string",
      "enum": ["invalid_request", "invalid_client", "invalid_grant", "unsupported_grant_type"],
      "description": "OAuth2 error code"
    },
    "error_description": {
      "type": "string",
      "description": "Human-readable error description"
    }
  },
  "required": ["error"]
}
```

## Contract Tests Required

1. **Valid client credentials** → 200 OK with valid token
2. **Missing client_id** → 400 Bad Request with invalid_request error
3. **Missing client_secret** → 400 Bad Request with invalid_request error
4. **Invalid client credentials** → 401 Unauthorized with invalid_client error
5. **Wrong grant_type** → 400 Bad Request with unsupported_grant_type error
6. **Wrong scope** → 400 Bad Request with invalid_scope error

## Security Requirements

- **TLS 1.2+** required for all requests
- **Client credentials** must not appear in logs or error messages
- **Rate limiting** may apply (429 Too Many Requests response)
- **Token validation** must verify signature and expiration