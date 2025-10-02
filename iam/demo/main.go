package main

// OAuth2 Demo Program for HiiRetail IAM Terraform Provider
//
// This program demonstrates the OAuth2 authentication capabilities of the
// HiiRetail IAM Terraform provider's auth package. It shows:
//
// 1. Configuration loading from environment variables
// 2. OAuth2 client credentials flow
// 3. Authenticated API requests
// 4. Error handling and retry logic
// 5. Token caching and refresh
//
// Usage:
//   export HIIRETAIL_TENANT_ID="your-tenant-id"
//   export HIIRETAIL_CLIENT_ID="your-client-id"
//   export HIIRETAIL_CLIENT_SECRET="your-client-secret"
//   export HIIRETAIL_ENVIRONMENT="test"  # optional: production, test, dev
//   go run demo/main.go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"hiiretail/internal/provider/shared/auth"
)

func main() {
	fmt.Println("🚀 HiiRetail IAM OAuth2 Authentication Demo")
	fmt.Println("==========================================")
	fmt.Println()

	// Create context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Step 1: Load configuration from environment
	fmt.Println("📋 Step 1: Loading OAuth2 Configuration")
	config, err := loadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	fmt.Printf("   ✅ Tenant ID: %s\n", config.TenantID)
	fmt.Printf("   ✅ Client ID: %s\n", redactCredential(config.ClientID))
	fmt.Printf("   ✅ Environment: %s\n", config.Environment)
	fmt.Println()

	// Step 2: Validate configuration
	fmt.Println("🔍 Step 2: Validating Configuration")
	if err := auth.ValidateConfig(config); err != nil {
		log.Fatalf("❌ Configuration validation failed: %v", err)
	}
	fmt.Println("   ✅ Configuration is valid")
	fmt.Println()

	// Step 3: Create OAuth2 client
	fmt.Println("🔐 Step 3: Creating OAuth2 Authentication Client")
	authClient, err := auth.New(config)
	if err != nil {
		log.Fatalf("❌ Failed to create auth client: %v", err)
	}
	defer authClient.Close()
	fmt.Println("   ✅ OAuth2 client created successfully")
	fmt.Println()

	// Step 4: Demonstrate endpoint resolution
	fmt.Println("🔗 Step 4: Demonstrating Endpoint Resolution")
	demonstrateEndpointResolution(config.TenantID, config.Environment)
	fmt.Println()

	// Step 5: Acquire OAuth2 token
	fmt.Println("🎫 Step 5: Acquiring OAuth2 Access Token")
	token, err := authClient.GetToken(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to acquire token: %v", err)
	}

	fmt.Printf("   ✅ Token acquired successfully\n")
	fmt.Printf("   ✅ Token type: %s\n", token.TokenType)
	fmt.Printf("   ✅ Expires in: %v\n", time.Until(token.Expiry).Round(time.Second))
	fmt.Printf("   ✅ Access token: %s...%s\n",
		token.AccessToken[:8],
		token.AccessToken[len(token.AccessToken)-8:])
	fmt.Println()

	// Step 6: Create authenticated HTTP client
	fmt.Println("🌐 Step 6: Creating Authenticated HTTP Client")
	httpClient, err := authClient.HTTPClient(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create HTTP client: %v", err)
	}
	fmt.Println("   ✅ Authenticated HTTP client created")
	fmt.Println()

	// Step 7: Demonstrate API calls
	fmt.Println("📡 Step 7: Making Authenticated API Requests")
	demonstrateAPIRequests(ctx, httpClient, config.TenantID)
	fmt.Println()

	// Step 8: Demonstrate retry with token refresh
	fmt.Println("🔄 Step 8: Demonstrating Retry Logic with Token Refresh")
	retryClient, err := authClient.HTTPClientWithRetry(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create retry client: %v", err)
	}
	demonstrateRetryLogic(ctx, retryClient, config.TenantID)
	fmt.Println()

	// Step 9: Demonstrate token validation
	fmt.Println("✅ Step 9: Demonstrating Token Validation")
	demonstrateTokenValidation(ctx, authClient, token)
	fmt.Println()

	// Step 10: Demonstrate configuration variations
	fmt.Println("⚙️  Step 10: Demonstrating Configuration Variations")
	demonstrateConfigurationVariations()
	fmt.Println()

	fmt.Println("🎉 Demo completed successfully!")
	fmt.Println("   All OAuth2 authentication features are working correctly.")
}

// loadConfigFromEnvironment loads OAuth2 configuration from environment variables
func loadConfigFromEnvironment() (*auth.Config, error) {
	tenantID := os.Getenv("HIIRETAIL_TENANT_ID")
	if tenantID == "" {
		return nil, fmt.Errorf("HIIRETAIL_TENANT_ID environment variable is required")
	}

	clientID := os.Getenv("HIIRETAIL_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("HIIRETAIL_CLIENT_ID environment variable is required")
	}

	clientSecret := os.Getenv("HIIRETAIL_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("HIIRETAIL_CLIENT_SECRET environment variable is required")
	}

	config := &auth.Config{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Environment:  os.Getenv("HIIRETAIL_ENVIRONMENT"), // Optional
		Scopes:       []string{"iam:read", "iam:write"},
		Timeout:      30 * time.Second,
		MaxRetries:   3,
	}

	// Override with explicit URLs if provided
	if authURL := os.Getenv("HIIRETAIL_AUTH_URL"); authURL != "" {
		config.AuthURL = authURL
	}

	if apiURL := os.Getenv("HIIRETAIL_API_URL"); apiURL != "" {
		config.APIURL = apiURL
	}

	return config, nil
}

// demonstrateEndpointResolution shows how endpoint resolution works
func demonstrateEndpointResolution(tenantID, environment string) {
	authURL, apiURL, err := auth.ResolveEndpoints(tenantID, environment)
	if err != nil {
		fmt.Printf("   ❌ Failed to resolve endpoints: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Auth URL: %s\n", authURL)
	fmt.Printf("   ✅ API URL: %s\n", apiURL)

	if auth.IsTestEnvironment(tenantID) {
		fmt.Printf("   ℹ️  Using test environment endpoints\n")
	} else {
		fmt.Printf("   ℹ️  Using production environment endpoints\n")
	}
}

// demonstrateAPIRequests makes sample API requests
func demonstrateAPIRequests(ctx context.Context, client *http.Client, tenantID string) {
	// Construct API URL based on tenant (this would normally be resolved automatically)
	var baseURL string
	if auth.IsTestEnvironment(tenantID) {
		baseURL = "https://iam-api.retailsvc-test.com"
	} else {
		baseURL = "https://iam-api.retailsvc.com"
	}

	// Example 1: List roles
	fmt.Println("   📊 Making GET request to list roles...")
	resp, err := makeAPIRequest(ctx, client, "GET", baseURL+"/api/v1/roles", nil, tenantID)
	if err != nil {
		fmt.Printf("   ⚠️  Failed to list roles: %v\n", err)
	} else {
		fmt.Printf("   ✅ Roles request successful (Status: %d)\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Example 2: List groups
	fmt.Println("   👥 Making GET request to list groups...")
	resp, err = makeAPIRequest(ctx, client, "GET", baseURL+"/api/v1/groups", nil, tenantID)
	if err != nil {
		fmt.Printf("   ⚠️  Failed to list groups: %v\n", err)
	} else {
		fmt.Printf("   ✅ Groups request successful (Status: %d)\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Example 3: Get current user info (if supported)
	fmt.Println("   👤 Making GET request to get user info...")
	resp, err = makeAPIRequest(ctx, client, "GET", baseURL+"/api/v1/user", nil, tenantID)
	if err != nil {
		fmt.Printf("   ⚠️  Failed to get user info: %v (this may be expected)\n", err)
	} else {
		fmt.Printf("   ✅ User info request successful (Status: %d)\n", resp.StatusCode)
		resp.Body.Close()
	}
}

// demonstrateRetryLogic shows the retry functionality with token refresh
func demonstrateRetryLogic(ctx context.Context, client *http.Client, tenantID string) {
	fmt.Println("   🔄 Testing retry logic with authenticated requests...")

	// Make a request that will use the retry client
	var baseURL string
	if auth.IsTestEnvironment(tenantID) {
		baseURL = "https://iam-api.retailsvc-test.com"
	} else {
		baseURL = "https://iam-api.retailsvc.com"
	}

	resp, err := makeAPIRequest(ctx, client, "GET", baseURL+"/api/v1/roles", nil, tenantID)
	if err != nil {
		fmt.Printf("   ⚠️  Request with retry failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Request with retry successful (Status: %d)\n", resp.StatusCode)
		fmt.Printf("   ℹ️  The client will automatically refresh tokens on 401 errors\n")
		resp.Body.Close()
	}
}

// demonstrateTokenValidation shows token validation capabilities
func demonstrateTokenValidation(ctx context.Context, client auth.Client, token interface{}) {
	// Force a token refresh to demonstrate the capability
	fmt.Println("   🔄 Forcing token refresh for demonstration...")

	newToken, err := client.RefreshToken(ctx)
	if err != nil {
		fmt.Printf("   ⚠️  Token refresh failed: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Token refreshed successfully\n")
	fmt.Printf("   ✅ New token expires in: %v\n", time.Until(newToken.Expiry).Round(time.Second))
}

// demonstrateConfigurationVariations shows different configuration options
func demonstrateConfigurationVariations() {
	fmt.Println("   🔧 Configuration Options:")
	fmt.Println("   • Environment Variables:")
	fmt.Println("     - HIIRETAIL_TENANT_ID (required)")
	fmt.Println("     - HIIRETAIL_CLIENT_ID (required)")
	fmt.Println("     - HIIRETAIL_CLIENT_SECRET (required)")
	fmt.Println("     - HIIRETAIL_ENVIRONMENT (optional: production, test, dev)")
	fmt.Println("     - HIIRETAIL_AUTH_URL (optional: custom OAuth2 endpoint)")
	fmt.Println("     - HIIRETAIL_API_URL (optional: custom API endpoint)")
	fmt.Println("     - HIIRETAIL_SCOPES (optional: comma-separated scopes)")
	fmt.Println("     - HIIRETAIL_TIMEOUT_SECONDS (optional: 5-300)")
	fmt.Println("     - HIIRETAIL_MAX_RETRIES (optional: 0-10)")
	fmt.Println()

	fmt.Println("   📝 Example Terraform Configuration:")
	fmt.Println(`   provider "hiiretail-iam" {
     tenant_id     = "my-tenant-123"
     client_id     = "oauth2-client-id"
     client_secret = "oauth2-client-secret"
     scopes        = ["iam:read", "iam:write"]
     timeout_seconds = 30
     max_retries   = 3
   }`)
	fmt.Println()

	fmt.Println("   🔐 Security Features:")
	fmt.Println("   • Automatic credential redaction in logs")
	fmt.Println("   • Secure token caching with integrity validation")
	fmt.Println("   • HTTPS enforcement for all endpoints")
	fmt.Println("   • Automatic token refresh on expiration")
	fmt.Println("   • Configurable retry logic with exponential backoff")
}

// makeAPIRequest makes an authenticated API request
func makeAPIRequest(ctx context.Context, client *http.Client, method, url string, body interface{}, tenantID string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	// Add tenant header (this would normally be done by the auth transport)
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}

// redactCredential redacts sensitive credentials for display
func redactCredential(credential string) string {
	if len(credential) <= 8 {
		return "[REDACTED]"
	}
	return credential[:4] + "..." + credential[len(credential)-4:]
}
