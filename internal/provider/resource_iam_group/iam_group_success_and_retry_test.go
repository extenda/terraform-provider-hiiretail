package resource_iam_group

import (
	"context"
	"net/http"
	"testing"
)

// Use existing helpers (seqRoundTripper, makeResp) defined in other test files.
func Test_makeAPIRequest_ExtraSuccess(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	rt := &seqRoundTripper{responses: []*http.Response{makeResp(200, `{"ok":true}`)}}
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func Test_makeAPIRequest_ExtraRetry(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	seq := []*http.Response{
		makeResp(500, `err`),
		makeResp(502, `err`),
		makeResp(200, `{"ok":true}`),
	}
	rt := &seqRoundTripper{responses: seq}
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
