package resource_iam_group

import (
    "context"
    "net/http"
    "testing"

    "github.com/stretchr/testify/require"
)

// Test that makeAPIRequest returns an error when resource is not configured
func Test_makeAPIRequest_MissingConfig_New(t *testing.T) {
    ctx := context.Background()
    r := NewIamGroupResource().(*IamGroupResource)

    // Ensure client and baseURL are nil/empty
    r.client = nil
    r.baseURL = ""

    _, err := r.makeAPIRequest(ctx, http.MethodGet, "/x", nil)
    require.Error(t, err)
    require.Contains(t, err.Error(), "missing HTTP client or base URL")
}

// Test that client errors (e.g., network) are mapped via mapHTTPError with status 0
func Test_makeAPIRequest_ClientErrorMapsToGeneric_New(t *testing.T) {
    ctx := context.Background()
    r := NewIamGroupResource().(*IamGroupResource)

    // Create a transport that returns an error by using a nil client and verifying behavior
    // Use a real client but with a transport that returns an error
    r.client = &http.Client{Transport: &errRT2{}}
    r.baseURL = "https://api.test"
    r.tenantID = "t1"

    _, err := r.makeAPIRequest(ctx, http.MethodGet, "/err", nil)
    require.Error(t, err)
    // mapHTTPError wraps the underlying error with an indication (status 0 path)
    require.Contains(t, err.Error(), "unexpected HTTP status")
}

// errRT is a RoundTripper that returns an error for testing purposes
// errRT2 is a RoundTripper that returns an error for testing purposes
type errRT2 struct{}

func (e errRT2) RoundTrip(req *http.Request) (*http.Response, error) {
    return nil, http.ErrHandlerTimeout
}
