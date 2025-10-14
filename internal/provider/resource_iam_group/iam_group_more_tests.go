package resource_iam_group

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// seqRT returns a sequence of responses/errors for each RoundTrip call.
type seqRT struct {
	calls int
	resps []*http.Response
	errs  []error
}

func (s *seqRT) RoundTrip(req *http.Request) (*http.Response, error) {
	idx := s.calls
	s.calls++
	if idx < len(s.resps) {
		return s.resps[idx], nil
	}
	if idx-len(s.resps) < len(s.errs) {
		return nil, s.errs[idx-len(s.resps)]
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
}

func Test_makeAPIRequest_retries_on_5xx_then_succeeds(t *testing.T) {
	rt := &seqRT{
		resps: []*http.Response{
			{StatusCode: 502, Body: io.NopCloser(strings.NewReader(`bad`))},
			{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))},
		},
	}
	r := &IamGroupResource{
		baseURL:  "https://api.test",
		tenantID: "tenant-x",
		client:   &http.Client{Transport: rt},
	}

	start := time.Now()
	resp, err := r.makeAPIRequest(context.Background(), "GET", "/retry", nil)
	dur := time.Since(start)
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Fatalf("unexpected resp: %v", resp)
	}
	if rt.calls < 2 {
		t.Fatalf("expected at least 2 attempts, got %d", rt.calls)
	}
	// ensure we waited some time for backoff (at least 1s)
	if dur < time.Second {
		t.Fatalf("expected backoff wait, duration: %v", dur)
	}
}

func Test_retryOperation_exhausts(t *testing.T) {
	r := &IamGroupResource{}
	attempts := 0
	op := func() error {
		attempts++
		return errors.New("server error: temporary")
	}
	start := time.Now()
	err := r.retryOperation(context.Background(), op)
	dur := time.Since(start)
	if err == nil || !strings.Contains(err.Error(), "operation failed after") {
		t.Fatalf("expected exhausted retry error, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if dur < 2*time.Second { // backoffs: 1s + 2s between retries ~ >=3s, but tolerantly check >2s
		t.Fatalf("expected some backoff delay, got %v", dur)
	}
}

func Test_isRetryableError_various_messages(t *testing.T) {
	r := &IamGroupResource{}
	cases := []struct {
		msg    string
		expect bool
	}{
		{"timeout while dialing", true},
		{"connection refused by peer", true},
		{"service temporarily unavailable", true},
		{"server error: status 500", true},
		{"authentication failed", false},
		{"random other error", false},
	}
	for _, c := range cases {
		ok := r.isRetryableError(errors.New(c.msg))
		if ok != c.expect {
			t.Fatalf("unexpected isRetryableError(%q) = %v, want %v", c.msg, ok, c.expect)
		}
	}
}

func Test_contains_findInString_edgecases(t *testing.T) {
	if contains("", "") != true {
		t.Fatalf("empty contains empty should be true")
	}
	if contains("a", "") != true {
		t.Fatalf("any contains empty should be true")
	}
	if contains("", "a") != false {
		t.Fatalf("empty contains non-empty should be false")
	}
	if findInString("abc", "d") != false {
		t.Fatalf("expected not found")
	}
}
