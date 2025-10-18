package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestCreateReadUpdateDeleteCustomRole_HTTPPaths(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// Create response
	created := CustomRoleResponse{ID: "c1", Name: "C1", TenantID: "tid", Permissions: []Permission{{ID: "p1"}}}
	cb, _ := json.Marshal(created)

	// Read response
	read := created
	rb, _ := json.Marshal(read)

	// Update response (same as created)
	ub, _ := json.Marshal(created)

	// Wire mock client to respond accordingly
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == "POST" && req.URL.Path == "/api/v1/tenants/tid/roles":
			return &http.Response{StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewBuffer(cb))}, nil
		case req.Method == "GET" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(rb))}, nil
		case req.Method == "PUT" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(ub))}, nil
		case req.Method == "DELETE" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusNoContent, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
		default:
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("err")))}, nil
		}
	}}}

	// modelToAPIRequest happy path
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	data := IamCustomRoleModel{Id: types.StringValue("c1"), Name: types.StringValue("C1"), Permissions: list}

	apiReq, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Equal(t, "c1", apiReq.ID)

	// createCustomRole
	cr, err := r.createCustomRole(ctx, apiReq)
	require.NoError(t, err)
	require.Equal(t, "c1", cr.ID)

	// readCustomRole
	rr, err := r.readCustomRole(ctx, "c1")
	require.NoError(t, err)
	require.Equal(t, "c1", rr.ID)

	// updateCustomRole
	ur, err := r.updateCustomRole(ctx, "c1", apiReq)
	require.NoError(t, err)
	require.Equal(t, "c1", ur.ID)

	// deleteCustomRole
	err = r.deleteCustomRole(ctx, "c1")
	require.NoError(t, err)
}

func TestAPIErrorPaths(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// mock client returns 400 with JSON error
	errResp := ErrorResponse{Message: "bad", Code: "BAD"}
	eb, _ := json.Marshal(errResp)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewBuffer(eb))}, nil
	}}}

	// createCustomRole should return API error
	_, err := r.createCustomRole(ctx, &CustomRoleRequest{ID: "x"})
	require.Error(t, err)

	// readCustomRole should return API error
	_, err = r.readCustomRole(ctx, "x")
	require.Error(t, err)

	// updateCustomRole should return API error
	_, err = r.updateCustomRole(ctx, "x", &CustomRoleRequest{ID: "x"})
	require.Error(t, err)

	// deleteCustomRole should return API error
	err = r.deleteCustomRole(ctx, "x")
	require.Error(t, err)
}
