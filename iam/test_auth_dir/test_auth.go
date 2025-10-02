package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// Get credentials from environment or use the ones from terraform.tfvars
	clientID := os.Getenv("HIIRETAIL_CLIENT_ID")
	clientSecret := os.Getenv("HIIRETAIL_CLIENT_SECRET")
	tenantID := os.Getenv("HIIRETAIL_TENANT_ID")

	if clientID == "" {
		clientID = "b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
	}
	if clientSecret == "" {
		clientSecret = "726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"
	}
	if tenantID == "" {
		tenantID = "CIR7nQwtS0rA6t0S6ejd"
	}

	tokenURL := "https://auth.retailsvc.com/oauth2/token"

	fmt.Printf("🔐 Testing OAuth2 Token Acquisition\n")
	fmt.Printf("Token URL: %s\n", tokenURL)
	fmt.Printf("Client ID: %s\n", clientID[:20]+"...")
	fmt.Printf("Tenant ID: %s\n", tenantID)
	fmt.Println()

	// Test 1: Manual Basic Auth approach
	fmt.Printf("📡 Test 1: Manual Basic Auth Request\n")
	err := testManualBasicAuth(clientID, clientSecret, tokenURL)
	if err != nil {
		fmt.Printf("❌ Manual Basic Auth failed: %v\n", err)
	} else {
		fmt.Printf("✅ Manual Basic Auth succeeded!\n")
	}
	fmt.Println()

	// Test 2: OAuth2 client credentials flow (what the provider uses)
	fmt.Printf("📡 Test 2: OAuth2 Client Credentials Flow\n")
	err = testOAuth2Flow(clientID, clientSecret, tokenURL)
	if err != nil {
		fmt.Printf("❌ OAuth2 flow failed: %v\n", err)
	} else {
		fmt.Printf("✅ OAuth2 flow succeeded!\n")
	}
}

func testManualBasicAuth(clientID, clientSecret, tokenURL string) error {
	// Create Basic Auth header manually
	auth := clientID + ":" + clientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "iam:read iam:write")

	// Create request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+basicAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	if resp.StatusCode == 200 {
		fmt.Printf("✅ Successfully received token!\n")
		return nil
	}

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func testOAuth2Flow(clientID, clientSecret, tokenURL string) error {
	// Create OAuth2 config (same as provider uses)
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{"iam:read", "iam:write"},
		AuthStyle:    oauth2.AuthStyleInHeader, // This uses Basic auth
	}

	// Get token
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := config.Token(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	fmt.Printf("✅ Token acquired successfully!\n")
	fmt.Printf("Token Type: %s\n", token.TokenType)
	fmt.Printf("Access Token: %s...\n", token.AccessToken[:20])
	fmt.Printf("Expires: %v\n", token.Expiry)
	fmt.Printf("Valid: %v\n", token.Valid())

	return nil
}
