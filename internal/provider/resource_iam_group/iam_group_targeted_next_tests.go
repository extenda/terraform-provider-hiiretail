package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure Content-Type is set when body is provided (exercise body branch)
func Test_makeAPIRequest_PostBodySetsContentType_Capture(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	resp := makeResp(201, `{"created":true}`)
	capt := &captureRT{resp: resp}
	r.client = &http.Client{Transport: capt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodPost, "/create", []byte(`{"name":"x"}`))
	require.NoError(t, err)
	require.NotNil(t, capt.req)
	require.Equal(t, "application/json", capt.req.Header.Get("Content-Type"))
}

// When the server returns 400 plain text, ensure mapping still occurs
func Test_makeAPIRequest_400PlainTextMaps(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	seq := &seqRoundTripper{responses: []*http.Response{makeResp(400, `bad request`)}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/bad", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid request")
}

// If resource not configured (client nil), ensure early error path is taken
func Test_makeAPIRequest_MissingConfigError(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = nil
	r.baseURL = ""

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/x", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "resource not properly configured")
}
