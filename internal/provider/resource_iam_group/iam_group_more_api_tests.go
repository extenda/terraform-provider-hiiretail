package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure 404 responses map to group not found error via makeAPIRequest.
func Test_makeAPIRequest_Maps404ToNotFound(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(404, `not found`)}}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/missing", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "group not found")
}

// When the server returns >=400, makeAPIRequest should map the error even if body is plain text.
func Test_makeAPIRequest_Maps400WithPlainBody(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(400, `bad request - plain text`)}}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/bad", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid request")
}

// unmarshalResponse should return an error when JSON is invalid.
func Test_unmarshalResponse_invalidJSON(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	resp := makeResp(200, `not-a-json`)
	var target map[string]interface{}
	err := r.unmarshalResponse(resp, &target)
	require.Error(t, err)
}

// marshalRequest should return an error when data is not JSON-marshalable (e.g., channel).
func Test_marshalRequest_error(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	_, err := r.marshalRequest(make(chan int))
	require.Error(t, err)
}
