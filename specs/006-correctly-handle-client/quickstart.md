# Quickstart Guide: OAuth2 Client Credentials Authentication

**Feature**: Correctly Handle Client Credentials  
**Target**: Development Team  
**Duration**: 30-45 minutes implementation

## Development Setup

### Prerequisites

1. **Go 1.21+** installed
2. **Terraform 1.0+** for testing
3. **Git** for version control
4. **Access to OCMS**: Test client credentials for Hii Retail OAuth Client Management Service

### Environment Setup

```bash
# Clone the repository (if not already done)
cd /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam

# Switch to feature branch
git checkout 006-correctly-handle-client

# Install dependencies
go mod tidy

# Set up test environment variables
export HIIRETAIL_TENANT_ID="test-tenant-123"
export HIIRETAIL_CLIENT_ID="test-client-456"  
export HIIRETAIL_CLIENT_SECRET="test-secret-789"
export HIIRETAIL_BASE_URL="https://iam-api.retailsvc-test.com"
```

## Implementation Steps

### Step 1: Create OAuth2 Discovery Client (15 minutes)

**File**: `internal/provider/auth/discovery.go`

```go
package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// OIDCDiscoveryResponse represents OAuth2 discovery response
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

// DiscoveryClient handles OAuth2 endpoint discovery
type DiscoveryClient struct {
    httpClient *http.Client
    cache      sync.Map // URL -> cached response
}

// NewDiscoveryClient creates a new discovery client
func NewDiscoveryClient(timeout time.Duration) *DiscoveryClient {
    return &DiscoveryClient{
        httpClient: &http.Client{Timeout: timeout},
    }
}

// Discover fetches OAuth2 configuration from discovery endpoint
func (d *DiscoveryClient) Discover(ctx context.Context, baseURL string) (*OIDCDiscoveryResponse, error) {
    discoveryURL := fmt.Sprintf("%s/.well-known/openid-configuration", baseURL)
    
    // Check cache first
    if cached, ok := d.cache.Load(discoveryURL); ok {
        if entry, ok := cached.(*cacheEntry); ok && time.Now().Before(entry.expiresAt) {
            return entry.response, nil
        }
    }
    
    // Fetch from endpoint
    req, err := http.NewRequestWithContext(ctx, "GET", discoveryURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create discovery request: %w", err)
    }
    
    resp, err := d.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch discovery configuration: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("discovery request failed with status %d", resp.StatusCode)
    }
    
    var config OIDCDiscoveryResponse
    if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
        return nil, fmt.Errorf("failed to decode discovery response: %w", err)
    }
    
    // Validate required fields
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid discovery response: %w", err)
    }
    
    // Cache result
    d.cache.Store(discoveryURL, &cacheEntry{
        response:  &config,
        expiresAt: time.Now().Add(1 * time.Hour),
    })
    
    return &config, nil
}

type cacheEntry struct {
    response  *OIDCDiscoveryResponse
    expiresAt time.Time
}

// Validate checks if the discovery response contains required fields
func (r *OIDCDiscoveryResponse) Validate() error {
    if r.Issuer == "" {
        return fmt.Errorf("missing issuer")
    }
    if r.TokenEndpoint == "" {
        return fmt.Errorf("missing token_endpoint")
    }
    
    // Check if client_credentials is supported
    for _, grant := range r.GrantTypesSupported {
        if grant == "client_credentials" {
            return nil
        }
    }
    
    return fmt.Errorf("client_credentials grant type not supported")
}
```

### Step 2: Enhance Authentication Client (15 minutes)

**File**: `internal/provider/auth/client.go`

```go
package auth

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
    
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/clientcredentials"
)

// AuthClient manages OAuth2 authentication
type AuthClient struct {
    config      *clientcredentials.Config
    tokenSource oauth2.TokenSource
    httpClient  *http.Client
    discovery   *DiscoveryClient
    
    tenantID    string
    baseURL     string
    tokenURL    string
    timeout     time.Duration
    maxRetries  int
    
    mutex       sync.RWMutex
}

// NewAuthClient creates a new authentication client
func NewAuthClient(config AuthConfig) (*AuthClient, error) {
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid auth configuration: %w", err)
    }
    
    client := &AuthClient{
        tenantID:   config.TenantID,
        baseURL:    config.BaseURL,
        tokenURL:   config.TokenURL,
        timeout:    config.Timeout,
        maxRetries: config.MaxRetries,
        discovery:  NewDiscoveryClient(config.Timeout),
    }
    
    if err := client.initialize(config); err != nil {
        return nil, fmt.Errorf("failed to initialize auth client: %w", err)
    }
    
    return client, nil
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
    TenantID     string
    ClientID     string
    ClientSecret string
    BaseURL      string
    TokenURL     string        // Optional: discovered if not set
    Scopes       []string      // OAuth2 scopes
    Timeout      time.Duration
    MaxRetries   int
}

// Validate validates the authentication configuration
func (c *AuthConfig) Validate() error {
    if c.TenantID == "" {
        return fmt.Errorf("tenant_id is required")
    }
    if c.ClientID == "" {
        return fmt.Errorf("client_id is required")
    }
    if c.ClientSecret == "" {
        return fmt.Errorf("client_secret is required")
    }
    if c.BaseURL == "" {
        return fmt.Errorf("base_url is required")
    }
    if c.Timeout <= 0 {
        c.Timeout = 30 * time.Second
    }
    if c.MaxRetries < 0 {
        c.MaxRetries = 3
    }
    if len(c.Scopes) == 0 {
        c.Scopes = []string{"iam:read", "iam:write"}
    }
    return nil
}

// initialize sets up the OAuth2 configuration
func (c *AuthClient) initialize(config AuthConfig) error {
    ctx := context.Background()
    
    // Determine token endpoint
    tokenEndpoint := config.TokenURL
    if tokenEndpoint == "" {
        // Use discovery to find token endpoint
        discovery, err := c.discovery.Discover(ctx, "https://auth.retailsvc.com")
        if err != nil {
            return fmt.Errorf("failed to discover OAuth2 endpoints: %w", err)
        }
        tokenEndpoint = discovery.TokenEndpoint
    }
    
    // Configure OAuth2 client credentials
    c.config = &clientcredentials.Config{
        ClientID:     config.ClientID,
        ClientSecret: config.ClientSecret,
        TokenURL:     tokenEndpoint,
        Scopes:       config.Scopes,
    }
    
    // Create token source with context that won't be canceled
    ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
        Timeout: c.timeout,
    })
    c.tokenSource = c.config.TokenSource(ctx)
    
    // Create authenticated HTTP client
    c.httpClient = c.config.Client(ctx)
    
    return nil
}

// GetHTTPClient returns an authenticated HTTP client
func (c *AuthClient) GetHTTPClient() *http.Client {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.httpClient
}

// ValidateCredentials tests the authentication configuration
func (c *AuthClient) ValidateCredentials(ctx context.Context) error {
    _, err := c.tokenSource.Token()
    if err != nil {
        return fmt.Errorf("credential validation failed: %w", err)
    }
    return nil
}
```

### Step 3: Update Provider Configuration (10 minutes)

**File**: `internal/provider/provider.go` (enhance existing Configure method)

```go
// Add to imports
import (
    "github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/auth"
)

// Enhance the Configure method
func (p *HiiRetailIamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var data HiiRetailIamProviderModel
    
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Extract configuration values (existing code...)
    
    // Create enhanced authentication configuration
    authConfig := auth.AuthConfig{
        TenantID:     tenantId,
        ClientID:     clientId,
        ClientSecret: clientSecret,
        BaseURL:      baseUrl,
        TokenURL:     "", // Will use discovery
        Scopes:       []string{"iam:read", "iam:write"},
        Timeout:      30 * time.Second,
        MaxRetries:   3,
    }
    
    // Create authentication client
    authClient, err := auth.NewAuthClient(authConfig)
    if err != nil {
        resp.Diagnostics.AddError(
            "Authentication Configuration Error",
            fmt.Sprintf("Failed to configure OAuth2 authentication: %s", err.Error()),
        )
        return
    }
    
    // Validate credentials
    if err := authClient.ValidateCredentials(ctx); err != nil {
        resp.Diagnostics.AddError(
            "Authentication Validation Error", 
            fmt.Sprintf("Failed to validate OAuth2 credentials: %s", err.Error()),
        )
        return
    }
    
    // Create API client with enhanced authentication
    apiClient := &APIClient{
        BaseURL:    baseUrl,
        TenantID:   tenantId,
        HTTPClient: authClient.GetHTTPClient(),
        AuthClient: authClient, // Add auth client reference
    }
    
    resp.DataSourceData = apiClient
    resp.ResourceData = apiClient
}
```

### Step 4: Add Tests (15 minutes)

**File**: `internal/provider/auth/discovery_test.go`

```go
package auth

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestDiscoveryClient_Discover(t *testing.T) {
    tests := []struct {
        name           string
        serverResponse string
        serverStatus   int
        expectError    bool
    }{
        {
            name: "successful discovery",
            serverResponse: `{
                "issuer": "https://auth.retailsvc.com",
                "token_endpoint": "https://auth.retailsvc.com/oauth2/token",
                "grant_types_supported": ["client_credentials"]
            }`,
            serverStatus: http.StatusOK,
            expectError:  false,
        },
        {
            name:         "server error",
            serverStatus: http.StatusInternalServerError,
            expectError:  true,
        },
        {
            name: "invalid JSON",
            serverResponse: `invalid json`,
            serverStatus:   http.StatusOK,
            expectError:    true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.serverStatus)
                if tt.serverResponse != "" {
                    w.Write([]byte(tt.serverResponse))
                }
            }))
            defer server.Close()
            
            client := NewDiscoveryClient(5 * time.Second)
            _, err := client.Discover(context.Background(), server.URL)
            
            if tt.expectError && err == nil {
                t.Error("expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

## Testing Commands

### Unit Tests
```bash
# Test OAuth2 discovery
go test ./internal/provider/auth -v

# Test provider configuration
go test ./internal/provider -v -run TestHiiRetailIamProvider_Configure
```

### Integration Tests
```bash
# Test with real OCMS endpoints (requires valid credentials)
export HIIRETAIL_TENANT_ID="your-tenant"
export HIIRETAIL_CLIENT_ID="your-client" 
export HIIRETAIL_CLIENT_SECRET="your-secret"

go test ./internal/provider -v -run TestProviderIntegration
```

### Manual Testing
```bash  
# Build provider
make build

# Test configuration
terraform init
terraform plan

# Check authentication logs
TF_LOG=DEBUG terraform plan 2>&1 | grep -E "(auth|token|oauth)"
```

## Validation Checklist

### ✅ Functional Requirements
- [ ] OAuth2 discovery client implemented
- [ ] Client credentials flow working
- [ ] Token caching and refresh automatic
- [ ] Error handling with clear messages
- [ ] Configuration validation
- [ ] Environment variable support

### ✅ Security Requirements  
- [ ] Client secrets marked as sensitive
- [ ] No credentials in logs or debug output
- [ ] TLS-only communication enforced
- [ ] Secure token storage

### ✅ Testing Requirements
- [ ] Unit tests for discovery client
- [ ] Unit tests for auth client
- [ ] Integration tests with real endpoints
- [ ] Provider configuration tests
- [ ] Error scenario tests

### ✅ Performance Requirements
- [ ] Discovery response caching
- [ ] Connection pooling
- [ ] Token reuse (no unnecessary acquisitions)
- [ ] Timeout configuration working

## Troubleshooting

### Common Issues

1. **Discovery Fails**
   ```bash
   # Check network connectivity
   curl https://auth.retailsvc.com/.well-known/openid-configuration
   ```

2. **Invalid Credentials**
   ```bash
   # Verify client credentials
   echo "Check HIIRETAIL_CLIENT_ID and HIIRETAIL_CLIENT_SECRET"
   ```

3. **Token Acquisition Fails**
   ```bash
   # Enable debug logging
   export TF_LOG=DEBUG
   terraform plan
   ```

4. **Network Timeouts**
   ```bash
   # Increase timeout in provider configuration
   provider "hiiretail_iam" {
     timeout = "60s"
   }
   ```

## Next Steps

After completing this quickstart:

1. **Review Tests**: Ensure all tests pass
2. **Security Audit**: Review credential handling
3. **Performance Test**: Test with concurrent operations
4. **Documentation**: Update provider documentation
5. **Integration**: Test with all existing resources

## Resources

- [OAuth2 RFC 6749](https://tools.ietf.org/html/rfc6749)
- [OAuth2 Discovery RFC 8414](https://tools.ietf.org/html/rfc8414)
- [golang.org/x/oauth2 Package](https://pkg.go.dev/golang.org/x/oauth2)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin)
- [OCMS Documentation](https://developer.hiiretail.com/docs/ocms/public/concepts/oauth2-authentication/)

---

**Quickstart Status**: ✅ Complete  
**Implementation Time**: 30-45 minutes  
**Next**: Task Generation Phase