# Data Model: OAuth2 Client Credentials Authentication

**Feature**: Correctly Handle Client Credentials  
**Date**: October 1, 2025  
**Status**: Phase 1 Design

## Core Data Structures

### 1. Provider Configuration Model

```go
// HiiRetailIamProviderModel represents the provider configuration schema
type HiiRetailIamProviderModel struct {
    // Required authentication parameters
    TenantId     types.String `tfsdk:"tenant_id"`
    ClientId     types.String `tfsdk:"client_id"`
    ClientSecret types.String `tfsdk:"client_secret"`
    
    // Optional configuration
    BaseUrl      types.String `tfsdk:"base_url"`
    TokenUrl     types.String `tfsdk:"token_url"`     // Optional: discovered if not set
    Scopes       types.List   `tfsdk:"scopes"`        // OAuth2 scopes
    
    // Timeout and retry configuration
    Timeout      types.String `tfsdk:"timeout"`       // HTTP timeout duration
    MaxRetries   types.Int64  `tfsdk:"max_retries"`   // Maximum retry attempts
}
```

### 2. OAuth2 Discovery Model

```go
// OIDCDiscoveryResponse represents the OpenID Connect discovery response
type OIDCDiscoveryResponse struct {
    Issuer                   string   `json:"issuer"`
    TokenEndpoint           string   `json:"token_endpoint"`
    AuthorizationEndpoint   string   `json:"authorization_endpoint"`
    JwksUri                 string   `json:"jwks_uri"`
    GrantTypesSupported     []string `json:"grant_types_supported"`
    TokenEndpointAuthMethods []string `json:"token_endpoint_auth_methods_supported"`
    ResponseTypesSupported  []string `json:"response_types_supported"`
    ScopesSupported         []string `json:"scopes_supported"`
}

// DiscoveryCache caches discovery results to minimize network calls
type DiscoveryCache struct {
    Response   *OIDCDiscoveryResponse
    ExpiresAt  time.Time
    mutex      sync.RWMutex
}
```

### 3. Authentication Client Model

```go
// AuthClient manages OAuth2 authentication and token lifecycle
type AuthClient struct {
    // Configuration
    config       *oauth2.Config
    tokenSource  oauth2.TokenSource
    httpClient   *http.Client
    
    // Discovery cache
    discovery    *DiscoveryCache
    
    // Provider settings
    tenantId     string
    baseUrl      string
    timeout      time.Duration
    maxRetries   int
    
    // Thread safety
    mutex        sync.RWMutex
}

// APIClient represents the enhanced API client with OAuth2 authentication
type APIClient struct {
    // Authentication
    authClient   *AuthClient
    
    // Legacy fields (maintained for compatibility)
    BaseURL      string
    TenantID     string
    HTTPClient   *http.Client
}
```

### 4. Error Models

```go
// AuthError represents authentication-specific errors
type AuthError struct {
    Type        AuthErrorType
    Message     string
    Underlying  error
    Retryable   bool
    RetryAfter  time.Duration // For rate limiting
}

type AuthErrorType int

const (
    AuthErrorUnknown AuthErrorType = iota
    AuthErrorConfiguration          // Invalid provider configuration
    AuthErrorDiscovery             // OAuth2 discovery failure
    AuthErrorCredentials           // Invalid client credentials
    AuthErrorNetwork               // Network/connectivity issues
    AuthErrorServerError           // OAuth2 server errors
    AuthErrorRateLimit             // Rate limiting
    AuthErrorTokenExpired          // Token expiration during operation
)

// ConfigValidationError represents configuration validation errors
type ConfigValidationError struct {
    Field       string
    Value       interface{}
    Constraint  string
    Suggestion  string
}
```

### 5. Token Management Models

```go
// TokenCache manages token caching and refresh logic
type TokenCache struct {
    token       *oauth2.Token
    mutex       sync.RWMutex
    lastRefresh time.Time
}

// RetryConfig defines retry behavior for different error types
type RetryConfig struct {
    MaxAttempts    int
    BaseDelay      time.Duration
    MaxDelay       time.Duration
    Multiplier     float64
    Jitter         bool
    RetryableErrors map[AuthErrorType]bool
}
```

## Data Flows

### 1. Provider Configuration Flow

```
User Configuration → Environment Variables → Validation → Provider Model
                                                      ↓
                  Schema Validation ← Configuration Merge ← Defaults
                                                      ↓
                                              AuthClient Creation
```

### 2. OAuth2 Discovery Flow

```
Provider Configure → Discovery URL Construction → HTTP GET Request
                                                       ↓
Error Handling ← Response Validation ← Discovery Response Parse
       ↓                                           ↓
Fallback Config ← Cache Miss/Expire? → Cache Discovery Results
       ↓                                           ↓
Manual Token URL ← Use Cached → Extract Token Endpoint
```

### 3. Token Acquisition Flow

```
API Call → Token Required? → Token Valid? → Use Cached Token
              ↓                   ↓              ↓
        Get New Token ← No → Refresh Token → API Call with Token
              ↓                               ↓
    OAuth2 Client Credentials Flow → Token Response → Cache Token
              ↓                               ↓
        Error Handling ← Network/Auth Error ← Success
```

### 4. Error Handling Flow

```
Error Occurred → Error Type Classification → Retryable?
                                                ↓
Rate Limited? → Extract Retry-After → Wait → Retry with Backoff
                                                ↓
Network Error? → Exponential Backoff → Retry Count < Max? → Retry
                                                ↓
Auth Error? → Credential Validation → User Error Message → Fail
```

## State Management

### 1. Provider State

```go
// Provider state maintained throughout Terraform execution
type ProviderState struct {
    // Configuration (immutable after Configure)
    config    *HiiRetailIamProviderModel
    
    // Runtime state (mutable, thread-safe)
    authClient *AuthClient
    apiClient  *APIClient
    
    // Metrics and monitoring
    metrics   *AuthMetrics
}

type AuthMetrics struct {
    TokenAcquisitions int64
    TokenRefreshes    int64
    AuthErrors        int64
    NetworkErrors     int64
    LastSuccessAuth   time.Time
    LastErrorAuth     time.Time
}
```

### 2. Token State

```go
// Token state managed by OAuth2 library with custom enhancements
type TokenState struct {
    // OAuth2 token (managed by golang.org/x/oauth2)
    token         *oauth2.Token
    
    // Custom metadata
    acquiredAt    time.Time
    refreshCount  int
    lastUsed      time.Time
    
    // Thread safety
    mutex         sync.RWMutex
}
```

## Validation Rules

### 1. Configuration Validation

```go
// Validation constraints for provider configuration
var ValidationRules = map[string]ValidationRule{
    "tenant_id": {
        Required: true,
        Pattern:  `^[a-zA-Z0-9\-]+$`,
        MinLength: 1,
        MaxLength: 255,
    },
    "client_id": {
        Required: true,
        Pattern:  `^[a-zA-Z0-9\-_]+$`,
        MinLength: 1,
        MaxLength: 255,
    },
    "client_secret": {
        Required: true,
        MinLength: 8,
        MaxLength: 1024,
        Sensitive: true,
    },
    "base_url": {
        Required: false,
        URLFormat: true,
        Schemes: []string{"https"},
    },
    "token_url": {
        Required: false,
        URLFormat: true,
        Schemes: []string{"https"},
    },
    "timeout": {
        Required: false,
        DurationFormat: true,
        MinValue: "1s",
        MaxValue: "300s",
        Default: "30s",
    },
    "max_retries": {
        Required: false,
        IntRange: [0, 10],
        Default: 3,
    },
}
```

### 2. Discovery Response Validation

```go
// Validation for OAuth2 discovery response
func (d *OIDCDiscoveryResponse) Validate() error {
    if d.Issuer == "" {
        return fmt.Errorf("missing required field: issuer")
    }
    
    if d.TokenEndpoint == "" {
        return fmt.Errorf("missing required field: token_endpoint")
    }
    
    if !contains(d.GrantTypesSupported, "client_credentials") {
        return fmt.Errorf("client_credentials grant type not supported")
    }
    
    if !isValidURL(d.TokenEndpoint) {
        return fmt.Errorf("invalid token_endpoint URL")
    }
    
    return nil
}
```

## Security Considerations

### 1. Credential Security

```go
// Secure credential handling
type SecureCredential struct {
    value     []byte // Encrypted in memory
    encrypted bool
    mutex     sync.RWMutex
}

func (s *SecureCredential) Set(plaintext string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Encrypt credential in memory (implementation depends on platform)
    s.value = encrypt([]byte(plaintext))
    s.encrypted = true
}

func (s *SecureCredential) Get() string {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    if !s.encrypted {
        return ""
    }
    
    return string(decrypt(s.value))
}
```

### 2. Token Security

```go
// Token security measures
type SecureToken struct {
    token     *oauth2.Token
    hash      string    // Hash for validation
    createdAt time.Time
    mutex     sync.RWMutex
}

func (t *SecureToken) IsValid() bool {
    t.mutex.RLock()
    defer t.mutex.RUnlock()
    
    if t.token == nil {
        return false
    }
    
    // Validate token hash to detect tampering
    if computeHash(t.token) != t.hash {
        return false
    }
    
    return t.token.Valid()
}
```

## Performance Optimizations

### 1. Connection Pooling

```go
// HTTP client with connection pooling
func createHTTPClient(timeout time.Duration) *http.Client {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   timeout,
    }
}
```

### 2. Discovery Caching

```go
// Discovery result caching
const (
    DiscoveryCacheTTL = 1 * time.Hour
    DiscoveryMaxAge   = 24 * time.Hour
)

func (c *DiscoveryCache) Get(url string) (*OIDCDiscoveryResponse, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    if c.Response == nil || time.Now().After(c.ExpiresAt) {
        return nil, false
    }
    
    return c.Response, true
}
```

## Testing Data

### 1. Test Configuration

```go
// Test configuration for unit tests
var TestConfig = &HiiRetailIamProviderModel{
    TenantId:     types.StringValue("test-tenant-123"),
    ClientId:     types.StringValue("test-client-456"),
    ClientSecret: types.StringValue("test-secret-789"),
    BaseUrl:      types.StringValue("https://iam-api.retailsvc-test.com"),
    TokenUrl:     types.StringValue("https://auth.retailsvc-test.com/oauth2/token"),
    Timeout:      types.StringValue("30s"),
    MaxRetries:   types.Int64Value(3),
}
```

### 2. Mock Responses

```go
// Mock OAuth2 discovery response
var MockDiscoveryResponse = &OIDCDiscoveryResponse{
    Issuer:                 "https://auth.retailsvc-test.com",
    TokenEndpoint:         "https://auth.retailsvc-test.com/oauth2/token",
    AuthorizationEndpoint: "https://auth.retailsvc-test.com/oauth2/authorize",
    JwksUri:               "https://auth.retailsvc-test.com/.well-known/jwks.json",
    GrantTypesSupported:   []string{"client_credentials", "authorization_code"},
    TokenEndpointAuthMethods: []string{"client_secret_basic", "client_secret_post"},
    ResponseTypesSupported:   []string{"code", "token"},
    ScopesSupported:          []string{"iam:read", "iam:write", "iam:admin"},
}

// Mock token response
var MockTokenResponse = &oauth2.Token{
    AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    TokenType:    "Bearer",
    RefreshToken: "",
    Expiry:       time.Now().Add(1 * time.Hour),
}
```

---

**Data Model Status**: ✅ Complete  
**Next**: Quickstart Guide and GitHub Copilot Instructions