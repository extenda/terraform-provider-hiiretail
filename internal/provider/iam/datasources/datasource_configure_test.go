package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestGroupsDataSource_Configure_NoPanicOnNilProviderData(t *testing.T) {
	d := &GroupsDataSource{}
	var resp datasource.ConfigureResponse
	// Should not panic or populate diagnostics when ProviderData is nil
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics error for nil ProviderData, got: %v", resp.Diagnostics)
	}
}

func TestGroupsDataSource_Configure_InvalidTypeProducesDiagnostic(t *testing.T) {
	d := &GroupsDataSource{}
	var resp datasource.ConfigureResponse
	// Pass an unexpected type for ProviderData and expect a diagnostic error
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "not-a-client"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics error for invalid ProviderData type, got none")
	}
}

func TestRolesDataSource_Configure_InvalidTypeProducesDiagnostic(t *testing.T) {
	d := &RolesDataSource{}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: 12345}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics error for invalid ProviderData type on roles datasource, got none")
	}
}

func TestGroupsDataSource_MetadataAndSchema(t *testing.T) {
	d := &GroupsDataSource{}
	var metaReq datasource.MetadataRequest
	metaReq.ProviderTypeName = "testprov"
	var metaResp datasource.MetadataResponse
	d.Metadata(context.Background(), metaReq, &metaResp)
	if metaResp.TypeName != "testprov_iam_groups" {
		t.Fatalf("unexpected type name: %s", metaResp.TypeName)
	}

	var schemaReq datasource.SchemaRequest
	var schemaResp datasource.SchemaResponse
	d.Schema(context.Background(), schemaReq, &schemaResp)
	if schemaResp.Schema.Attributes == nil {
		t.Fatalf("expected schema attributes to be set")
	}
	// quick checks for expected attribute names
	if _, ok := schemaResp.Schema.Attributes["groups"]; !ok {
		t.Fatalf("expected 'groups' attribute in schema")
	}
	if _, ok := schemaResp.Schema.Attributes["filter"]; !ok {
		t.Fatalf("expected 'filter' attribute in schema")
	}
}

func TestRolesDataSource_MetadataAndSchema(t *testing.T) {
	d := &RolesDataSource{}
	var metaReq datasource.MetadataRequest
	metaReq.ProviderTypeName = "testprov"
	var metaResp datasource.MetadataResponse
	d.Metadata(context.Background(), metaReq, &metaResp)
	if metaResp.TypeName != "testprov_iam_roles" {
		t.Fatalf("unexpected type name: %s", metaResp.TypeName)
	}

	var schemaReq datasource.SchemaRequest
	var schemaResp datasource.SchemaResponse
	d.Schema(context.Background(), schemaReq, &schemaResp)
	if schemaResp.Schema.Attributes == nil {
		t.Fatalf("expected schema attributes to be set")
	}
	if _, ok := schemaResp.Schema.Attributes["roles"]; !ok {
		t.Fatalf("expected 'roles' attribute in schema")
	}
}
