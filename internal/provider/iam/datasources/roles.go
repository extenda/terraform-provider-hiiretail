package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RolesDataSource{}

// RolesDataSource defines the data source implementation for IAM roles
type RolesDataSource struct {
	client     *client.Client
	iamService *iam.Service
}

// RolesDataSourceModel describes the data source data model
type RolesDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Filter types.String `tfsdk:"filter"`
	Roles  types.List   `tfsdk:"roles"`
}

// RoleDataModel describes a single role in the data source
type RoleDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Stage       types.String `tfsdk:"stage"`
	Type        types.String `tfsdk:"type"`
}

// NewRolesDataSource creates a new roles data source
func NewRolesDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

// Metadata returns the data source type name
func (d *RolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_roles"
}

// Schema defines the schema for the data source
func (d *RolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a list of IAM roles within HiiRetail.",
		MarkdownDescription: "Retrieves a list of IAM roles within HiiRetail. Supports filtering to narrow down results by role type or other criteria.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the data source.",
				MarkdownDescription: "Unique identifier for the data source.",
				Computed:            true,
			},
			"filter": schema.StringAttribute{
				Description:         "Optional filter to narrow down the roles (e.g., 'type:custom' or 'name:iam.*').",
				MarkdownDescription: "Optional filter to narrow down the roles (e.g., `type:custom` or `name:iam.*`).",
				Optional:            true,
			},
			"roles": schema.ListNestedAttribute{
				Description:         "List of IAM roles matching the filter criteria.",
				MarkdownDescription: "List of IAM roles matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Unique identifier for the role.",
							MarkdownDescription: "Unique identifier for the role.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Name of the role.",
							MarkdownDescription: "Name of the role.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							Description:         "Human-readable title of the role.",
							MarkdownDescription: "Human-readable title of the role.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Description of the role.",
							MarkdownDescription: "Description of the role.",
							Computed:            true,
						},
						"stage": schema.StringAttribute{
							Description:         "Development stage of the role (ALPHA, BETA, GA).",
							MarkdownDescription: "Development stage of the role (ALPHA, BETA, GA).",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Type of the role (basic or custom).",
							MarkdownDescription: "Type of the role (`basic` or `custom`).",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *RolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
	d.iamService = iam.NewService(client, client.TenantID())

	tflog.Info(ctx, "Configured IAM Roles Data Source")
}

// Read refreshes the Terraform state with the latest data
func (d *RolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RolesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get filter value
	filter := ""
	if !config.Filter.IsNull() {
		filter = config.Filter.ValueString()
	}

	// Get roles from API
	roles, err := d.iamService.ListRoles(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read IAM Roles",
			err.Error(),
		)
		return
	}

	// Map API response to data source model
	config.ID = types.StringValue("roles")

	// Convert roles to list of objects
	roleElements := make([]attr.Value, len(roles))
	for i, role := range roles {
		roleObj := map[string]attr.Value{
			"id":          types.StringValue(role.ID),
			"name":        types.StringValue(role.Name),
			"title":       types.StringValue(role.Title),
			"description": types.StringValue(role.Description),
			"stage":       types.StringValue(role.Stage),
			"type":        types.StringValue(role.Type),
		}

		objType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":          types.StringType,
				"name":        types.StringType,
				"title":       types.StringType,
				"description": types.StringType,
				"stage":       types.StringType,
				"type":        types.StringType,
			},
		}

		objValue, diags := types.ObjectValue(objType.AttrTypes, roleObj)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		roleElements[i] = objValue
	}

	listType := types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":          types.StringType,
				"name":        types.StringType,
				"title":       types.StringType,
				"description": types.StringType,
				"stage":       types.StringType,
				"type":        types.StringType,
			},
		},
	}

	listValue, diags := types.ListValue(listType.ElemType, roleElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Roles = listValue

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

	tflog.Trace(ctx, "read roles data source")
}
