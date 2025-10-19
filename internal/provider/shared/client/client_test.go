package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	authpkg "github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
)

// fake RoundTripper that returns a simple 200 response
type fakeRT struct{}

func (f fakeRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Header:     make(http.Header),
	}, nil
}

func TestBuildURLAndBase(t *testing.T) {
	cfg := DefaultConfig()
	ac := &authpkg.Config{TenantID: "tenant-1", TestToken: "tok"}
	c, err := New(ac, cfg)
	if err != nil {
		t.Fatalf("New client: %v", err)
	}

	// build a url and ensure the base and tenant appear
	u := c.buildURL("/some/path")
	if u.String() == "" {
		t.Fatalf("expected URL, got empty")
	}

	// ensure IAMClient wraps and doesn't panic
	sc := c.IAMClient()
	if sc == nil || sc.client == nil {
		t.Fatalf("expected service client")
	}

	// shouldRetry basic cases
	if !c.shouldRetry(500) {
		t.Fatalf("expected retry for 500")
	}
	if c.shouldRetry(200) {
		t.Fatalf("didn't expect retry for 200")
	}

	// basic checks
	if c.TenantID() != ac.TenantID {
		t.Fatalf("tenant mismatch: got %s", c.TenantID())
	}
	if c.BaseURL() == "" {
		t.Fatalf("base URL empty")
	}

	// install fake transport to avoid real network calls
	c.httpClient.Transport = fakeRT{}

	// call Get which will use the fake RoundTripper
	resp, err := sc.Get(context.Background(), "/ping", nil)
	if err != nil {
		t.Fatalf("unexpected error from Get: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
