package resource_iam_resource

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
)

// mockRawClient implements iam.RawClient for testing
type mockRawClient struct{}

func (m *mockRawClient) Do(ctx context.Context, req *client.Request) (*client.Response, error) {
	// Simulate a successful response for SetResource, GetResource, DeleteResource
	// Extract ID from path, e.g., /api/v1/tenants/test-tenant/resources/{id}
	var id string
	if strings.Contains(req.Path, "/resources/") {
		parts := strings.Split(req.Path, "/resources/")
		if len(parts) > 1 {
			id = strings.Split(parts[1], "/")[0]
		}
	}
	if id == "" {
		id = "test-resource"
	}

	switch req.Method {
	case "PUT":
		// SetResource: return a Resource JSON, use name from body, ID from path
		var name string
		if req.Body != nil {
			if dto, ok := req.Body.(*iam.SetResourceDto); ok {
				name = dto.Name
			} else {
				name = "test-resource"
			}
		} else {
			name = "test-resource"
		}
		body, _ := json.Marshal(&iam.Resource{ID: id, Name: name, Props: map[string]interface{}{"key": "value"}})
		return &client.Response{Body: body, StatusCode: 200}, nil
	case "GET":
		// GetResource: return a Resource JSON
		body, _ := json.Marshal(&iam.Resource{ID: id, Name: "test-resource", Props: map[string]interface{}{"key": "value"}})
		return &client.Response{Body: body, StatusCode: 200}, nil
	case "DELETE":
		// DeleteResource: return empty response
		return &client.Response{StatusCode: 200}, nil
	}
	return &client.Response{StatusCode: 200}, nil
}

// mockRawClientNotFound returns 404 for delete operations
type mockRawClientNotFound struct{}

func (m *mockRawClientNotFound) Do(ctx context.Context, req *client.Request) (*client.Response, error) {
	if req.Method == "DELETE" {
		return &client.Response{StatusCode: 404, Body: []byte(`{"error": "Resource not found"}`)}, nil
	}
	return &client.Response{StatusCode: 200}, nil
}

// mockRawClientError returns an error for delete operations
type mockRawClientError struct{}

func (m *mockRawClientError) Do(ctx context.Context, req *client.Request) (*client.Response, error) {
	if req.Method == "DELETE" {
		return &client.Response{StatusCode: 500, Body: []byte(`{"error": "Internal server error"}`)}, nil
	}
	return &client.Response{StatusCode: 200}, nil
}

func newTestService() *iam.Service {
	s := &iam.Service{}
	// Use unsafe to set private fields if needed
	// For this test, set rawClient to mockRawClient
	v := reflect.ValueOf(s).Elem()
	rawClientField := v.FieldByName("rawClient")
	rawClientPtr := unsafe.Pointer(rawClientField.UnsafeAddr())
	*(*iam.RawClient)(rawClientPtr) = &mockRawClient{}

	tenantIDField := v.FieldByName("tenantID")
	tenantIDPtr := unsafe.Pointer(tenantIDField.UnsafeAddr())
	*(*string)(tenantIDPtr) = "test-tenant"

	return s
}

func newTestServiceNotFound() *iam.Service {
	s := &iam.Service{}
	v := reflect.ValueOf(s).Elem()
	rawClientField := v.FieldByName("rawClient")
	rawClientPtr := unsafe.Pointer(rawClientField.UnsafeAddr())
	*(*iam.RawClient)(rawClientPtr) = &mockRawClientNotFound{}

	tenantIDField := v.FieldByName("tenantID")
	tenantIDPtr := unsafe.Pointer(tenantIDField.UnsafeAddr())
	*(*string)(tenantIDPtr) = "test-tenant"

	return s
}

func newTestServiceError() *iam.Service {
	s := &iam.Service{}
	v := reflect.ValueOf(s).Elem()
	rawClientField := v.FieldByName("rawClient")
	rawClientPtr := unsafe.Pointer(rawClientField.UnsafeAddr())
	*(*iam.RawClient)(rawClientPtr) = &mockRawClientError{}

	tenantIDField := v.FieldByName("tenantID")
	tenantIDPtr := unsafe.Pointer(tenantIDField.UnsafeAddr())
	*(*string)(tenantIDPtr) = "test-tenant"

	return s
}

func setServiceField(r *IAMResourceResource, service interface{}) {
	v := reflect.ValueOf(r).Elem()
	field := v.FieldByName("service")
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	svc := service.(*iam.Service)
	*(*uintptr)(fieldPtr) = uintptr(unsafe.Pointer(svc))
}

func TestIAMResource_MetadataAndSchema(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	var mr resource.MetadataResponse
	r.Metadata(nil, resource.MetadataRequest{ProviderTypeName: "hiiretail"}, &mr)
	require.Contains(t, mr.TypeName, "hiiretail_iam_resource")

	var sr resource.SchemaResponse
	r.Schema(nil, resource.SchemaRequest{}, &sr)
	require.NotNil(t, sr.Schema)
}

func TestIAMResource_Configure(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)

	// Invalid provider data
	var cr resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "invalid"}, &cr)
	require.True(t, cr.Diagnostics.HasError())
}

func TestIAMResource_Create_Success(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	// Fail-fast check: ensure service is set to *mockService
	if r.service == nil {
		t.Fatalf("service field was not set; test cannot proceed")
	}

	data := IAMResourceResourceModel{
		ID:       types.StringNull(), // empty ID to test generation
		Name:     types.StringValue("test-resource"),
		Props:    types.StringValue(`{"key":"value"}`),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var creq resource.CreateRequest
	creq.Plan.Schema = schema
	diags := creq.Plan.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var cresp resource.CreateResponse
	cresp.State.Schema = schema

	r.Create(context.Background(), creq, &cresp)
	if cresp.Diagnostics.HasError() {
		t.Logf("Create diagnostics: %v", cresp.Diagnostics.Errors())
	}
	require.False(t, cresp.Diagnostics.HasError())

	var out IAMResourceResourceModel
	diags = cresp.State.Get(context.Background(), &out)
	require.False(t, diags.HasError())
	require.Equal(t, "test-resource", out.ID.ValueString())
	require.Equal(t, "test-resource", out.Name.ValueString())
	require.Equal(t, "test-tenant", out.TenantID.ValueString())
}

func TestIAMResource_Create_InvalidProps(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("test-resource"),
		Props:    types.StringValue(`invalid json`), // invalid JSON
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var creq resource.CreateRequest
	creq.Plan.Schema = schema
	diags := creq.Plan.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var cresp resource.CreateResponse
	cresp.State.Schema = schema

	r.Create(context.Background(), creq, &cresp)
	require.True(t, cresp.Diagnostics.HasError()) // expect error due to invalid JSON
}

func TestIAMResource_Read_Success(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("test-resource"),
		Props:    types.StringNull(),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var rreq resource.ReadRequest
	rreq.State.Schema = schema
	diags := rreq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var rresp resource.ReadResponse
	rresp.State.Schema = schema

	r.Read(context.Background(), rreq, &rresp)
	if rresp.Diagnostics.HasError() {
		t.Logf("Read diagnostics: %v", rresp.Diagnostics.Errors())
	}
	require.False(t, rresp.Diagnostics.HasError())

	var out IAMResourceResourceModel
	diags = rresp.State.Get(context.Background(), &out)
	require.False(t, diags.HasError())
	require.Equal(t, "test-id", out.ID.ValueString())
	require.Equal(t, "test-resource", out.Name.ValueString())
	require.Equal(t, "test-tenant", out.TenantID.ValueString())
}

func TestIAMResource_Update_Success(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("updated-resource"),
		Props:    types.StringValue(`{"updated":"true"}`),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var ureq resource.UpdateRequest
	ureq.Plan.Schema = schema
	diags := ureq.Plan.Set(context.Background(), data)
	require.False(t, diags.HasError())
	ureq.State.Schema = schema
	diags = ureq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var uresp resource.UpdateResponse
	uresp.State.Schema = schema

	r.Update(context.Background(), ureq, &uresp)
	if uresp.Diagnostics.HasError() {
		t.Logf("Update diagnostics: %v", uresp.Diagnostics.Errors())
	}
	require.False(t, uresp.Diagnostics.HasError())

	var out IAMResourceResourceModel
	diags = uresp.State.Get(context.Background(), &out)
	require.False(t, diags.HasError())
	require.Equal(t, "test-id", out.ID.ValueString())
	require.Equal(t, "updated-resource", out.Name.ValueString())
	require.Equal(t, "test-tenant", out.TenantID.ValueString())
}

func TestIAMResource_Update_InvalidProps(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("updated-resource"),
		Props:    types.StringValue(`invalid json`), // invalid JSON
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var ureq resource.UpdateRequest
	ureq.Plan.Schema = schema
	diags := ureq.Plan.Set(context.Background(), data)
	require.False(t, diags.HasError())
	ureq.State.Schema = schema
	diags = ureq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var uresp resource.UpdateResponse
	uresp.State.Schema = schema

	r.Update(context.Background(), ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError()) // expect error due to invalid JSON
}

func TestIAMResource_Delete_Success(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("test-resource"),
		Props:    types.StringNull(),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var dreq resource.DeleteRequest
	dreq.State.Schema = schema
	diags := dreq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(context.Background(), dreq, &dresp)
	if dresp.Diagnostics.HasError() {
		t.Logf("Delete diagnostics: %v", dresp.Diagnostics.Errors())
	}
	require.False(t, dresp.Diagnostics.HasError())
}

func TestIAMResource_Delete_EmptyID(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestService())

	data := IAMResourceResourceModel{
		ID:       types.StringValue(""), // empty ID
		Name:     types.StringValue("test-resource"),
		Props:    types.StringNull(),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var dreq resource.DeleteRequest
	dreq.State.Schema = schema
	diags := dreq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(context.Background(), dreq, &dresp)
	require.True(t, dresp.Diagnostics.HasError())
	require.Contains(t, dresp.Diagnostics.Errors()[0].Detail(), "Resource ID is required")
}

func TestIAMResource_Delete_AlreadyGone(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestServiceNotFound())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("test-resource"),
		Props:    types.StringNull(),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var dreq resource.DeleteRequest
	dreq.State.Schema = schema
	diags := dreq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(context.Background(), dreq, &dresp)
	// Should not have errors for 404 (resource already gone)
	require.False(t, dresp.Diagnostics.HasError())
}

func TestIAMResource_Delete_APIError(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)
	setServiceField(r, newTestServiceError())

	data := IAMResourceResourceModel{
		ID:       types.StringValue("test-id"),
		Name:     types.StringValue("test-resource"),
		Props:    types.StringNull(),
		TenantID: types.StringNull(),
	}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	var dreq resource.DeleteRequest
	dreq.State.Schema = schema
	diags := dreq.State.Set(context.Background(), data)
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(context.Background(), dreq, &dresp)
	// Should have errors for API failures (non-404)
	require.True(t, dresp.Diagnostics.HasError())
}

func TestIAMResource_ImportState_InvalidID(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)

	var ireq resource.ImportStateRequest
	ireq.ID = "" // invalid empty ID

	var iresp resource.ImportStateResponse
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	iresp.State = tfsdk.State{Schema: schemaResp.Schema}

	r.ImportState(context.Background(), ireq, &iresp)
	require.True(t, iresp.Diagnostics.HasError())
}

func TestIAMResource_ImportState_InvalidIDFormat(t *testing.T) {
	r := NewIAMResourceResource().(*IAMResourceResource)

	var ireq resource.ImportStateRequest
	ireq.ID = "invalid/id" // invalid ID with slash

	var iresp resource.ImportStateResponse
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	iresp.State = tfsdk.State{Schema: schemaResp.Schema}

	r.ImportState(context.Background(), ireq, &iresp)
	require.True(t, iresp.Diagnostics.HasError())
}
