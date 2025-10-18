package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// When the server returns a JSON error body with 400, ensure mapping includes message.
func Test_makeAPIRequest_Maps400WithJSONBody(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	body := `{"error":"invalid name"}`
	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(400, body)}}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/badjson", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid request")
}

// Ensure makeAPIRequest retries when response status >=500 and succeeds later.
func Test_makeAPIRequest_502Then200_retries(t *testing.T) {
	// seqRoundTripper will return a 502 then a 200
	seq := &seqRoundTripper{responses: []*http.Response{
		makeResp(502, `temporary`),
		makeResp(200, `{"ok":true}`),
	}}

	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/retry502", nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 200, resp.StatusCode)
	// seqRoundTripper tracks calls internally; ensure at least 2 calls were made
	require.True(t, seq.calls >= 2)
}

// Extra unit tests for contains/findInString edgecases to bump coverage.
func Test_contains_and_findInString_additional(t *testing.T) {
	// substr same as s
	if !contains("abc", "abc") {
		t.Fatalf("expected contains to find equal strings")
	}
	// substr at end
	if !contains("hello world", "world") {
		t.Fatalf("expected contains to find suffix")
	}
	// findInString with multi occurrence
	if !findInString("aaaaab", "ab") {
		t.Fatalf("expected findInString to find substring")
	}
	// negative case
	if findInString("short", "long") {
		t.Fatalf("did not expect to find substring")
	}
}
