package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUpdate_APIErrorAndMalformedResponse(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	// 1) network error during update
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return nil, http.ErrHandlerTimeout
	}}}

	data := IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list}
	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq.Plan.Set(ctx, data)

	var uresp resource.UpdateResponse
	r.Update(ctx, ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError())

	// 2) API returns 400 with JSON error
	errResp := ErrorResponse{Message: "bad", Code: "BAD"}
	eb, _ := json.Marshal(errResp)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewBuffer(eb))}, nil
	}}}

	var ureq2 resource.UpdateRequest
	ureq2.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq2.Plan.Set(ctx, IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list})
	var uresp2 resource.UpdateResponse
	r.Update(ctx, ureq2, &uresp2)
	require.True(t, uresp2.Diagnostics.HasError())

	// 3) API returns 200 but malformed JSON
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("not-json")))}, nil
	}}}

	var ureq3 resource.UpdateRequest
	ureq3.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq3.Plan.Set(ctx, IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list})
	var uresp3 resource.UpdateResponse
	r.Update(ctx, ureq3, &uresp3)
	require.True(t, uresp3.Diagnostics.HasError())
}
