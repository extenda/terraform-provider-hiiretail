package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ResourceDataSource{}

// ResourceDataSource defines the data source implementation for IAM resource
type ResourceDataSource struct {
	client     *client.Client
	iamService *iam.Service
}

// ResourceDataSourceModel describes the data source data model
type ResourceDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Properties types.String `tfsdk:"properties"`
}

// NewResourceDataSource creates a new resource data source
func NewResourceDataSource() datasource.DataSource {
	return &ResourceDataSource{}
}

// Metadata returns the data source type name
func (d *ResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_resource"
}

// Schema defines the schema for the data source
func (d *ResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches an IAM resource by ID from HiiRetail.",
		MarkdownDescription: "Fetches an IAM resource by ID from HiiRetail. This is useful for referencing resources that are managed outside of Terraform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the resource in the format 'type:name'.",
				MarkdownDescription: "The unique identifier of the resource in the format `type:name`.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the resource.",
				MarkdownDescription: "The name of the resource.",
				Computed:            true,
			},
			"properties": schema.StringAttribute{
				Description:         "JSON string containing additional properties of the resource.",
				MarkdownDescription: "JSON string containing additional properties of the resource.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *ResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected map[string]interface{}, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	client, ok := clientData["client"].(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Client Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", clientData["client"]),
		)
		return
	}

	iamService, ok := clientData["iam_service"].(*iam.Service)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected IAM Service Type",
			fmt.Sprintf("Expected *iam.Service, got: %T. Please report this issue to the provider developers.", clientData["iam_service"]),
		)
		return
	}

	d.client = client
	d.iamService = iamService
}

// Read refreshes the Terraform state with the latest data
func (d *ResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ResourceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource ID
	resourceID := config.ID.ValueString()

	tflog.Debug(ctx, "Fetching IAM resource", map[string]interface{}{
		"resource_id": resourceID,
	})

	// Fetch the resource from the API
	resource, err := d.iamService.GetResource(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Fetch Resource",
			fmt.Sprintf("Could not fetch resource %s: %s", resourceID, err.Error()),
		)
		return
	}

	// Map response to model
	config.ID = types.StringValue(resource.ID)
	config.Name = types.StringValue(resource.Name)

	// Set properties as JSON string
	if resource.Props != nil {
		propsJSON, err := json.Marshal(resource.Props)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Unable to serialize properties",
				fmt.Sprintf("Could not serialize properties to JSON: %s", err.Error()),
			)
			config.Properties = types.StringValue("{}")
		} else {
			config.Properties = types.StringValue(string(propsJSON))
		}
	} else {
		config.Properties = types.StringValue("{}")
	}

	tflog.Debug(ctx, "Successfully fetched IAM resource", map[string]interface{}{
		"resource_id": resourceID,
		"name":        resource.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
