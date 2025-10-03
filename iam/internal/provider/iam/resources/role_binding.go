package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/client"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/validators"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RoleBindingResource{}
var _ resource.ResourceWithImportState = &RoleBindingResource{}

// RoleBindingResource defines the resource implementation for IAM role bindings
type RoleBindingResource struct {
	client     *client.Client
	iamService *iam.Service
}

// RoleBindingResourceModel describes the resource data model
type RoleBindingResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Role      types.String `tfsdk:"role"`
	Members   types.Set    `tfsdk:"members"`
	Condition types.String `tfsdk:"condition"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// NewRoleBindingResource creates a new role binding resource
func NewRoleBindingResource() resource.Resource {
	return &RoleBindingResource{}
}

// Metadata returns the resource type name
func (r *RoleBindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_binding"
}

// Schema defines the schema for the resource
func (r *RoleBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an IAM role binding within HiiRetail.",
		MarkdownDescription: "Manages an IAM role binding within HiiRetail. Role bindings grant roles to users and groups.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the role binding.",
				MarkdownDescription: "Unique identifier for the role binding.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the role binding. Must be unique within the tenant.",
				MarkdownDescription: "Name of the role binding. Must be unique within the tenant.",
				Required:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(1, 128),
					validators.IAMResourceName(),
					validators.StringNoLeadingTrailingSpaces(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role": schema.StringAttribute{
				Description:         "The role to be granted (e.g., 'roles/iam.viewer', 'roles/custom.developer').",
				MarkdownDescription: "The role to be granted (e.g., `roles/iam.viewer`, `roles/custom.developer`).",
				Required:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(1, 200),
				},
			},
			"members": schema.SetAttribute{
				ElementType:         types.StringType,
				Description:         "Set of member identifiers to grant the role to (user:email@domain.com or group:groupname).",
				MarkdownDescription: "Set of member identifiers to grant the role to in the format `user:email@domain.com` or `group:groupname`.",
				Required:            true,
				Validators:          []validator.Set{
					// TODO: Add set validators for member identifiers
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"condition": schema.StringAttribute{
				Description:         "Optional conditional expression to limit when the role binding applies.",
				MarkdownDescription: "Optional conditional expression to limit when the role binding applies.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(0, 1000),
				},
			},
			"created_at": schema.StringAttribute{
				Description:         "Timestamp when the role binding was created.",
				MarkdownDescription: "Timestamp when the role binding was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				Description:         "Timestamp when the role binding was last updated.",
				MarkdownDescription: "Timestamp when the role binding was last updated.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *RoleBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
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

	tflog.Info(ctx, "Configured IAM Role Binding Resource")
}

// Create creates the resource and sets the initial Terraform state
func (r *RoleBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API role binding object
	binding := &iam.RoleBinding{
		Name:      data.Name.ValueString(),
		Role:      data.Role.ValueString(),
		Condition: data.Condition.ValueString(),
	}

	// Convert members from Terraform set to string slice
	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		members := make([]string, 0, len(data.Members.Elements()))
		for _, elem := range data.Members.Elements() {
			if str, ok := elem.(types.String); ok {
				members = append(members, str.ValueString())
			}
		}
		binding.Members = members
	}

	// Create the role binding via API
	createdBinding, err := r.iamService.CreateRoleBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating IAM Role Binding",
			"Could not create role binding, unexpected error: "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.ID = types.StringValue(createdBinding.ID)
	// Preserve the configured name instead of using the API response
	// Note: We explicitly DO NOT set data.Name here to preserve the configured value
	data.Role = types.StringValue(createdBinding.Role)
	// Set condition properly - use null if empty to maintain Terraform consistency
	if createdBinding.Condition == "" {
		data.Condition = types.StringNull()
	} else {
		data.Condition = types.StringValue(createdBinding.Condition)
	}

	if len(createdBinding.Members) > 0 {
		memberElements := make([]attr.Value, len(createdBinding.Members))
		for i, member := range createdBinding.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(createdBinding.CreatedAt)
	data.UpdatedAt = types.StringValue(createdBinding.UpdatedAt)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created IAM role binding resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data
func (r *RoleBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RoleBindingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Printf("=== DEBUG Read METHOD START: ID=%s ===\n", data.ID.ValueString())
	fmt.Printf("CRITICAL: This debug output should appear if Read method is called\n")

	// Get role binding from API
	binding, err := r.iamService.GetRoleBinding(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Role binding no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading IAM Role Binding",
			"Could not read role binding "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response to resource model
	data.ID = types.StringValue(binding.ID)
	// NEVER set name - always preserve the configured name from Terraform
	fmt.Printf("DEBUG Read: NOT setting data.Name at all, preserving configured value: %s\n", data.Name.ValueString())
	// NOTE: data.Name is intentionally NOT set here to preserve the configuration value
	data.Role = types.StringValue(binding.Role)
	// Set condition properly - use null if empty to maintain Terraform consistency
	if binding.Condition == "" {
		data.Condition = types.StringNull()
	} else {
		data.Condition = types.StringValue(binding.Condition)
	}

	if len(binding.Members) > 0 {
		memberElements := make([]attr.Value, len(binding.Members))
		for i, member := range binding.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Handle timestamps - if API doesn't return them, preserve existing ones or use empty
	if binding.CreatedAt != "" {
		data.CreatedAt = types.StringValue(binding.CreatedAt)
	}
	// If API doesn't provide CreatedAt, keep the current value (don't override)
	
	if binding.UpdatedAt != "" {
		data.UpdatedAt = types.StringValue(binding.UpdatedAt)
	}
	// If API doesn't provide UpdatedAt, keep the current value (don't override)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *RoleBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RoleBindingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API role binding object
	binding := &iam.RoleBinding{
		Name:      data.Name.ValueString(),
		Role:      data.Role.ValueString(),
		Condition: data.Condition.ValueString(),
	}

	// Convert members from Terraform set to string slice
	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		members := make([]string, 0, len(data.Members.Elements()))
		for _, elem := range data.Members.Elements() {
			if str, ok := elem.(types.String); ok {
				members = append(members, str.ValueString())
			}
		}
		binding.Members = members
	}

	// Update the role binding via API
	updatedBinding, err := r.iamService.UpdateRoleBinding(ctx, data.ID.ValueString(), binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IAM Role Binding",
			"Could not update role binding "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.Role = types.StringValue(updatedBinding.Role)
	// Set condition properly - use null if empty to maintain Terraform consistency
	if updatedBinding.Condition == "" {
		data.Condition = types.StringNull()
	} else {
		data.Condition = types.StringValue(updatedBinding.Condition)
	}

	if len(updatedBinding.Members) > 0 {
		memberElements := make([]attr.Value, len(updatedBinding.Members))
		for i, member := range updatedBinding.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(updatedBinding.CreatedAt)
	data.UpdatedAt = types.StringValue(updatedBinding.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *RoleBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleBindingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the role binding via API
	err := r.iamService.DeleteRoleBinding(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Role binding already deleted, nothing to do
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting IAM Role Binding",
			"Could not delete role binding "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted IAM role binding resource")
}

// ImportState imports an existing resource into Terraform state
func (r *RoleBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the name as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
