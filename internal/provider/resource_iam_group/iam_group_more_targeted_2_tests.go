package resource_iam_group

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Validate unmarshalResponse returns an error on invalid JSON bodies
func Test_unmarshalResponse_InvalidJSON(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// makeResp uses application/json header and returns Body
	resp := makeResp(200, `{ not-json `)
	err := r.unmarshalResponse(resp, &map[string]interface{}{})
	require.Error(t, err)
}

// Explicitly test mapHTTPError branches for 401, 403, 409
func Test_mapHTTPError_AuthForbiddenConflict(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	e := r.mapHTTPError(401, errors.New("x"))
	require.Error(t, e)
	require.Contains(t, e.Error(), "authentication failed")

	e = r.mapHTTPError(403, errors.New("x"))
	require.Error(t, e)
	require.Contains(t, e.Error(), "access denied")

	e = r.mapHTTPError(409, errors.New("x"))
	require.Error(t, e)
	require.Contains(t, e.Error(), "group already exists")
}

// Check that when a body is provided the Content-Type header is set (captureRT used)
func Test_makeAPIRequest_BodySetsContentType_Capture(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	resp := makeResp(200, `{"ok":true}`)
	capt := &captureRT{resp: resp}
	r.client = &http.Client{Transport: capt}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodPut, "/x", []byte(`{"a":1}`))
	require.NoError(t, err)
	require.NotNil(t, capt.req)
	require.Equal(t, "application/json", capt.req.Header.Get("Content-Type"))
}
