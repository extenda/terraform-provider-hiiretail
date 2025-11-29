package resource_iam_group

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// fakeRoundTripper simulates an http.RoundTripper that always returns an error
type fakeRoundTripper struct{}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("connection refused")
}

func Test_isRetryableError_and_findInString(t *testing.T) {
	r := &IamGroupResource{}

	if r.isRetryableError(nil) {
		t.Fatalf("expected nil error to be non-retryable")
	}

	if !r.isRetryableError(errors.New("timeout while connecting")) {
		t.Fatalf("expected timeout to be retryable")
	}

	if !findInString("hello world", "world") {
		t.Fatalf("expected findInString to find substring")
	}

	if findInString("short", "longer-substr") {
		t.Fatalf("expected findInString to return false when substr longer than source")
	}
}

func Test_makeAPIRequest_MisconfiguredAndDoError(t *testing.T) {
	ctx := context.Background()
	r := &IamGroupResource{}

	// Unconfigured resource should return configuration error
	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/", nil)
	if err == nil {
		t.Fatalf("expected error when resource not configured")
	}

	// Configure resource with fake client that errors via transport
	r.client = &http.Client{Transport: &fakeRoundTripper{}}
	r.baseURL = "https://example.com"
	r.tenantID = "tenant-x"

	_, err = r.makeAPIRequest(ctx, http.MethodGet, "/", nil)
	if err == nil {
		t.Fatalf("expected error when underlying client Do returns error")
	}
}
