package resource_iam_group

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ... existing code ...
func Test_makeAPIRequest_DoError_ConnectionRefused(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: &errRT{err: fmt.Errorf("connection refused by peer")}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/x", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "connection refused")
}

func Test_retryOperation_BackoffThenSuccess(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	calls := 0
	err := r.retryOperation(context.Background(), func() error {
		calls++
		if calls == 1 {
			return errors.New("server error: status 500")
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}
