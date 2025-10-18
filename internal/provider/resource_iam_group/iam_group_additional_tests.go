package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeAPIRequest_HTTP400_MapsToInvalidRequest(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// Transport returns a single 400 response
	rt := &seqRoundTripper{responses: []*http.Response{makeResp(400, `bad request`)}}
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/bad", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid request")
}

func Test_makeAPIRequest_RetryThenSuccess(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// seqRT: first 500, then 502, then 200
	seq := []*http.Response{
		makeResp(500, `err`),
		makeResp(502, `err`),
		makeResp(200, `{"ok":true}`),
	}
	rt := &seqRoundTripper{responses: seq}
	r.client = &http.Client{Transport: rt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(ctx, http.MethodGet, "/retry", nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 200, resp.StatusCode)
}

func Test_makeAPIRequest_InvalidConfig(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// No client or baseURL configured
	r.client = nil
	r.baseURL = ""

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/x", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "resource not properly configured")
}

// (retry exhaustion test is defined in iam_group_more_tests.go)

// Inline capture transport to assert headers without relying on other test files.
type localCaptureRT struct {
	req  *http.Request
	resp *http.Response
}

func (c *localCaptureRT) RoundTrip(req *http.Request) (*http.Response, error) {
	c.req = req
	return c.resp, nil
}

func Test_makeAPIRequest_PostWithBody_SetsContentType_Additional(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	resp := makeResp(201, `{"created":true}`)
	capt := &localCaptureRT{resp: resp}
	r.client = &http.Client{Transport: capt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	body := []byte(`{"name":"example"}`)
	gotResp, err := r.makeAPIRequest(ctx, http.MethodPost, "/create", body)
	require.NoError(t, err)
	require.NotNil(t, gotResp)
	require.Equal(t, 201, gotResp.StatusCode)

	require.NotNil(t, capt.req)
	if ct := capt.req.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}
