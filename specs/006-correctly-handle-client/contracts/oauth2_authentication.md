# OAuth2 Authentication Contract

**Service**: Hii Retail OAuth Client Management Service (OCMS)  
**Base URL**: https://auth.retailsvc.com  
**Discovery**: https://auth.retailsvc.com/.well-known/openid-configuration

## Discovery Endpoint Contract

### GET /.well-known/openid-configuration

**Request**:
```http
GET https://auth.retailsvc.com/.well-known/openid-configuration
Accept: application/json
```

**Response** (Expected):
```json
{
  "issuer": "https://auth.retailsvc.com",
  "token_endpoint": "https://auth.retailsvc.com/oauth2/token",
  "authorization_endpoint": "https://auth.retailsvc.com/oauth2/authorize",
  "jwks_uri": "https://auth.retailsvc.com/.well-known/jwks.json",
  "grant_types_supported": [
    "client_credentials",
    "authorization_code",
    "refresh_token"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "client_secret_post"
  ],
  "response_types_supported": [
    "code",
    "token"
  ],
  "scopes_supported": [
    "iam:read",
    "iam:write",
    "iam:admin"
  ]
}
```

**Error Responses**:
- `404 Not Found`: Discovery endpoint not available
- `5xx Server Error`: OCMS service unavailable

## Token Endpoint Contract

### POST /oauth2/token (Client Credentials)

**Request**:
```http
POST https://auth.retailsvc.com/oauth2/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic <base64(client_id:client_secret)>

grant_type=client_credentials&scope=iam:read+iam:write
```

**Alternative Request** (POST body credentials):
```http
POST https://auth.retailsvc.com/oauth2/token
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&client_id=<client_id>&client_secret=<client_secret>&scope=iam:read+iam:write
```

**Success Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "iam:read iam:write"
}
```

**Error Responses**:

**400 Bad Request** (Invalid request):
```json
{
  "error": "invalid_request",
  "error_description": "Missing required parameter: grant_type"
}
```

**400 Bad Request** (Unsupported grant):
```json
{
  "error": "unsupported_grant_type",
  "error_description": "The authorization grant type is not supported"
}
```

**401 Unauthorized** (Invalid credentials):
```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

**403 Forbidden** (Insufficient scope):
```json
{
  "error": "invalid_scope",
  "error_description": "The requested scope is invalid, unknown, or malformed"
}
```

**429 Too Many Requests**:
```json
{
  "error": "rate_limited",
  "error_description": "Too many token requests",
  "retry_after": 60
}
```

**500 Server Error**:
```json
{
  "error": "server_error",
  "error_description": "The authorization server encountered an unexpected condition"
}
```

## API Authentication Contract

### Using Bearer Token

**Request**:
```http
GET https://iam-api.retailsvc-test.com/iam/v1/groups
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
X-Tenant-ID: <tenant_id>
```

**Token Expired Response** (401 Unauthorized):
```json
{
  "error": "invalid_token",
  "error_description": "The access token expired"
}
```

**Token Invalid Response** (401 Unauthorized):
```json
{
  "error": "invalid_token", 
  "error_description": "The access token is malformed or invalid"
}
```

**Insufficient Scope Response** (403 Forbidden):
```json
{
  "error": "insufficient_scope",
  "error_description": "The request requires higher privileges than provided by the access token"
}
```

## Provider Implementation Contract

### Configuration Schema

```hcl
provider "hiiretail_iam" {
  # Required authentication parameters
  tenant_id     = string  # Tenant identifier for API calls
  client_id     = string  # OAuth2 client identifier
  client_secret = string  # OAuth2 client secret (sensitive)
  
  # Optional configuration
  base_url      = string  # API base URL (defaults to production)
  token_url     = string  # OAuth2 token endpoint (uses discovery if not set)
  scopes        = list(string)  # OAuth2 scopes (defaults to ["iam:read", "iam:write"])
  
  # Timeout and retry configuration
  timeout       = string  # HTTP timeout (default: "30s")
  max_retries   = number  # Maximum retry attempts (default: 3)
}
```

### Environment Variable Support

```bash
# Authentication (required)
export HIIRETAIL_TENANT_ID="tenant-123"
export HIIRETAIL_CLIENT_ID="client-456"
export HIIRETAIL_CLIENT_SECRET="secret-789"

# Optional configuration
export HIIRETAIL_BASE_URL="https://iam-api.retailsvc-test.com"
export HIIRETAIL_TOKEN_URL="https://auth.retailsvc.com/oauth2/token"
export HIIRETAIL_SCOPES="iam:read,iam:write"
export HIIRETAIL_TIMEOUT="30s"
export HIIRETAIL_MAX_RETRIES="3"
```

### Error Handling Contract

**Configuration Errors**:
- Missing required parameters → Clear error message with parameter name
- Invalid URLs → URL validation error with expected format
- Invalid timeout/retry values → Range validation error

**Authentication Errors**:
- Discovery failure → Fallback to manual configuration or clear error
- Invalid credentials → Authentication error with troubleshooting guidance
- Token acquisition failure → Network vs credential error differentiation
- Token refresh failure → Retry with exponential backoff

**Runtime Errors**:
- API call with expired token → Automatic token refresh and retry
- Network timeouts → Retry with backoff (up to max_retries)
- Rate limiting → Respect retry-after header and backoff

### Token Management Contract

**Token Acquisition**:
1. Attempt OAuth2 discovery if token_url not explicitly configured
2. Use discovered or configured token endpoint for client credentials flow
3. Cache token with expiration time
4. Return authenticated HTTP client

**Token Refresh**:
1. Detect token expiration from API responses (401 with invalid_token)
2. Acquire new token using cached credentials
3. Retry original API call with new token
4. Fail after max_retries exceeded

**Token Security**:
1. Never log client_secret or access_token values
2. Mark client_secret as sensitive in Terraform schema
3. Use secure HTTP client with TLS verification
4. Clear tokens from memory on provider destruction

## Test Scenarios

### Unit Tests
- OAuth2 discovery endpoint parsing
- Token acquisition with valid/invalid credentials
- Error response handling and mapping
- Configuration validation
- Token caching and expiration logic

### Integration Tests
- Real OCMS endpoint discovery
- Token acquisition with test credentials
- API calls with acquired tokens
- Token refresh scenarios
- Error handling with real service responses

### Acceptance Tests
- Provider configuration with various parameter combinations
- Environment variable configuration
- Authentication failure scenarios
- Long-running operations with token refresh
- Concurrent operations with shared token cache

---

**Contract Status**: ✅ Complete  
**Next**: Data Model and Quickstart Guide