package resource_iam_custom_role

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestModelToAPIRequest_PermissionsElementsAsError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build a permissions list with the wrong inner type (strings instead of PermissionsValue)
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	// Construct a list using raw strings which will not convert to PermissionsValue
	badList, _ := types.ListValueFrom(ctx, permType, []interface{}{"not-a-perm"})

	data := IamCustomRoleModel{Id: types.StringValue("x"), Permissions: badList}

	_, err := r.modelToAPIRequest(ctx, data)
	require.Error(t, err)
}

func TestCreateReadUpdate_MalformedJSONResponse(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// createCustomRole: return 201 with invalid JSON
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		if req.Method == "POST" {
			return &http.Response{StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("not-json")))}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	_, err := r.createCustomRole(ctx, &CustomRoleRequest{ID: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode response")

	// readCustomRole: return 200 with invalid JSON
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("not-json")))}, nil
	}}}

	_, err = r.readCustomRole(ctx, "x")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode response")

	// updateCustomRole: return 200 with invalid JSON
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("not-json")))}, nil
	}}}

	_, err = r.updateCustomRole(ctx, "x", &CustomRoleRequest{ID: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode response")
}
