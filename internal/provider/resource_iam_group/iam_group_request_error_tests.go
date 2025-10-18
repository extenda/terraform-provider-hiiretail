package resource_iam_group

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Body tracker to detect Close calls
type trackCloser struct {
	closed bool
}

func (t *trackCloser) Read(p []byte) (int, error) { return 0, io.EOF }
func (t *trackCloser) Close() error               { t.closed = true; return nil }

func Test_makeAPIRequest_NewRequestError(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// invalid baseURL to cause http.NewRequestWithContext to fail
	r.baseURL = ":"
	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(200, `ok`)}}}
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/x", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected HTTP status 0")
}

func Test_makeAPIRequest_ClosesBodyOn5xx(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	tc := &trackCloser{}
	resp1 := &http.Response{StatusCode: 502, Body: tc}
	resp2 := makeResp(200, `{"ok":true}`)

	seq := &seqRoundTripper{responses: []*http.Response{resp1, resp2}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	got, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/close", nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, tc.closed, "expected response body to be closed on 5xx")
}
