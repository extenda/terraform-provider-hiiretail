# Provider Configuration Contract

## Provider Schema Contract

```hcl
provider "hiiretail" {
  # Required OAuth2 authentication
  client_id     = string # sensitive
  client_secret = string # sensitive
  
  # Optional OAuth2 configuration  
  auth_url = string # default: "https://auth.retailsvc.com/oauth2/token"
  audience = string # default: "https://iam-api.retailsvc.com"
  
  # Optional connection settings
  timeout_seconds = number # default: 30, range: 5-300
  max_retries     = number # default: 3, range: 0-10
  
  # Optional service endpoint overrides
  iam_endpoint = string # default: "https://iam-api.retailsvc.com"
  ccc_endpoint = string # default: "https://ccc-api.retailsvc.com" (future)
}
```

## Provider Validation Contract

**Authentication Validation**:
- `client_id` and `client_secret` must be non-empty strings
- OAuth2 token endpoint must respond with valid access token
- Access token must be valid for specified audience

**URL Validation**:
- All endpoint URLs must use HTTPS protocol
- URLs must be well-formed and reachable
- Service endpoints must respond to health checks

**Configuration Validation**:
- `timeout_seconds` must be integer between 5 and 300
- `max_retries` must be integer between 0 and 10
- Invalid configurations must fail with clear error messages

## Provider Behavior Contract

**Authentication Flow**:
1. Validate configuration parameters
2. Perform OAuth2 client credentials flow
3. Store access token securely (never logged)  
4. Refresh token automatically when expired
5. Handle authentication failures with retry logic

**Service Registration**:
1. Register all available service modules (IAM, CCC, etc.)
2. Initialize service-specific clients with appropriate endpoints
3. Validate service availability during provider configuration
4. Fail fast with clear error messages for unavailable services

**Error Handling**:
- Network errors: Retry with exponential backoff
- Authentication errors: Clear error message with troubleshooting steps
- Service unavailable: Graceful degradation with informative errors
- Invalid configuration: Immediate failure with specific validation errors