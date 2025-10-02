package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GroupsDataSource{}

// GroupsDataSource defines the data source implementation for IAM groups
type GroupsDataSource struct {
	client     *client.Client
	iamService *iam.Service
}

// GroupsDataSourceModel describes the data source data model
type GroupsDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Filter types.String `tfsdk:"filter"`
	Groups types.List   `tfsdk:"groups"`
}

// GroupDataModel describes a single group in the data source
type GroupDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	MemberCount types.Int64  `tfsdk:"member_count"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

// NewGroupsDataSource creates a new groups data source
func NewGroupsDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

// Metadata returns the data source type name
func (d *GroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_groups"
}

// Schema defines the schema for the data source
func (d *GroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a list of IAM groups within HiiRetail.",
		MarkdownDescription: "Retrieves a list of IAM groups within HiiRetail. Supports filtering to narrow down results.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the data source.",
				MarkdownDescription: "Unique identifier for the data source.",
				Computed:            true,
			},
			"filter": schema.StringAttribute{
				Description:         "Optional filter to narrow down the groups (e.g., 'name:dev-*').",
				MarkdownDescription: "Optional filter to narrow down the groups (e.g., `name:dev-*`).",
				Optional:            true,
			},
			"groups": schema.ListNestedAttribute{
				Description:         "List of IAM groups matching the filter criteria.",
				MarkdownDescription: "List of IAM groups matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Unique identifier for the group.",
							MarkdownDescription: "Unique identifier for the group.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Name of the group.",
							MarkdownDescription: "Name of the group.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Description of the group.",
							MarkdownDescription: "Description of the group.",
							Computed:            true,
						},
						"member_count": schema.Int64Attribute{
							Description:         "Number of members in the group.",
							MarkdownDescription: "Number of members in the group.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							Description:         "Timestamp when the group was created.",
							MarkdownDescription: "Timestamp when the group was created.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *GroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	tflog.Info(ctx, "Configured IAM Groups Data Source")
}

// Read refreshes the Terraform state with the latest data
func (d *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GroupsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create list groups request
	listReq := &iam.ListGroupsRequest{}
	if !config.Filter.IsNull() {
		listReq.Filter = config.Filter.ValueString()
	}

	// Get groups from API
	listResp, err := d.iamService.ListGroups(ctx, listReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read IAM Groups",
			err.Error(),
		)
		return
	}

	// Map API response to data source model
	config.ID = types.StringValue("groups")

	// Convert groups to list of objects
	groupElements := make([]attr.Value, len(listResp.Groups))
	for i, group := range listResp.Groups {
		groupObj := map[string]attr.Value{
			"id":           types.StringValue(group.ID),
			"name":         types.StringValue(group.Name),
			"description":  types.StringValue(group.Description),
			"member_count": types.Int64Value(int64(len(group.Members))),
			"created_at":   types.StringValue(group.CreatedAt),
		}

		objType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":           types.StringType,
				"name":         types.StringType,
				"description":  types.StringType,
				"member_count": types.Int64Type,
				"created_at":   types.StringType,
			},
		}

		objValue, diags := types.ObjectValue(objType.AttrTypes, groupObj)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		groupElements[i] = objValue
	}

	listType := types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":           types.StringType,
				"name":         types.StringType,
				"description":  types.StringType,
				"member_count": types.Int64Type,
				"created_at":   types.StringType,
			},
		},
	}

	listValue, diags := types.ListValue(listType.ElemType, groupElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Groups = listValue

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

	tflog.Trace(ctx, "read groups data source")
}
