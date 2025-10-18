package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeAPIRequest_Various4xxMap(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	cases := []struct {
		code int
		want string
	}{
		{401, "authentication failed"},
		{403, "access denied"},
		{409, "group already exists"},
		{400, "invalid request"},
		{418, "unexpected HTTP status 418"},
	}

	for _, tc := range cases {
		r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(tc.code, `err`)}}}
		r.baseURL = "https://api.test"
		r.tenantID = "t1"

		_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/code", nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), tc.want)
	}
}

func Test_makeAPIRequest_NewRequestParseError(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// malformed base URL should cause NewRequestWithContext to return an error
	r.baseURL = "http://[::1"
	r.tenantID = "t1"
	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{makeResp(200, `ok`)}}}

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/badurl", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}
