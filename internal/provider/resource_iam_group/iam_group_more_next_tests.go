package resource_iam_group

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Verify Accept and User-Agent headers are set on requests with no body.
func Test_makeAPIRequest_SetsAcceptAndUserAgent(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	resp := makeResp(200, `{"ok":true}`)
	capt := &captureRT{resp: resp}
	r.client = &http.Client{Transport: capt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	got, err := r.makeAPIRequest(ctx, http.MethodGet, "/hdrs", nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, capt.req)
	require.Equal(t, "application/json", capt.req.Header.Get("Accept"))
	require.Contains(t, capt.req.Header.Get("User-Agent"), "terraform-provider-hiiretail-iam")
}

// Ensure multiple 5xx responses trigger retries and eventual success when a 200 appears.
func Test_makeAPIRequest_Multiple5xxThen200(t *testing.T) {
	seq := &seqRoundTripper{responses: []*http.Response{
		makeResp(500, `err1`),
		makeResp(503, `err2`),
		makeResp(200, `{"ok":true}`),
	}}
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/multi", nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 200, resp.StatusCode)
	require.True(t, seq.calls >= 3)
}

// If client.Do returns an immediate non-retryable error, makeAPIRequest should map it with status 0.
type errRT struct {
	err error
}

func (e *errRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, e.err
}

func Test_makeAPIRequest_DoErrorMapsToStatusZero(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: &errRT{err: errors.New("dial tcp: connection refused")}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/err", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected HTTP status 0")
}
