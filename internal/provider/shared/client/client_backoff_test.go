package client

import (
	"testing"

	authpkg "github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
)

func TestShouldRetry_StatusCodes(t *testing.T) {
	cfg := DefaultConfig()
	ac := &authpkg.Config{TenantID: "t", TestToken: "tok"}
	c, err := New(ac, cfg)
	if err != nil {
		t.Fatalf("New client error: %v", err)
	}

	// Status codes that should trigger a retry
	retryCodes := []int{429, 502, 503, 504}
	for _, code := range retryCodes {
		if !c.shouldRetry(code) {
			t.Fatalf("expected shouldRetry to be true for status %d", code)
		}
	}

	// Status codes that should not trigger a retry
	okCodes := []int{200, 201, 400, 401, 403}
	for _, code := range okCodes {
		if c.shouldRetry(code) {
			t.Fatalf("expected shouldRetry to be false for status %d", code)
		}
	}
}

func TestCalculateBackoff_MinMax(t *testing.T) {
	cfg := DefaultConfig()
	ac := &authpkg.Config{TenantID: "t", TestToken: "tok"}
	c, err := New(ac, cfg)
	if err != nil {
		t.Fatalf("New client error: %v", err)
	}

	// Low attempt should be >= min
	if d := c.calculateBackoff(1); d < cfg.RetryWaitMin {
		t.Fatalf("backoff too small: %v", d)
	}

	// Very large attempt should be <= max
	if d := c.calculateBackoff(1000); d > cfg.RetryWaitMax {
		t.Fatalf("backoff exceeded max: %v", d)
	}
}
