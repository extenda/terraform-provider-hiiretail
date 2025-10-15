package resource_iam_group

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// seqRoundTripper returns a sequence of predefined responses for successive RoundTrip calls.
type seqRoundTripper struct {
	mu        sync.Mutex
	responses []*http.Response
	calls     int
}

func (s *seqRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.calls
	if idx >= len(s.responses) {
		idx = len(s.responses) - 1
	}
	s.calls++
	// Return a copy so callers can Close the Body independently
	resp := s.responses[idx]
	return resp, nil
}

func makeResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": {"application/json"}},
	}
}

func Test_makeAPIRequest_SuccessAnd404AndRetry(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// 1) Success path
	rt1 := &seqRoundTripper{responses: []*http.Response{makeResp(200, `{"ok":true}`)}}
	r.client = &http.Client{Transport: rt1}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(ctx, http.MethodGet, "/123", nil)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	// exercise unmarshalResponse success
	var payload map[string]bool
	err = r.unmarshalResponse(resp, &payload)
	require.NoError(t, err)
	require.True(t, payload["ok"])

	// 2) 404 should be mapped to a not found style error
	rt2 := &seqRoundTripper{responses: []*http.Response{makeResp(404, `not found`)}}
	r.client = &http.Client{Transport: rt2}
	resp, err = r.makeAPIRequest(ctx, http.MethodGet, "/notfound", nil)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "group not found")

	// 3) Retry on 5xx then success
	seq := []*http.Response{
		makeResp(500, `err`),
		makeResp(502, `err`),
		makeResp(200, `{"ok":true}`),
	}
	rt3 := &seqRoundTripper{responses: seq}
	r.client = &http.Client{Transport: rt3}
	resp, err = r.makeAPIRequest(ctx, http.MethodGet, "/retry", nil)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func Test_unmarshal_and_marshal_errors_and_helpers(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	// unmarshal invalid JSON
	bad := makeResp(200, `{invalid json`)
	err := r.unmarshalResponse(bad, &map[string]interface{}{})
	require.Error(t, err)

	// marshal success
	data := map[string]string{"key": "value"}
	body, err := r.marshalRequest(data)
	require.NoError(t, err)
	require.Contains(t, string(body), `"key":"value"`)

	// contains and findInString edge cases
	require.True(t, contains("hello world", "world"))
	require.False(t, contains("hello", "world"))
	require.True(t, findInString("hello world", "world"))
	require.False(t, findInString("hello", "world"))

	// isRetryableError
	require.True(t, r.isRetryableError(errors.New("timeout occurred")))
	require.True(t, r.isRetryableError(errors.New("connection refused")))
	require.True(t, r.isRetryableError(errors.New("server error")))
	require.False(t, r.isRetryableError(errors.New("not found")))
}

func Test_mapHTTPError_AllCases(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	testCases := []struct {
		statusCode int
		expected   string
	}{
		{404, "group not found"},
		{401, "authentication failed"},
		{403, "access denied"},
		{409, "group already exists"},
		{400, "invalid request"},
		{500, "server error"},
		{503, "service temporarily unavailable"},
		{418, "unexpected HTTP status 418"}, // default case
	}

	for _, tc := range testCases {
		err := r.mapHTTPError(tc.statusCode, errors.New("test"))
		require.Error(t, err)
		require.Contains(t, err.Error(), tc.expected)
	}
}
