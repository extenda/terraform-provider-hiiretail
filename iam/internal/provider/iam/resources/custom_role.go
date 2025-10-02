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
var _ resource.Resource = &CustomRoleResource{}
var _ resource.ResourceWithImportState = &CustomRoleResource{}

// CustomRoleResource defines the resource implementation for IAM custom roles
type CustomRoleResource struct {
	client     *client.Client
	iamService *iam.Service
}

// CustomRoleResourceModel describes the resource data model
type CustomRoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Permissions types.Set    `tfsdk:"permissions"`
	Stage       types.String `tfsdk:"stage"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

// NewCustomRoleResource creates a new custom role resource
func NewCustomRoleResource() resource.Resource {
	return &CustomRoleResource{}
}

// Metadata returns the resource type name
func (r *CustomRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_custom_role"
}

// Schema defines the schema for the resource
func (r *CustomRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an IAM custom role within HiiRetail.",
		MarkdownDescription: "Manages an IAM custom role within HiiRetail. Custom roles allow you to define granular permissions for specific use cases.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the custom role.",
				MarkdownDescription: "Unique identifier for the custom role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the custom role. Must be unique within the tenant.",
				MarkdownDescription: "Name of the custom role. Must be unique within the tenant.",
				Required:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(1, 64),
					validators.IAMResourceName(),
					validators.StringNoLeadingTrailingSpaces(),
				},
			},
			"title": schema.StringAttribute{
				Description:         "Human-readable title for the custom role.",
				MarkdownDescription: "Human-readable title for the custom role.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(0, 100),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the custom role.",
				MarkdownDescription: "Description of the custom role.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(0, 500),
				},
			},
			"permissions": schema.SetAttribute{
				ElementType:         types.StringType,
				Description:         "Set of permissions for the custom role in format 'service.resource.action'.",
				MarkdownDescription: "Set of permissions for the custom role in format `service.resource.action` (e.g., `iam.groups.list`).",
				Required:            true,
				Validators:          []validator.Set{
					// TODO: Add set validators for permissions
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"stage": schema.StringAttribute{
				Description:         "Development stage of the custom role (ALPHA, BETA, GA).",
				MarkdownDescription: "Development stage of the custom role. Valid values are `ALPHA`, `BETA`, `GA`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					validators.StringOneOf("ALPHA", "BETA", "GA"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description:         "Timestamp when the custom role was created.",
				MarkdownDescription: "Timestamp when the custom role was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				Description:         "Timestamp when the custom role was last updated.",
				MarkdownDescription: "Timestamp when the custom role was last updated.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *CustomRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.iamService = iam.NewService(client)

	tflog.Info(ctx, "Configured IAM Custom Role Resource")
}

// Create creates the resource and sets the initial Terraform state
func (r *CustomRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CustomRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API custom role object
	role := &iam.CustomRole{
		Name:        data.Name.ValueString(),
		Title:       data.Title.ValueString(),
		Description: data.Description.ValueString(),
		Stage:       data.Stage.ValueString(),
	}

	// Convert permissions from Terraform set to string slice
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions := make([]string, 0, len(data.Permissions.Elements()))
		for _, elem := range data.Permissions.Elements() {
			if str, ok := elem.(types.String); ok {
				permissions = append(permissions, str.ValueString())
			}
		}
		role.Permissions = permissions
	}

	// Create the custom role via API
	createdRole, err := r.iamService.CreateCustomRole(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating IAM Custom Role",
			"Could not create custom role, unexpected error: "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.ID = types.StringValue(createdRole.ID)
	data.Name = types.StringValue(createdRole.Name)
	data.Title = types.StringValue(createdRole.Title)
	data.Description = types.StringValue(createdRole.Description)
	data.Stage = types.StringValue(createdRole.Stage)

	if len(createdRole.Permissions) > 0 {
		permissionElements := make([]attr.Value, len(createdRole.Permissions))
		for i, permission := range createdRole.Permissions {
			permissionElements[i] = types.StringValue(permission)
		}
		data.Permissions = types.SetValueMust(types.StringType, permissionElements)
	} else {
		data.Permissions = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(createdRole.CreatedAt)
	data.UpdatedAt = types.StringValue(createdRole.UpdatedAt)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created IAM custom role resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data
func (r *CustomRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CustomRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get custom role from API
	role, err := r.iamService.GetCustomRole(ctx, data.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Custom role no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading IAM Custom Role",
			"Could not read custom role "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response to resource model
	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)
	data.Title = types.StringValue(role.Title)
	data.Description = types.StringValue(role.Description)
	data.Stage = types.StringValue(role.Stage)

	if len(role.Permissions) > 0 {
		permissionElements := make([]attr.Value, len(role.Permissions))
		for i, permission := range role.Permissions {
			permissionElements[i] = types.StringValue(permission)
		}
		data.Permissions = types.SetValueMust(types.StringType, permissionElements)
	} else {
		data.Permissions = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(role.CreatedAt)
	data.UpdatedAt = types.StringValue(role.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *CustomRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CustomRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API custom role object
	role := &iam.CustomRole{
		Name:        data.Name.ValueString(),
		Title:       data.Title.ValueString(),
		Description: data.Description.ValueString(),
		Stage:       data.Stage.ValueString(),
	}

	// Convert permissions from Terraform set to string slice
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions := make([]string, 0, len(data.Permissions.Elements()))
		for _, elem := range data.Permissions.Elements() {
			if str, ok := elem.(types.String); ok {
				permissions = append(permissions, str.ValueString())
			}
		}
		role.Permissions = permissions
	}

	// Update the custom role via API
	updatedRole, err := r.iamService.UpdateCustomRole(ctx, data.Name.ValueString(), role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IAM Custom Role",
			"Could not update custom role "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.Title = types.StringValue(updatedRole.Title)
	data.Description = types.StringValue(updatedRole.Description)
	data.Stage = types.StringValue(updatedRole.Stage)

	if len(updatedRole.Permissions) > 0 {
		permissionElements := make([]attr.Value, len(updatedRole.Permissions))
		for i, permission := range updatedRole.Permissions {
			permissionElements[i] = types.StringValue(permission)
		}
		data.Permissions = types.SetValueMust(types.StringType, permissionElements)
	} else {
		data.Permissions = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.UpdatedAt = types.StringValue(updatedRole.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *CustomRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CustomRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the custom role via API
	err := r.iamService.DeleteCustomRole(ctx, data.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Custom role already deleted, nothing to do
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting IAM Custom Role",
			"Could not delete custom role "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted IAM custom role resource")
}

// ImportState imports an existing resource into Terraform state
func (r *CustomRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the name as the import identifier (custom roles are identified by name)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
