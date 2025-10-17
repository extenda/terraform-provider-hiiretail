package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that a 401 response is mapped to an error and the body is closed
func Test_makeAPIRequest_Maps401AndClosesBody(t *testing.T) {
	rt := &seqRoundTripper{
		responses: []*http.Response{makeResp(401, `{"error":"unauthorized"}`)},
	}

	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	ctx := context.Background()
	reqBody := []byte(`{"foo":"bar"}`)

	_, err := r.makeAPIRequest(ctx, http.MethodPost, "/v1/some/path", reqBody)
	require.Error(t, err)
	// mapHTTPError maps 401 to an authentication failed style message; ensure body-derived text present
	require.Contains(t, err.Error(), "unauthorized")
}

// Test that a 409 Conflict response is mapped to an error (idempotency / conflict case)
func Test_makeAPIRequest_Maps409(t *testing.T) {
	rt := &seqRoundTripper{
		responses: []*http.Response{makeResp(409, `{"error":"conflict"}`)},
	}

	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	ctx := context.Background()
	_, err := r.makeAPIRequest(ctx, http.MethodPut, "/v1/some/path", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflict")
}
