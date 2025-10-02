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
var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}

// GroupResource defines the resource implementation for IAM groups
type GroupResource struct {
	client     *client.Client
	iamService *iam.Service
}

// GroupResourceModel describes the resource data model
type GroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Members     types.Set    `tfsdk:"members"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

// NewGroupResource creates a new group resource
func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

// Metadata returns the resource type name
func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

// Schema defines the schema for the resource
func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an IAM group within HiiRetail.",
		MarkdownDescription: "Manages an IAM group within HiiRetail. Groups are collections of users that can be granted permissions through role bindings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the group.",
				MarkdownDescription: "Unique identifier for the group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the group. Must be unique within the tenant.",
				MarkdownDescription: "Name of the group. Must be unique within the tenant.",
				Required:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(1, 128),
					validators.IAMResourceName(),
					validators.StringNoLeadingTrailingSpaces(),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the group.",
				MarkdownDescription: "Description of the group.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringLengthBetween(0, 500),
				},
			},
			"members": schema.SetAttribute{
				ElementType:         types.StringType,
				Description:         "Set of member identifiers (user:email@domain.com or group:groupname).",
				MarkdownDescription: "Set of member identifiers in the format `user:email@domain.com` or `group:groupname`.",
				Optional:            true,
				Validators:          []validator.Set{
					// TODO: Add set validators for member identifiers
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description:         "Timestamp when the group was created.",
				MarkdownDescription: "Timestamp when the group was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				Description:         "Timestamp when the group was last updated.",
				MarkdownDescription: "Timestamp when the group was last updated.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	tflog.Info(ctx, "Configured IAM Group Resource")
}

// Create creates the resource and sets the initial Terraform state
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API group object
	group := &iam.Group{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Convert members from Terraform set to string slice
	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		members := make([]string, 0, len(data.Members.Elements()))
		for _, elem := range data.Members.Elements() {
			if str, ok := elem.(types.String); ok {
				members = append(members, str.ValueString())
			}
		}
		group.Members = members
	}

	// Create the group via API
	createdGroup, err := r.iamService.CreateGroup(ctx, group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating IAM Group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.ID = types.StringValue(createdGroup.ID)
	data.Name = types.StringValue(createdGroup.Name)
	data.Description = types.StringValue(createdGroup.Description)

	if len(createdGroup.Members) > 0 {
		memberElements := make([]attr.Value, len(createdGroup.Members))
		for i, member := range createdGroup.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(createdGroup.CreatedAt)
	data.UpdatedAt = types.StringValue(createdGroup.UpdatedAt)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created IAM group resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group from API
	group, err := r.iamService.GetGroup(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Group no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading IAM Group",
			"Could not read group ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response to resource model
	data.Name = types.StringValue(group.Name)
	data.Description = types.StringValue(group.Description)

	if len(group.Members) > 0 {
		memberElements := make([]attr.Value, len(group.Members))
		for i, member := range group.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.CreatedAt = types.StringValue(group.CreatedAt)
	data.UpdatedAt = types.StringValue(group.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API group object
	group := &iam.Group{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Convert members from Terraform set to string slice
	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		members := make([]string, 0, len(data.Members.Elements()))
		for _, elem := range data.Members.Elements() {
			if str, ok := elem.(types.String); ok {
				members = append(members, str.ValueString())
			}
		}
		group.Members = members
	}

	// Update the group via API
	updatedGroup, err := r.iamService.UpdateGroup(ctx, data.ID.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IAM Group",
			"Could not update group ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map API response back to resource model
	data.Name = types.StringValue(updatedGroup.Name)
	data.Description = types.StringValue(updatedGroup.Description)

	if len(updatedGroup.Members) > 0 {
		memberElements := make([]attr.Value, len(updatedGroup.Members))
		for i, member := range updatedGroup.Members {
			memberElements[i] = types.StringValue(member)
		}
		data.Members = types.SetValueMust(types.StringType, memberElements)
	} else {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.UpdatedAt = types.StringValue(updatedGroup.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the group via API
	err := r.iamService.DeleteGroup(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			// Group already deleted, nothing to do
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting IAM Group",
			"Could not delete group ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted IAM group resource")
}

// ImportState imports an existing resource into Terraform state
func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
