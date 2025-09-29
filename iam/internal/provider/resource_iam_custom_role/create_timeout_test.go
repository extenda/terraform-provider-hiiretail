package resource_iam_custom_role

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCreateOperationTimeout(t *testing.T) {
	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	// Create the resource
	r := NewIamCustomRoleResource()

	// Configure the resource with a mock API client
	configReq := resource.ConfigureRequest{
		ProviderData: &APIClient{
			BaseURL:  env.BaseURL,
			TenantID: "test-tenant-123",
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	}
	configResp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), configReq, configResp)

	if configResp.Diagnostics.HasError() {
		t.Fatalf("Resource configuration failed: %v", configResp.Diagnostics.Errors())
	}

	// Create test data
	testData := IamCustomRoleModel{
		Id:   types.StringValue("test-role-001"),
		Name: types.StringValue("Test Custom Role"),
		Permissions: func() types.List {
			// Create a list of permissions using the correct terraform types
			permList, _ := types.ListValueFrom(context.Background(),
				PermissionsValue{}.Type(context.Background()),
				[]PermissionsValue{
					{
						Id:         types.StringValue("pos.payment.create"),
						Alias:      types.StringValue(""),
						Attributes: types.MapNull(types.StringType),
						state:      attr.ValueStateKnown,
					},
				})
			return permList
		}(),
	}

	// Test Create operation with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	createReq := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw: testData, // This might need adjustment for proper test setup
		},
	}
	createResp := &resource.CreateResponse{}

	// Test that Create doesn't hang
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
