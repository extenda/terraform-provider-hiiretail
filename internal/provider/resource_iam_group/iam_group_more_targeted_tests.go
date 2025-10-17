package resource_iam_group

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure a 201 Created response is returned as-is and body can be read by caller
func Test_makeAPIRequest_ReturnsCreatedAndBodyAccessible(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	resp := makeResp(201, `{"created":true}`)
	seq := &seqRoundTripper{responses: []*http.Response{resp}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	got, err := r.makeAPIRequest(context.Background(), http.MethodPost, "/create", []byte(`{"name":"x"}`))
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, 201, got.StatusCode)
	// read body via unmarshalResponse to ensure body accessible and closed
	var out map[string]bool
	err = r.unmarshalResponse(got, &out)
	require.NoError(t, err)
	require.True(t, out["created"])
}

// Unexpected 418 should hit the default branch in mapHTTPError via makeAPIRequest
func Test_makeAPIRequest_Unexpected418MapsDefault(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	seq := &seqRoundTripper{responses: []*http.Response{makeResp(418, `nope`)}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	got, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/teapot", nil)
	require.Error(t, err)
	require.Nil(t, got)
	require.Contains(t, err.Error(), "unexpected HTTP status 418")
}

// Call makeAPIRequest concurrently to exercise potential concurrent branches and transport usage
func Test_makeAPIRequest_ConcurrentCalls(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// prepare multiple responses for concurrent calls
	seq := &seqRoundTripper{responses: []*http.Response{
		makeResp(200, `{"ok":true}`),
		makeResp(200, `{"ok":true}`),
		makeResp(200, `{"ok":true}`),
	}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	var wg sync.WaitGroup
	errs := make(chan error, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/con", nil)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		require.NoError(t, e)
	}
}
