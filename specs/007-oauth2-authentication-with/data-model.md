# Data Model: OAuth2 Authentication System
**Feature**: OAuth2 Authentication with Environment-Specific Endpoints  
**Date**: October 1, 2025

## Core Entities

### 1. OAuth2 Configuration
**Purpose**: Holds OAuth2 client credentials and configuration parameters

**Fields**:
- `ClientID` (string, required): OAuth2 client identifier for HiiRetail IAM
- `ClientSecret` (string, required, sensitive): OAuth2 client secret for authentication
- `TenantID` (string, required): HiiRetail tenant identifier for environment detection
- `TokenURL` (string, optional): Override for OAuth2 token endpoint (testing only)
- `APIURL` (string, optional): Override for IAM API endpoint (testing only)
- `MockMode` (bool, optional): Enable mock server mode for testing

**Validation Rules**:
- ClientID must be non-empty string
- ClientSecret must be non-empty string and marked as sensitive
- TenantID must be non-empty string
- TokenURL must be valid HTTPS URL if provided
- APIURL must be valid HTTPS URL if provided
- MockMode can only be true in non-production environments

**Relationships**: Used by AuthClient for authentication operations

### 2. AuthClient
**Purpose**: Manages OAuth2 authentication lifecycle and token operations

**Fields**:
- `config` (OAuth2Configuration): OAuth2 configuration parameters
- `httpClient` (*http.Client): HTTP client for API requests
- `tokenCache` (TokenCache): In-memory token storage
- `mutex` (sync.RWMutex): Thread-safety for concurrent access

**State Transitions**:
```
Initialized → Authenticating → Authenticated → Token Expired → Refresh → Authenticated
                     ↓              ↓               ↓
                   Failed        Active         Refresh Failed
                     ↓              ↓               ↓
                 Error State    Valid State    Re-authenticate
```

**Validation Rules**:
- Must have valid OAuth2Configuration before authentication
- HTTP client must enforce TLS 1.2+
- Token cache must be thread-safe
- Authentication failures must be properly categorized

**Relationships**: Uses OAuth2Configuration, creates TokenCache, integrates with Provider

### 3. TokenCache
**Purpose**: Manages OAuth2 access token storage and lifecycle

**Fields**:
- `AccessToken` (string, sensitive): Current OAuth2 access token
- `TokenType` (string): Token type (typically "Bearer")
- `ExpiresAt` (time.Time): Token expiration timestamp
- `RefreshAt` (time.Time): Calculated refresh time (60s before expiration)
- `mutex` (sync.RWMutex): Thread-safety for token access

**State Transitions**:
```
Empty → Token Stored → Valid → Near Expiry → Refreshing → Updated
          ↓              ↓         ↓             ↓          ↓
      Accessible    Accessible  Refresh     Refreshing   Accessible
                                Required
```

**Validation Rules**:
- AccessToken must be non-empty when cached
- ExpiresAt must be in the future for valid tokens
- RefreshAt must be before ExpiresAt
- All token operations must be thread-safe

**Relationships**: Owned by AuthClient, used by HTTP requests

### 4. EndpointResolver
**Purpose**: Determines correct API endpoints based on tenant configuration

**Fields**:
- `TenantID` (string): Tenant identifier for environment detection
- `AuthURL` (string): Resolved OAuth2 authentication endpoint
- `APIURL` (string): Resolved IAM API endpoint
- `IsTestEnvironment` (bool): Detected environment type

**Resolution Logic**:
```
Input: TenantID
↓
Pattern Analysis:
- Contains "test|dev|staging" (case-insensitive) → Test Environment
- Environment override HIIRETAIL_FORCE_TEST_ENV=true → Test Environment
- Default → Live Environment
↓
Endpoint Selection:
- Test: auth.retailsvc.com + iam-api.retailsvc-test.com
- Live: auth.retailsvc.com + iam-api.retailsvc.com
- Mock: Use override URLs from configuration
```

**Validation Rules**:
- TenantID must be non-empty
- Resolved URLs must be valid HTTPS endpoints
- Environment detection must be deterministic
- Mock overrides only allowed in test mode

**Relationships**: Used by AuthClient for endpoint determination

### 5. AuthError
**Purpose**: Represents authentication-specific error conditions

**Fields**:
- `Type` (AuthErrorType): Classification of authentication error
- `Message` (string): Human-readable error description
- `Code` (string): Machine-readable error code
- `Retryable` (bool): Whether the operation can be retried
- `Cause` (error): Underlying error cause

**Error Types**:
```go
type AuthErrorType int

const (
    ErrInvalidCredentials AuthErrorType = iota
    ErrTokenExpired
    ErrNetworkError
    ErrConfigurationError
    ErrTenantNotFound
    ErrEndpointUnreachable
)
```

**Classification Rules**:
- HTTP 401/403 → ErrInvalidCredentials (non-retryable)
- Token expiration → ErrTokenExpired (retryable with refresh)
- HTTP 5xx/timeout → ErrNetworkError (retryable with backoff)
- Invalid configuration → ErrConfigurationError (non-retryable)
- Tenant resolution failure → ErrTenantNotFound (non-retryable)
- Connection failures → ErrEndpointUnreachable (retryable)

**Relationships**: Used throughout authentication system for error handling

## Entity Relationships

```
OAuth2Configuration
    ↓ (used by)
AuthClient
    ↓ (creates/manages)
TokenCache
    ↓ (provides tokens to)
HTTP Requests

AuthClient
    ↓ (uses)
EndpointResolver
    ↓ (determines)
API Endpoints

AuthClient
    ↓ (generates)
AuthError
    ↓ (consumed by)
Provider Error Handling
```

## Data Flow

### Authentication Flow
```
1. Provider initialization
   → OAuth2Configuration created from Terraform config
   
2. AuthClient creation
   → EndpointResolver determines URLs from TenantID
   → HTTP client configured with TLS settings
   
3. Token acquisition
   → OAuth2 client credentials flow to auth.retailsvc.com
   → Token stored in TokenCache with expiration
   
4. API requests
   → TokenCache provides valid token
   → Requests sent to environment-specific IAM API
   
5. Token refresh (automatic)
   → Near expiration detected
   → New token acquired and cached
   → Transparent to API operations
```

### Error Handling Flow
```
1. Operation failure
   → AuthError created with classification
   
2. Error analysis
   → Retryable errors trigger exponential backoff
   → Non-retryable errors propagated immediately
   
3. Recovery actions
   → Token refresh for expiration errors
   → Re-authentication for credential errors
   → Endpoint resolution for network errors
```

## Validation Summary

All entities include comprehensive validation rules ensuring:
- Security: Sensitive fields marked and protected
- Reliability: Thread-safe operations and proper error handling
- Performance: Efficient caching and minimal network operations
- Maintainability: Clear separation of concerns and error classification