package resource_iam_group

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type fakeRT struct {
	resp *http.Response
	err  error
}

func (f *fakeRT) RoundTrip(req *http.Request) (*http.Response, error) { return f.resp, f.err }

func Test_makeAPIRequest_SuccessAnd404(t *testing.T) {
	r := &IamGroupResource{
		baseURL:  "https://api.test",
		tenantID: "tenant-x",
		client:   &http.Client{Transport: &fakeRT{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}}},
	}

	resp, err := r.makeAPIRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Fatalf("expected 200 response, got %v %v", resp, err)
	}

	// Now simulate 404
	r.client = &http.Client{Transport: &fakeRT{resp: &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(`{}`))}}}
	resp2, err2 := r.makeAPIRequest(context.Background(), "GET", "/notfound", nil)
	if err2 == nil || !strings.Contains(err2.Error(), "group not found") {
		t.Fatalf("expected group not found error on 404, got %v (resp=%v)", err2, resp2)
	}
}

func Test_retryOperation_nonRetryableImmediate(t *testing.T) {
	r := &IamGroupResource{}
	err := r.retryOperation(context.Background(), func() error {
		return fmt.Errorf("permanent failure")
	})
	if err == nil || !strings.Contains(err.Error(), "permanent failure") {
		t.Fatalf("expected immediate permanent failure, got %v", err)
	}
}

func Test_retryOperation_singleRetry_then_success(t *testing.T) {
	r := &IamGroupResource{}
	calls := 0
	op := func() error {
		calls++
		if calls == 1 {
			return fmt.Errorf("server error")
		}
		return nil
	}

	start := time.Now()
	err := r.retryOperation(context.Background(), op)
	duration := time.Since(start)
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	// Ensure we waited at least 1s because of the backoff between attempts
	if duration < time.Second {
		t.Fatalf("expected at least 1s duration due to backoff, got %v", duration)
	}
}
