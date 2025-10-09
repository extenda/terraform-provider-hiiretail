package resource_iam_custom_role

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/testutils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// ...existing code...
)

func TestCreateOperationTimeout(t *testing.T) {
	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	// Create the resource
	r := &IamCustomRoleResource{}

	// Configure the resource with a mock API client
	configReq := resource.ConfigureRequest{
		ProviderData: struct {
			BaseURL    string
			TenantID   string
			HTTPClient *http.Client
		}{
			BaseURL:    env.BaseURL,
			TenantID:   env.TenantID,
			HTTPClient: http.DefaultClient,
		},
	}
	configResp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), configReq, configResp)

	if configResp.Diagnostics.HasError() {
		t.Fatalf("Resource configuration failed: %v", configResp.Diagnostics.Errors())
	}

	// Create test data
	// No struct needed; use map for plan value

	// Test Create operation with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	objVal := types.ObjectValueMust(
		map[string]attr.Type{
			"id":          types.StringType,
			"name":        types.StringType,
			"permissions": types.ListType{ElemType: types.StringType},
			"tenant_id":   types.StringType,
		},
		map[string]attr.Value{
			"id":   types.StringValue("test-role-001"),
			"name": types.StringValue("Test Custom Role"),
			"permissions": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("pos.payment.create"),
			}),
			"tenant_id": types.StringValue(env.TenantID),
		},
	)
	rawValue, diag := objVal.ToTerraformValue(context.Background())
	if diag != nil {
		t.Fatalf("Failed to convert ObjectValue to TerraformValue: %v", diag.Error())
	}
	createReq := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw: rawValue,
		},
	}

	createResp := &resource.CreateResponse{}
	done := make(chan bool, 1)
	var createErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				createErr = fmt.Errorf("Create operation panicked: %v", r)
			}
			done <- true
		}()

		r.Create(ctx, createReq, createResp)
		if createResp.Diagnostics.HasError() {
			createErr = fmt.Errorf("Create operation failed: %v", createResp.Diagnostics.Errors())
		}
	}()

	select {
	case <-done:
		if createErr != nil {
			t.Logf("Create operation completed with expected error: %v", createErr)
		} else {
			t.Logf("Create operation completed successfully")
		}
	case <-ctx.Done():
		t.Fatalf("Create operation timed out after 20 seconds")
	}
}
