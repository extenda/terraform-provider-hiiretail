package resource_iam_role_binding

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SimpleIamRoleBindingResource{}
var _ resource.ResourceWithImportState = &SimpleIamRoleBindingResource{}

func NewSimpleIamRoleBindingResource() resource.Resource {
	return &SimpleIamRoleBindingResource{}
}

// SimpleIamRoleBindingResource defines the simple 1:1 role binding resource implementation.
type SimpleIamRoleBindingResource struct {
	client     *client.Client
	iamService *iam.Service
}

// Metadata returns the resource type name.
func (r *SimpleIamRoleBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_binding"
}

func (r *SimpleIamRoleBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SimpleIamRoleBindingResourceSchema(ctx)
}

func (r *SimpleIamRoleBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	r.iamService = iam.NewService(client, client.TenantID())
}

func (r *SimpleIamRoleBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SimpleRoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Creating simple IAM role binding resource")

	// Extract values from the model
	groupId := data.GroupID.ValueString()
	roleId := data.RoleID.ValueString()
	isCustom := data.IsCustom.ValueBool()

	// Extract bindings (optional)
	var bindings []string
	if !data.Bindings.IsNull() && !data.Bindings.IsUnknown() {
		elements := data.Bindings.Elements()
		for _, element := range elements {
			if bindingStr, ok := element.(types.String); ok {
				bindings = append(bindings, bindingStr.ValueString())
			}
		}
	}

	tflog.Debug(ctx, "Adding role to group", map[string]interface{}{
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
		"bindings":  bindings,
	})

	// Use the AddRoleToGroup method
	err := r.iamService.AddRoleToGroup(ctx, groupId, roleId, isCustom, bindings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Adding Role to Group",
			fmt.Sprintf("Could not add role %s to group %s, unexpected error: %s", roleId, groupId, err.Error()),
		)
		return
	}

	// Generate a unique ID for this role binding
	compositeId := GenerateResourceId(r.client.TenantID(), groupId, roleId)

	// Update the model with response data
	data.ID = types.StringValue(compositeId)
	data.TenantID = types.StringValue(r.client.TenantID())

	tflog.Trace(ctx, "Created simple IAM role binding resource", map[string]interface{}{
		"id":        compositeId,
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SimpleIamRoleBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SimpleRoleBindingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the composite ID to extract components
	id := data.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(
			"Error Reading Role Binding",
			"Role binding ID is empty",
		)
		return
	}

	// Parse the ID: format is "tenantId-groupId-roleId-hash"
	parts := strings.Split(id, "-")
	if len(parts) < 4 {
		resp.Diagnostics.AddError(
			"Error Reading Role Binding",
			fmt.Sprintf("Invalid role binding ID format: %s", id),
		)
		return
	}

	tenantId := parts[0]
	groupId := parts[1]
	roleId := parts[2]
	// hash is parts[3] (not needed for parsing)

	// Determine if the role is custom by checking if it exists as a custom role
	isCustom := false

	// Check if this role exists as a custom role
	_, err := r.iamService.GetCustomRole(ctx, roleId)
	if err != nil {
		// If GetCustomRole fails, it's likely a builtin role
		tflog.Debug(ctx, "Role not found as custom role, assuming builtin", map[string]interface{}{
			"role_id": roleId,
			"error":   err.Error(),
		})
		isCustom = false
	} else {
		// Role was found as a custom role
		isCustom = true
		tflog.Debug(ctx, "Role identified as custom", map[string]interface{}{
			"role_id": roleId,
		})
	}

	tflog.Debug(ctx, "Role type determination complete", map[string]interface{}{
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	// Update the model with parsed data
	data.TenantID = types.StringValue(tenantId)
	data.GroupID = types.StringValue(groupId)
	data.RoleID = types.StringValue(roleId)
	data.IsCustom = types.BoolValue(isCustom)

	// For bindings, we can't easily reconstruct them from just the ID
	// So we'll leave them as they were in the original state
	// In a production system, you might want to fetch them from the API

	tflog.Trace(ctx, "Read simple IAM role binding resource", map[string]interface{}{
		"id":        id,
		"tenant_id": tenantId,
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SimpleIamRoleBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SimpleRoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Since this is a simple 1:1 binding, updates are not really meaningful
	// The only things that could change are the bindings list or description
	// For now, we'll just log the update and save the state

	tflog.Trace(ctx, "Updated simple IAM role binding resource", map[string]interface{}{
		"id":        data.ID.ValueString(),
		"group_id":  data.GroupID.ValueString(),
		"role_id":   data.RoleID.ValueString(),
		"is_custom": data.IsCustom.ValueBool(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SimpleIamRoleBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SimpleRoleBindingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract values from the model
	groupId := data.GroupID.ValueString()
	roleId := data.RoleID.ValueString()
	isCustom := data.IsCustom.ValueBool()

	tflog.Debug(ctx, "Removing role from group", map[string]interface{}{
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	// Call the IAM service to remove the role from the group
	err := r.removeRoleFromGroup(ctx, groupId, roleId, isCustom)
	if err != nil {
		// Check if the error is a 404, which means the role binding doesn't exist
		// This is actually success since the desired state is that it doesn't exist
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not assigned to this group") {
			tflog.Debug(ctx, "Role binding already removed from group (404 response), treating as successful deletion", map[string]interface{}{
				"group_id":  groupId,
				"role_id":   roleId,
				"is_custom": isCustom,
				"error":     err.Error(),
			})
			return // Success - desired state achieved
		}

		resp.Diagnostics.AddError(
			"Error Deleting Role Binding",
			"Could not delete role binding, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted simple IAM role binding resource", map[string]interface{}{
		"id":        data.ID.ValueString(),
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})
}

// removeRoleFromGroup removes a role from a group using the V1 API
func (r *SimpleIamRoleBindingResource) removeRoleFromGroup(ctx context.Context, groupId, roleId string, isCustom bool) error {
	// Use the V1 API endpoint: DELETE /api/v1/tenants/{tenantId}/groups/{id}/roles/{roleId}
	// This endpoint properly supports role removal and has the required permissions
	path := fmt.Sprintf("/api/v1/tenants/%s/groups/%s/roles/%s", r.client.TenantID(), groupId, roleId)

	tflog.Debug(ctx, "Making DELETE request to remove role from group (V1 API)", map[string]interface{}{
		"path":      path,
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	// Build the request with isCustom query parameter if needed
	req := &client.Request{
		Method: "DELETE",
		Path:   path,
	}

	// Add isCustom query parameter if the role is custom
	if isCustom {
		req.Query = map[string]string{
			"isCustom": "true",
		}
	}

	resp, err := r.client.Do(ctx, req)
	if err != nil {
		tflog.Error(ctx, "HTTP request failed (V1 API)", map[string]interface{}{
			"error":    err.Error(),
			"path":     path,
			"group_id": groupId,
			"role_id":  roleId,
		})
		return fmt.Errorf("failed to remove role %s from group %s: %w", roleId, groupId, err)
	}

	tflog.Debug(ctx, "DELETE request completed (V1 API)", map[string]interface{}{
		"status_code": resp.StatusCode,
		"path":        path,
		"group_id":    groupId,
		"role_id":     roleId,
	})

	// Check if the request was successful
	if err := client.CheckResponse(resp); err != nil {
		tflog.Error(ctx, "DELETE request failed with error response (V1 API)", map[string]interface{}{
			"error":       err.Error(),
			"status_code": resp.StatusCode,
			"path":        path,
			"group_id":    groupId,
			"role_id":     roleId,
		})
		return fmt.Errorf("failed to remove role %s from group %s: %w", roleId, groupId, err)
	}

	tflog.Debug(ctx, "Successfully removed role from group (V1 API)", map[string]interface{}{
		"path":      path,
		"group_id":  groupId,
		"role_id":   roleId,
		"is_custom": isCustom,
	})

	return nil
}

func (r *SimpleIamRoleBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import state using the composite ID format: tenantId-groupId-roleId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
