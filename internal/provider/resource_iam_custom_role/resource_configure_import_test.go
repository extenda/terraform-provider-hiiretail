package resource_iam_custom_role

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConfigure_WithAPIClientPointer(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	client := &APIClient{BaseURL: "http://api", TenantID: "tid", HTTPClient: &http.Client{}}

	var req resource.ConfigureRequest
	var resp resource.ConfigureResponse
	req.ProviderData = client

	r.Configure(context.Background(), req, &resp)

	if r.baseURL != "http://api" {
		t.Fatalf("expected baseURL set, got %s", r.baseURL)
	}
	if r.tenantID != "tid" {
		t.Fatalf("expected tenantID set, got %s", r.tenantID)
	}
	if r.client == nil {
		t.Fatalf("expected HTTP client set")
	}
}

func TestConfigure_WithReflectedStruct(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Create a struct type that matches provider.APIClient fields
	providerData := struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}{
		BaseURL:    "http://reflected",
		TenantID:   "ref-tid",
		HTTPClient: &http.Client{},
	}

	var req resource.ConfigureRequest
	var resp resource.ConfigureResponse
	req.ProviderData = &providerData

	r.Configure(context.Background(), req, &resp)

	if r.baseURL != "http://reflected" {
		t.Fatalf("expected baseURL set from reflected struct, got %s", r.baseURL)
	}
	if r.tenantID != "ref-tid" {
		t.Fatalf("expected tenantID set from reflected struct, got %s", r.tenantID)
	}
	if r.client == nil {
		t.Fatalf("expected HTTP client set from reflected struct")
	}
}

func TestConfigure_WithInvalidProviderData(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var req resource.ConfigureRequest
	var resp resource.ConfigureResponse
	req.ProviderData = "not-a-struct"

	r.Configure(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics error when provider data is invalid")
	}
}

func TestImportState_MissingID(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var req resource.ImportStateRequest
	var resp resource.ImportStateResponse

	req.ID = ""

	r.ImportState(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics when import ID is missing")
	}
}
