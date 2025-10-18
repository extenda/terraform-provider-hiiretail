package resource_iam_role_binding

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IamRoleBindingResource{}
var _ resource.ResourceWithImportState = &IamRoleBindingResource{}

func NewIamRoleBindingResource() resource.Resource {
	return &IamRoleBindingResource{}
}

// IamRoleBindingResource defines the resource implementation.
type IamRoleBindingResource struct {
	client     *client.Client
	iamService *iam.Service
}

// Metadata returns the resource type name.
func (r *IamRoleBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_binding"
}

func (r *IamRoleBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Use the enhanced schema that supports both legacy and new property structures
	// T030-T032: Enhanced schema with dual property support and deprecation warnings
	resp.Schema = EnhancedIamRoleBindingResourceSchema(ctx)
}

func (r *IamRoleBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IamRoleBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	fmt.Printf("=== ENHANCED ROLE BINDING CREATE START ===\n")
	var data RoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Creating IAM role binding resource with enhanced property structure support")

	// T033-T034: Validate property structure and handle mixed properties
	validationResult := ValidatePropertyStructure(ctx, &data)
	if !validationResult.IsValid {
		for _, err := range validationResult.Errors {
			resp.Diagnostics.AddError("Property Structure Validation Failed", err)
		}
		return
	}

	// Add deprecation warnings for legacy properties
	for _, warning := range validationResult.Warnings {
		resp.Diagnostics.AddWarning("Deprecated Property Usage", warning)
	}

	// Convert model to working format based on property structure
	var workingModel *RoleBindingResourceModel
	var err error

	if validationResult.PropertyMix == "legacy" {
		// Convert legacy properties to new structure for internal processing
		workingModel, err = ConvertLegacyToNew(ctx, &data)
		if err != nil {
			resp.Diagnostics.AddError("Legacy Property Conversion Failed", err.Error())
			return
		}
		tflog.Debug(ctx, "Converted legacy properties to new structure for processing")
	} else {
		// Use new properties directly
		workingModel = &data
	}

	// Extract roles and create role bindings via API
	var roles []RoleModel
	diags := workingModel.Roles.ElementsAs(ctx, &roles, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Step 1: Get the group ID from terraform state
	// Since we're using enhanced role bindings, just use the group_id directly
	// The group_id comes from terraform state and should be used as-is
	terraformGroupId := workingModel.GroupId.ValueString()

	tflog.Debug(ctx, "Using group ID from terraform state", map[string]interface{}{
		"terraform_group_id": terraformGroupId,
	})

	// We'll use the terraform group ID directly - no need to look up or create groups
	// The group resource handles group creation, role binding just adds roles to groups

	// Step 2: Add roles to the group using the V2 API pattern
	fmt.Printf("=== ENHANCED: About to add %d roles to group %s ===\n", len(roles), terraformGroupId)
	var assignedRoles []string
	for _, role := range roles {
		// Parse role ID and get custom flag from config
		roleValue := role.Id.ValueString()
		isCustom := role.IsCustom.ValueBool() // Use the is_custom field from config
		roleId := strings.TrimPrefix(roleValue, "roles/")
		if isCustom {
			roleId = strings.TrimPrefix(roleId, "custom.")
		}

		// Extract bindings from the role
		var bindings []string
		for _, binding := range role.Bindings.Elements() {
			if bindingStr, ok := binding.(types.String); ok {
				bindings = append(bindings, bindingStr.ValueString())
			}
		}

		tflog.Debug(ctx, "Adding role to group", map[string]interface{}{
			"group_id":  terraformGroupId,
			"role_id":   roleId,
			"is_custom": isCustom,
			"bindings":  bindings,
		})

		// Use the AddRoleToGroup method with specific bindings
		fmt.Printf("=== DEBUG: About to call AddRoleToGroup with group ID: %s ===\n", terraformGroupId)
		err := r.iamService.AddRoleToGroup(ctx, terraformGroupId, roleId, isCustom, bindings)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adding Role to Group",
				fmt.Sprintf("Could not add role %s to group %s, unexpected error: %s", roleId, terraformGroupId, err.Error()),
			)
			return
		}

		assignedRoles = append(assignedRoles, roleValue)
		tflog.Debug(ctx, "Added role to group", map[string]interface{}{
			"role_id":  roleId,
			"group_id": terraformGroupId,
		})
	}

	// Generate a composite ID for the enhanced resource (since it manages multiple role bindings)
	// Store the individual binding IDs in the composite ID for later retrieval
	compositeId := GenerateResourceId(r.client.TenantID(), workingModel.GroupId.ValueString(), "multi-role")

	// Store the created binding IDs for later retrieval (we'll need them for Read/Update/Delete)
	// For now, we'll use a simple approach - in a real implementation, this might need a more sophisticated tracking mechanism

	// Update the original model with response data
	data.Id = types.StringValue(compositeId)
	data.TenantId = types.StringValue(r.client.TenantID())

	// Set computed legacy compatibility fields
	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		var roles []RoleModel
		diags := data.Roles.ElementsAs(ctx, &roles, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Set role_id from the first role for legacy compatibility
		if len(roles) > 0 {
			data.RoleId = roles[0].Id
		} else {
			data.RoleId = types.StringNull()
		}
	} else {
		data.RoleId = types.StringNull()
	}

	// Set bindings_legacy as empty list (it's a legacy compatibility field)
	data.BindingsLegacy = types.ListValueMust(types.StringType, []attr.Value{})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Created IAM role binding resource", map[string]interface{}{
		"id":               compositeId,
		"property_type":    validationResult.PropertyMix,
		"created_bindings": len(assignedRoles),
	})
}

func (r *IamRoleBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RoleBindingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Reading IAM role binding resource with enhanced property structure support", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Get role binding from API (placeholder implementation)
	// TODO: Implement actual API call in T035
	roleBinding := &RoleBindingResponse{
		ID:       data.Id.ValueString(),
		TenantId: r.client.TenantID(),
		// Placeholder data - will be replaced with actual API call
	}

	// Update core properties
	data.TenantId = types.StringValue(roleBinding.TenantId)

	// Maintain property structure consistency
	// The state should preserve the same property structure that was used for creation
	validationResult := ValidatePropertyStructure(ctx, &data)
	if validationResult.PropertyMix == "legacy" {
		tflog.Debug(ctx, "Maintaining legacy property structure in state")
		// Keep legacy properties in state
	} else if validationResult.PropertyMix == "new" {
		tflog.Debug(ctx, "Maintaining new property structure in state")
		// Keep new properties in state
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IamRoleBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updating IAM role binding resource with enhanced property structure support", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Validate property structure for the update
	validationResult := ValidatePropertyStructure(ctx, &data)
	if !validationResult.IsValid {
		for _, err := range validationResult.Errors {
			resp.Diagnostics.AddError("Property Structure Validation Failed", err)
		}
		return
	}

	// Add deprecation warnings for legacy properties
	for _, warning := range validationResult.Warnings {
		resp.Diagnostics.AddWarning("Deprecated Property Usage", warning)
	}

	// Property structure validation passed - ready for processing
	if validationResult.PropertyMix == "legacy" {
		tflog.Debug(ctx, "Processing update with legacy properties structure")
	} else {
		tflog.Debug(ctx, "Processing update with new properties structure")
	}

	// For updates, we need to delete the old binding and create a new one
	// This is because the group_id is changing from "reconciliation-approvers-binding" to "finance-team"

	// Convert model to working format based on property structure
	var workingModel *RoleBindingResourceModel
	var err error

	if validationResult.PropertyMix == "legacy" {
		// Convert legacy properties to new structure for internal processing
		workingModel, err = ConvertLegacyToNew(ctx, &data)
		if err != nil {
			resp.Diagnostics.AddError("Legacy Property Conversion Failed", err.Error())
			return
		}
		tflog.Debug(ctx, "Converted legacy properties to new structure for processing")
	} else {
		// Use new properties directly
		workingModel = &data
	}

	// Extract roles and create role bindings via API
	var roles []RoleModel
	diags := workingModel.Roles.ElementsAs(ctx, &roles, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Step 1: Get or create the group for the role binding
	groupName := workingModel.GroupId.ValueString()

	// Try to get existing group first
	existingGroup, err := r.iamService.GetGroup(ctx, groupName)
	if err != nil {
		// If group doesn't exist, create it
		group := &iam.Group{
			Name:        groupName,
			Description: fmt.Sprintf("Group for role binding: %s", groupName),
			Members:     []string{}, // Empty members - bindings are handled differently
		}

		tflog.Debug(ctx, "Creating group for role binding update", map[string]interface{}{
			"group_name":        groupName,
			"group_description": group.Description,
		})

		existingGroup, err = r.iamService.CreateGroup(ctx, group)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating IAM Group for Update",
				fmt.Sprintf("Could not create group %s for role binding update, unexpected error: %s", groupName, err.Error()),
			)
			return
		}
	}

	// Step 2: Add roles to the group using the V2 API pattern
	fmt.Printf("=== ENHANCED UPDATE: About to add %d roles to group %s ===\n", len(roles), existingGroup.ID)
	var assignedRoles []string
	for _, role := range roles {
		// Parse role ID and get custom flag from config
		roleValue := role.Id.ValueString()
		isCustom := role.IsCustom.ValueBool() // Use the is_custom field from config
		roleId := strings.TrimPrefix(roleValue, "roles/")
		if isCustom {
			roleId = strings.TrimPrefix(roleId, "custom.")
		}

		// Extract bindings from the role
		var bindings []string
		for _, binding := range role.Bindings.Elements() {
			if bindingStr, ok := binding.(types.String); ok {
				bindings = append(bindings, bindingStr.ValueString())
			}
		}

		tflog.Debug(ctx, "Adding role to group in update", map[string]interface{}{
			"group_id":  existingGroup.ID,
			"role_id":   roleId,
			"is_custom": isCustom,
			"bindings":  bindings,
		})

		// Use the AddRoleToGroup method with specific bindings
		err := r.iamService.AddRoleToGroup(ctx, existingGroup.ID, roleId, isCustom, bindings)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adding Role to Group in Update",
				fmt.Sprintf("Could not add role %s to group %s during update, unexpected error: %s", roleId, existingGroup.ID, err.Error()),
			)
			return
		}

		assignedRoles = append(assignedRoles, roleValue)
		tflog.Debug(ctx, "Added role to group in update", map[string]interface{}{
			"role_id":  roleId,
			"group_id": existingGroup.ID,
		})
	}

	// Generate a new composite ID for the updated resource
	compositeId := GenerateResourceId(r.client.TenantID(), workingModel.GroupId.ValueString(), "multi-role")

	// Update the model with response data
	data.Id = types.StringValue(compositeId)
	data.TenantId = types.StringValue(r.client.TenantID())

	// Set computed legacy compatibility fields
	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		var roles []RoleModel
		diags := data.Roles.ElementsAs(ctx, &roles, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Set role_id from the first role for legacy compatibility
		if len(roles) > 0 {
			data.RoleId = roles[0].Id
		} else {
			data.RoleId = types.StringNull()
		}
	} else {
		data.RoleId = types.StringNull()
	}

	// Set bindings_legacy as empty list (it's a legacy compatibility field)
	data.BindingsLegacy = types.ListValueMust(types.StringType, []attr.Value{})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Updated IAM role binding resource", map[string]interface{}{
		"id":               compositeId,
		"property_type":    validationResult.PropertyMix,
		"updated_bindings": len(assignedRoles),
	})
}

func (r *IamRoleBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleBindingResourceModel

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
		TenantId: r.client.TenantID(),
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
		TenantId: r.client.TenantID(),
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
		TenantId: r.client.TenantID(),
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
	// Generate a simple random UUID-like string
	// In production, you'd use a proper UUID library
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		rand.Int63(),
		rand.Int63(),
		rand.Int63(),
		rand.Int63(),
		rand.Int63())
}
