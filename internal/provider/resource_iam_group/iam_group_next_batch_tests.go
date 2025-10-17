package resource_iam_group

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// When server returns 400 with a JSON error body, makeAPIRequest should map to invalid request and include body message
func Test_makeAPIRequest_Maps400WithJSONErrorBody(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	seq := &seqRoundTripper{responses: []*http.Response{makeResp(400, `{"message":"bad input"}`)}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/badjson", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad input")
}

// Test longer retry sequence: several 5xx then success to ensure retryOperation exercised
func Test_makeAPIRequest_Several5xxThen200Retries(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	seq := &seqRoundTripper{responses: []*http.Response{
		makeResp(500, `err1`),
		makeResp(502, `err2`),
		makeResp(503, `err3`),
		makeResp(200, `{"ok":true}`),
	}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/many", nil)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// If client.Do returns an error with a custom message, the mapping should include that message
func Test_makeAPIRequest_DoErrorIncludesMessage(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: &errRT{err: errors.New("custom dial error")}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/errmsg", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "custom dial error")
}
