package resource_iam_group

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that retryOperation eventually fails after exhausting retries when given a retryable error
func Test_retryOperation_ExhaustsRetries(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	err := r.retryOperation(context.Background(), func() error {
		return fmt.Errorf("server error: status 500")
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "operation failed after 3 attempts")
}

// Test that makeAPIRequest retries on 5xx and ultimately returns a mapped error when all retries fail
func Test_makeAPIRequest_All500_ExhaustsAndMaps(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// seqRoundTripper and makeResp are defined in iam_group_api_tests.go
	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{
		makeResp(500, `{"error":"oops"}`),
		makeResp(502, `{"error":"still oops"}`),
		makeResp(503, `{"error":"nope"}`),
	}}}
	r.baseURL = "http://example"
	r.tenantID = "tenant"

	_, err := r.makeAPIRequest(context.Background(), "GET", "/", nil)
	require.Error(t, err)
	// When retryOperation fails it maps to status 0 in makeAPIRequest
	require.Contains(t, err.Error(), "unexpected HTTP status 0")
}

// Test that a 400 response is mapped to an "invalid request" error
func Test_makeAPIRequest_400MapsToInvalidRequest(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: &seqRoundTripper{responses: []*http.Response{
		makeResp(400, `{"message":"bad"}`),
	}}}
	r.baseURL = "http://example"
	r.tenantID = "tenant"

	_, err := r.makeAPIRequest(context.Background(), "GET", "/", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid request")
}
