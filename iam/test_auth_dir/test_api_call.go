package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

	fmt.Printf("🔐 Testing IAM API with OAuth2 Token\n")
	fmt.Printf("Tenant ID: %s\n\n", tenantID)

	// Get OAuth2 token
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://auth.retailsvc.com/oauth2/token",
		Scopes:       []string{"iam:read", "iam:write"},
	}

	token, err := config.Token(context.Background())
	if err != nil {
		fmt.Printf("❌ Failed to get token: %v\n", err)
		return
	}

	fmt.Printf("✅ Got token: %s...\n", token.AccessToken[:50])

	// Test 1: Try to create a group
	fmt.Printf("\n📡 Test 1: Creating IAM Group\n")
	err = testCreateGroup(token.AccessToken, tenantID)
	if err != nil {
		fmt.Printf("❌ Create group failed: %v\n", err)
	} else {
		fmt.Printf("✅ Create group succeeded!\n")
	}

	// Test 2: Try to list groups (read operation)
	fmt.Printf("\n📡 Test 2: Listing IAM Groups\n")
	err = testListGroups(token.AccessToken, tenantID)
	if err != nil {
		fmt.Printf("❌ List groups failed: %v\n", err)
	} else {
		fmt.Printf("✅ List groups succeeded!\n")
	}

	// Test 3: Try to create a custom role
	fmt.Printf("\n📡 Test 3: Creating Custom Role\n")
	err = testCreateCustomRole(token.AccessToken, tenantID)
	if err != nil {
		fmt.Printf("❌ Create custom role failed: %v\n", err)
	} else {
		fmt.Printf("✅ Create custom role succeeded!\n")
	}
}

func testCreateGroup(token, tenantID string) error {
	payload := map[string]interface{}{
		"name":        "test-api-group-" + fmt.Sprintf("%d", time.Now().Unix()),
		"description": "Test group created via direct API call",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://iam-api.retailsvc.com/api/v1/tenants/%s/groups", tenantID), bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("   Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func testListGroups(token, tenantID string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://iam-api.retailsvc.com/api/v1/tenants/%s/groups", tenantID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("   Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body)[:min(200, len(string(body)))])

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func testCreateCustomRole(token, tenantID string) error {
	roleId := "test-api-role-" + fmt.Sprintf("%d", time.Now().Unix())
	payload := map[string]interface{}{
		"id":   roleId,
		"name": "Test API Custom Role",
		"permissions": []map[string]interface{}{
			{"id": "iam.group.list"},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://iam-api.retailsvc.com/api/v1/tenants/%s/roles", tenantID), bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("   Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
