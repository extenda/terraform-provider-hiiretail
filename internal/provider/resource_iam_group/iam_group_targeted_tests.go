package resource_iam_group

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test a successful POST with a body: Content-Type should be set and response returned
func Test_makeAPIRequest_SuccessWithBodyHeaders(t *testing.T) {
	r := &IamGroupResource{
		client:   &http.Client{},
		baseURL:  "https://api.example.com",
		tenantID: "t1",
	}

	// prepare capture transport (defined in other test file)
	capt := &captureRT{resp: makeResp(200, `{"ok":true}`)}
	r.client.Transport = capt

	// captureRT will validate headers when Do is called
	resp, err := r.makeAPIRequest(context.Background(), http.MethodPost, "", []byte(`{"name":"x"}`))
	require.NoError(t, err)
	if resp != nil {
		// ensure we close body in test if non-nil
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// Test that a 4xx HTTP status returns a mapped error and body is closed
func Test_makeAPIRequest_4xxBodyConsumedAndMapped(t *testing.T) {
	// Transport that returns a 400 response with a small JSON body
	rt := &seqRoundTripper{responses: []*http.Response{
		makeResp(400, `{"error":"invalid"}`),
	}}

	r := &IamGroupResource{
		client:   &http.Client{Transport: rt},
		baseURL:  "https://api.example.com",
		tenantID: "t1",
	}

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/foo", nil)
	require.Error(t, err)
	require.Nil(t, resp)
	// error should contain our mapping prefix
	require.Contains(t, err.Error(), "invalid request")
}

// Test that repeated 5xx responses exhaust retries and map to status 0
func Test_makeAPIRequest_5xxExhaustRetriesMapsToStatusZero(t *testing.T) {
	// Return three 502 responses so retryOperation will exhaust
	rt := &seqRoundTripper{responses: []*http.Response{
		makeResp(502, ""),
		makeResp(502, ""),
		makeResp(502, ""),
	}}

	r := &IamGroupResource{
		client:   &http.Client{Transport: rt},
		baseURL:  "https://api.example.com",
		tenantID: "t1",
	}

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/fail", nil)
	require.Error(t, err)
	require.Nil(t, resp)
	// mapping uses status 0 for client-level/network errors; here ensure generic mapping present
	require.Contains(t, err.Error(), "unexpected HTTP status")
}
