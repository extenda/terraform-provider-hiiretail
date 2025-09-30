package resource_iam_role_binding

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// APIClient represents the configuration for making API calls
// This matches the APIClient from the provider package
type APIClient struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IamRoleBindingResource{}
var _ resource.ResourceWithImportState = &IamRoleBindingResource{}

func NewIamRoleBindingResource() resource.Resource {
	return &IamRoleBindingResource{}
}

// IamRoleBindingResource defines the resource implementation.
type IamRoleBindingResource struct {
	client *APIClient
}

func (r *IamRoleBindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_binding"
}

func (r *IamRoleBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IamRoleBindingResourceSchema(ctx)
}

func (r *IamRoleBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IamRoleBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IamRoleBindingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Creating IAM role binding resource")

	// Extract bindings from types.List to []string
	var bindings []string
	resp.Diagnostics.Append(data.Bindings.ElementsAs(ctx, &bindings, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the role binding model
	err := ValidateRoleBindingModel(ctx, data.RoleId.ValueString(), data.IsCustom.ValueBool(), bindings)
	if err != nil {
		resp.Diagnostics.AddError("Role Binding Validation Failed", err.Error())
		return
	}

	// Validate maximum bindings
	err = ValidateMaxBindings(bindings)
	if err != nil {
		resp.Diagnostics.AddError("Maximum Bindings Validation Failed", err.Error())
		return
	}

	// Validate binding formats
	err = ValidateBindingFormat(bindings)
	if err != nil {
		resp.Diagnostics.AddError("Binding Format Validation Failed", err.Error())
		return
	}

	// Create the role binding via API
	roleBinding, err := r.createRoleBinding(ctx, data.RoleId.ValueString(), data.IsCustom.ValueBool(), bindings)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create role binding, got error: %s", err))
		return
	}

	// Update the model with response data
	data.Id = types.StringValue(roleBinding.ID)
	data.TenantId = types.StringValue(roleBinding.TenantId)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Created IAM role binding resource", map[string]interface{}{
		"id":      roleBinding.ID,
		"role_id": roleBinding.RoleId,
	})
}

func (r *IamRoleBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IamRoleBindingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Reading IAM role binding resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Get role binding from API
	roleBinding, err := r.readRoleBinding(ctx, data.Id.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			// Role binding was deleted outside of Terraform
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read role binding, got error: %s", err))
		return
	}

	// Update model with API response
	data.RoleId = types.StringValue(roleBinding.RoleId)
	data.IsCustom = types.BoolValue(roleBinding.IsCustom)
	data.TenantId = types.StringValue(roleBinding.TenantId)

	// Convert bindings to types.List
	bindingsElements := make([]types.String, len(roleBinding.Bindings))
	for i, binding := range roleBinding.Bindings {
		bindingsElements[i] = types.StringValue(binding)
	}
	bindingsList, diags := types.ListValueFrom(ctx, types.StringType, bindingsElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Bindings = bindingsList

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IamRoleBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IamRoleBindingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updating IAM role binding resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Extract bindings from types.List to []string
	var bindings []string
	resp.Diagnostics.Append(data.Bindings.ElementsAs(ctx, &bindings, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the updated role binding model
	err := ValidateRoleBindingModel(ctx, data.RoleId.ValueString(), data.IsCustom.ValueBool(), bindings)
	if err != nil {
		resp.Diagnostics.AddError("Role Binding Validation Failed", err.Error())
		return
	}

	// Validate maximum bindings
	err = ValidateMaxBindings(bindings)
	if err != nil {
		resp.Diagnostics.AddError("Maximum Bindings Validation Failed", err.Error())
		return
	}

	// Validate binding formats
	err = ValidateBindingFormat(bindings)
	if err != nil {
		resp.Diagnostics.AddError("Binding Format Validation Failed", err.Error())
		return
	}

	// Update the role binding via API
	roleBinding, err := r.updateRoleBinding(ctx, data.Id.ValueString(), data.RoleId.ValueString(), data.IsCustom.ValueBool(), bindings)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update role binding, got error: %s", err))
		return
	}

	// Update the model with response data
	data.TenantId = types.StringValue(roleBinding.TenantId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Updated IAM role binding resource", map[string]interface{}{
		"id":      roleBinding.ID,
		"role_id": roleBinding.RoleId,
	})
}

func (r *IamRoleBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IamRoleBindingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Deleting IAM role binding resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Delete role binding via API
	err := r.deleteRoleBinding(ctx, data.Id.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			// Role binding was already deleted outside of Terraform
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete role binding, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Deleted IAM role binding resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})
}

func (r *IamRoleBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// API interaction methods

type RoleBindingResponse struct {
	ID       string   `json:"id"`
	RoleId   string   `json:"role_id"`
	IsCustom bool     `json:"is_custom"`
	Bindings []string `json:"bindings"`
	TenantId string   `json:"tenant_id"`
}

func (r *IamRoleBindingResource) createRoleBinding(ctx context.Context, roleId string, isCustom bool, bindings []string) (*RoleBindingResponse, error) {
	// TODO: Implement actual API call to create role binding
	// This is a placeholder that will be implemented in Phase 3.3
	return &RoleBindingResponse{
		ID:       "rb-" + generateUUID(),
		RoleId:   roleId,
		IsCustom: isCustom,
		Bindings: bindings,
		TenantId: r.client.TenantID,
	}, nil
}

func (r *IamRoleBindingResource) readRoleBinding(ctx context.Context, id string) (*RoleBindingResponse, error) {
	// TODO: Implement actual API call to read role binding
	// This is a placeholder that will be implemented in Phase 3.3
	return &RoleBindingResponse{
		ID:       id,
		RoleId:   "placeholder-role",
		IsCustom: true,
		Bindings: []string{"user:placeholder"},
		TenantId: r.client.TenantID,
	}, nil
}

func (r *IamRoleBindingResource) updateRoleBinding(ctx context.Context, id, roleId string, isCustom bool, bindings []string) (*RoleBindingResponse, error) {
	// TODO: Implement actual API call to update role binding
	// This is a placeholder that will be implemented in Phase 3.3
	return &RoleBindingResponse{
		ID:       id,
		RoleId:   roleId,
		IsCustom: isCustom,
		Bindings: bindings,
		TenantId: r.client.TenantID,
	}, nil
}

func (r *IamRoleBindingResource) deleteRoleBinding(ctx context.Context, id string) error {
	// TODO: Implement actual API call to delete role binding
	// This is a placeholder that will be implemented in Phase 3.3
	return nil
}

// Helper functions

func isNotFoundError(err error) bool {
	// TODO: Implement proper error detection
	// This should check if the error indicates a 404 Not Found response
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404")
}

func generateUUID() string {
	// TODO: Implement proper UUID generation
	// This is a placeholder for demonstration
	return "550e8400-e29b-41d4-a716-446655440000"
}
