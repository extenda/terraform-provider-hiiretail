package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test various HTTP error/status branches for create/read/update/delete helpers
func TestAPIHelpers_ErrorAndEdgeStatuses(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// 1) readCustomRole -> 404 should return "custom role not found"
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	_, err := r.readCustomRole(ctx, "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "custom role not found")

	// 2) deleteCustomRole -> 404 should be treated as not found error
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	err = r.deleteCustomRole(ctx, "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "custom role not found")

	// 3) deleteCustomRole -> unexpected status with malformed JSON body
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusTeapot, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("not-json")))}, nil
	}}}

	err = r.deleteCustomRole(ctx, "x")
	require.Error(t, err)
	require.Contains(t, err.Error(), "API request failed with status")

	// 4) createCustomRole -> non-201 with valid JSON error body
	errResp := ErrorResponse{Message: "bad request", Code: "BAD"}
	eb, _ := json.Marshal(errResp)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewBuffer(eb))}, nil
	}}}

	_, err = r.createCustomRole(ctx, &CustomRoleRequest{ID: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "API error")

	// 5) updateCustomRole -> 404 should return "custom role not found"
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	_, err = r.updateCustomRole(ctx, "x", &CustomRoleRequest{ID: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "custom role not found")
}
