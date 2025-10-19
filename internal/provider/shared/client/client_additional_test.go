package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	authpkg "github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
)

// fakeRoundTripper simulates transient failures then success
type fakeRoundTripper struct {
	calls int
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.calls++
	if f.calls <= 1 {
		return nil, errors.New("simulated network error")
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Header:     make(http.Header),
	}, nil
}

func TestCalculateBackoffBounds(t *testing.T) {
	cfg := DefaultConfig()
	ac := &authpkg.Config{TenantID: "t", TestToken: "tok"}
	c, err := New(ac, cfg)
	if err != nil {
		t.Fatalf("New client error: %v", err)
	}

	// small attempt should be at least min
	d := c.calculateBackoff(1)
	if d < cfg.RetryWaitMin {
		t.Fatalf("backoff too small: %v", d)
	}

	// large attempt should be capped at max
	d = c.calculateBackoff(20)
	if d > cfg.RetryWaitMax {
		t.Fatalf("backoff exceeded max: %v", d)
	}
}

func TestDoWithRetry_NetworkErrorThenSuccess(t *testing.T) {
	cfg := DefaultConfig()
	ac := &authpkg.Config{TenantID: "t", TestToken: "tok"}
	c, err := New(ac, cfg)
	if err != nil {
		t.Fatalf("New client error: %v", err)
	}

	// Install fake transport
	frt := &fakeRoundTripper{}
	c.httpClient.Transport = frt

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, c.buildURL("/ping").String(), nil)

	// This should retry once then succeed
	resp, err := c.doWithRetry(context.Background(), req)
	if err != nil {
		t.Fatalf("doWithRetry failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}
