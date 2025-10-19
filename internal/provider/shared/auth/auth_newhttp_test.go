package auth

import (
	"context"
	"testing"
)

func TestNewHTTPClient_WithTestToken(t *testing.T) {
	cfg := &Config{TenantID: "t1", TestToken: "dummy"}
	_, err := NewHTTPClient(context.Background(), cfg)
	// NewHTTPClient validates the config and will return a configuration error
	// when required fields (like client_id) are missing. Ensure we get an error.
	if err == nil {
		t.Fatalf("expected configuration error creating http client with missing client_id, got nil")
	}
}
