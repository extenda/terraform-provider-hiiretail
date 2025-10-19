package auth

import (
	"context"
	"strings"
	"testing"
)

// NewHTTPClient should return an error if given a nil config
func TestNewHTTPClient_NilConfig(t *testing.T) {
	_, err := NewHTTPClient(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error when calling NewHTTPClient with nil config")
	}
}

// ResolveEndpoints should return non-empty URLs for a reasonable tenant
func TestResolveEndpoints_Basic(t *testing.T) {
	authURL, apiURL, err := ResolveEndpoints("test-tenant-123", "")
	if err != nil {
		t.Fatalf("ResolveEndpoints returned error: %v", err)
	}
	if authURL == "" || apiURL == "" {
		t.Fatalf("expected non-empty endpoints, got auth=%q api=%q", authURL, apiURL)
	}
}

// ValidateConfig should return an error when required fields are missing
func TestValidateConfig_MissingClientID(t *testing.T) {
	cfg := &Config{ClientID: "", ClientSecret: "secret", TenantID: "tenant"}
	if err := ValidateConfig(cfg); err == nil {
		t.Fatalf("expected ValidateConfig to return error when ClientID is missing")
	}
}

// ResolveEndpoints respects explicit environment strings
func TestResolveEndpoints_Environment(t *testing.T) {
	authURL, apiURL, err := ResolveEndpoints("tenant-x", "staging")
	if err != nil {
		t.Fatalf("ResolveEndpoints returned error for staging: %v", err)
	}
	if !strings.Contains(authURL, "staging") && !strings.Contains(apiURL, "staging") {
		t.Fatalf("expected staging urls, got auth=%q api=%q", authURL, apiURL)
	}
}
