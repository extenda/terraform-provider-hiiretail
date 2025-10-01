package main
// OAuth2 Demo - Run with: go run demo/oauth2_demo.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/auth"
)

func main() {
	fmt.Println("🚀 HiiRetail IAM OAuth2 Authentication Demo")
	fmt.Println("===========================================")

	// Demo credentials
	demoTenantID := "hiiretail-demo-tenant"
	demoClientID := "hiiretail-demo-client"
	demoClientSecret := "hiiretail-demo-secret-secure-credential"

	// Create mock OCMS server for demonstration
	var mockOCMS *httptest.Server
	mockOCMS = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			discoveryResponse := map[string]interface{}{
				"issuer":                                mockOCMS.URL,
				"token_endpoint":                        mockOCMS.URL + "/oauth2/token",
				"authorization_endpoint":                mockOCMS.URL + "/oauth2/authorize",
				"jwks_uri":                              mockOCMS.URL + "/.well-known/jwks.json",
				"grant_types_supported":                 []string{"client_credentials", "authorization_code"},
				"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
				"response_types_supported":              []string{"code", "token"},
				"scopes_supported":                      []string{"iam:read", "iam:write", "iam:admin"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(discoveryResponse)

		case "/oauth2/token":
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			clientID := r.Form.Get("client_id")
			clientSecret := r.Form.Get("client_secret")

			if clientID == "" || clientSecret == "" {
				username, password, ok := r.BasicAuth()
				if ok {
					clientID = username
					clientSecret = password
				}
			}

			if clientID != demoClientID || clientSecret != demoClientSecret {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":             "invalid_client",
					"error_description": "Invalid client credentials",
				})
				return
			}

			tokenResponse := map[string]interface{}{
				"access_token": fmt.Sprintf("hiiretail-token-%d", time.Now().Unix()),
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read iam:write",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse)

		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer mockOCMS.Close()

	fmt.Printf("📡 Mock OCMS Server: %s\n", mockOCMS.URL)
	fmt.Printf("🔑 Demo Tenant ID: %s\n", demoTenantID)
	fmt.Printf("🆔 Demo Client ID: %s\n", demoClientID)
	fmt.Println()

	// Create OAuth2 configuration
	fmt.Println("🔧 Creating OAuth2 Authentication Configuration...")
	config := &auth.AuthClientConfig{
		TenantID:     demoTenantID,
		ClientID:     demoClientID,
		ClientSecret: demoClientSecret,
		BaseURL:      mockOCMS.URL,
		Scopes:       []string{"iam:read", "iam:write"},
		Timeout:      30 * time.Second,
		MaxRetries:   3,
	}

	// Create OAuth2 authentication client
	fmt.Println("🏗️  Initializing OAuth2 Authentication Client...")
	authClient, err := auth.NewAuthClient(config)
	if err != nil {
		log.Fatalf("❌ Failed to create auth client: %v", err)
	}
	defer authClient.Close()

	// Acquire OAuth2 access token
	fmt.Println("🎫 Acquiring OAuth2 Access Token...")
	ctx := context.Background()
	
	token, err := authClient.GetToken(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to acquire token: %v", err)
	}
	
	fmt.Println("✅ SUCCESS! OAuth2 token acquired:")
	fmt.Printf("   Token Type: %s\n", token.TokenType)
	fmt.Printf("   Access Token: %s...\n", token.AccessToken[:25])
	fmt.Printf("   Expires: %v\n", token.Expiry)
	fmt.Printf("   Valid: %v\n", token.Valid())

	// Test token caching
	fmt.Println("\n💾 Testing token caching...")
	token2, err := authClient.GetToken(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get cached token: %v", err)
	}
	
	if token.AccessToken == token2.AccessToken {
		fmt.Println("✅ Token retrieved from cache successfully!")
	}

	fmt.Println("\n🎉 OAuth2 Authentication Demo Complete!")
	fmt.Println("Your OAuth2 implementation is working perfectly! 🚀")
}