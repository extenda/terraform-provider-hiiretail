package resource_iam_resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IAMResourceResource{}
var _ resource.ResourceWithImportState = &IAMResourceResource{}

func NewIAMResourceResource() resource.Resource {
	return &IAMResourceResource{}
}

// IAMResourceResource defines the resource implementation.
type IAMResourceResource struct {
	service *iam.Service
}

// IAMResourceResourceModel describes the resource data model.
type IAMResourceResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Props    types.String `tfsdk:"props"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func (r *IAMResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_resource"
}

func (r *IAMResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "HiiRetail IAM Resource for granular access control management.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the resource within the tenant. Must match pattern `^(?!\\.\\..?$)(?!.*__.*__)([^/]{1,1500})$`",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1500),
					// Custom validation for resource ID format rules
					// Note: Using our custom validateResourceID function in runtime validation
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable display name for the resource",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"props": schema.StringAttribute{
				MarkdownDescription: "Flexible properties object as JSON string that can contain additional metadata",
				Optional:            true,
				Validators: []validator.String{
					&jsonValidator{},
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "Tenant identifier, inherited from provider configuration",
				Computed:            true,
			},
		},
	}
}

func (r *IAMResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.service = iam.NewService(client, client.TenantID())
}

func (r *IAMResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IAMResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate ID if not provided
	resourceId := data.ID.ValueString()
	if resourceId == "" {
		// Generate a unique ID based on name (simple implementation)
		resourceId = data.Name.ValueString()
	}

	// Validate props JSON if provided
	var propsData interface{}
	if !data.Props.IsNull() && !data.Props.IsUnknown() {
		propsStr := data.Props.ValueString()
		if err := ValidateJSONString(propsStr); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Props",
				fmt.Sprintf("Props must be valid JSON: %s", err.Error()),
			)
			return
		}

		// Parse JSON for API call
		if propsStr != "" {
			if err := json.Unmarshal([]byte(propsStr), &propsData); err != nil {
				resp.Diagnostics.AddError(
					"Props JSON Parse Error",
					fmt.Sprintf("Failed to parse props JSON: %s", err.Error()),
				)
				return
			}
		}
	}

	// Create API request
	createRequest := &iam.SetResourceDto{
		Name:  data.Name.ValueString(),
		Props: propsData,
	}

	// Call SetResource API (PUT endpoint)
	createdResource, err := r.service.SetResource(ctx, resourceId, createRequest)
	if err != nil {
		title, detail := handleAPIError(err, "create", resourceId)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	// Map response to model
	data.ID = types.StringValue(createdResource.ID)
	data.Name = types.StringValue(createdResource.Name)
	data.TenantID = types.StringValue(r.service.TenantID())

	// Handle props response
	if createdResource.Props != nil {
		propsJson, err := json.Marshal(createdResource.Props)
		if err != nil {
			resp.Diagnostics.AddError(
				"Props Serialization Error",
				fmt.Sprintf("Failed to serialize props: %s", err.Error()),
			)
			return
		}
		data.Props = types.StringValue(string(propsJson))
	} else {
		data.Props = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IAMResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IAMResourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource ID from state
	resourceId := data.ID.ValueString()
	if resourceId == "" {
		resp.Diagnostics.AddError(
			"Missing Resource ID",
			"Resource ID is required for read operation",
		)
		return
	}

	// Call GetResource API
	resource, err := r.service.GetResource(ctx, resourceId)
	if err != nil {
		// Handle 404 errors by removing from state
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		title, detail := handleAPIError(err, "read", resourceId)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	// Map API response to model
	data.ID = types.StringValue(resource.ID)
	data.Name = types.StringValue(resource.Name)
	data.TenantID = types.StringValue(r.service.TenantID())

	// Handle props response
	if resource.Props != nil {
		propsJson, err := json.Marshal(resource.Props)
		if err != nil {
			resp.Diagnostics.AddError(
				"Props Serialization Error",
				fmt.Sprintf("Failed to serialize props: %s", err.Error()),
			)
			return
		}
		data.Props = types.StringValue(string(propsJson))
	} else {
		data.Props = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IAMResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IAMResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource ID from state
	resourceId := data.ID.ValueString()
	if resourceId == "" {
		resp.Diagnostics.AddError(
			"Missing Resource ID",
			"Resource ID is required for update operation",
		)
		return
	}

	// Validate props JSON if provided
	var propsData interface{}
	if !data.Props.IsNull() && !data.Props.IsUnknown() {
		propsStr := data.Props.ValueString()
		if err := ValidateJSONString(propsStr); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Props",
				fmt.Sprintf("Props must be valid JSON: %s", err.Error()),
			)
			return
		}

		// Parse JSON for API call
		if propsStr != "" {
			if err := json.Unmarshal([]byte(propsStr), &propsData); err != nil {
				resp.Diagnostics.AddError(
					"Props JSON Parse Error",
					fmt.Sprintf("Failed to parse props JSON: %s", err.Error()),
				)
				return
			}
		}
	}

	// Create API request
	updateRequest := &iam.SetResourceDto{
		Name:  data.Name.ValueString(),
		Props: propsData,
	}

	// Call SetResource API (PUT endpoint - same as create)
	updatedResource, err := r.service.SetResource(ctx, resourceId, updateRequest)
	if err != nil {
		title, detail := handleAPIError(err, "update", resourceId)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	// Map response to model
	data.ID = types.StringValue(updatedResource.ID)
	data.Name = types.StringValue(updatedResource.Name)
	data.TenantID = types.StringValue(r.service.TenantID())

	// Handle props response
	if updatedResource.Props != nil {
		propsJson, err := json.Marshal(updatedResource.Props)
		if err != nil {
			resp.Diagnostics.AddError(
				"Props Serialization Error",
				fmt.Sprintf("Failed to serialize props: %s", err.Error()),
			)
			return
		}
		data.Props = types.StringValue(string(propsJson))
	} else {
		data.Props = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IAMResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IAMResourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource ID from state
	resourceId := data.ID.ValueString()
	if resourceId == "" {
		resp.Diagnostics.AddError(
			"Missing Resource ID",
			"Resource ID is required for delete operation",
		)
		return
	}

	// Call DeleteResource API
	err := r.service.DeleteResource(ctx, resourceId)
	if err != nil {
		// If resource is already gone (404), that's okay for delete
		if client.IsNotFoundError(err) {
			return
		}

		title, detail := handleAPIError(err, "delete", resourceId)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	// Resource is automatically removed from state by Terraform framework
}

func (r *IAMResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate ID format according to resource ID rules
	if err := validateResourceID(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Resource ID",
			fmt.Sprintf("Resource ID '%s' is invalid: %s", req.ID, err.Error()),
		)
		return
	}

	// Use standard passthrough import after validation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// validateResourceID validates a resource ID according to the rules:
// - Must be 1-1500 characters
// - Cannot contain forward slashes
// - Cannot be '.' or '..'
// - Cannot contain consecutive underscores '__'
func validateResourceID(id string) error {
	if id == "" {
		return fmt.Errorf("resource ID cannot be empty")
	}

	if len(id) > 1500 {
		return fmt.Errorf("resource ID too long (max 1500 characters)")
	}

	if id == "." || id == ".." {
		return fmt.Errorf("resource ID cannot be '.' or '..'")
	}

	if strings.Contains(id, "/") {
		return fmt.Errorf("resource ID cannot contain forward slashes")
	}

	if strings.Contains(id, "__") {
		return fmt.Errorf("resource ID cannot contain consecutive underscores '__'")
	}

	return nil
}

// handleAPIError provides comprehensive error handling for IAM API responses
func handleAPIError(err error, operation, resourceId string) (string, string) {
	if err == nil {
		return "", ""
	}

	errStr := err.Error()

	// Check for common HTTP status codes in error messages
	switch {
	case strings.Contains(errStr, "400") || strings.Contains(errStr, "Bad Request"):
		return "Invalid Request",
			fmt.Sprintf("The %s request for resource '%s' was invalid. This usually means the resource data doesn't meet the API requirements. Please check that:\n"+
				"• Resource ID follows the pattern (1-1500 chars, no slashes, no '.', '..', or '__')\n"+
				"• Resource name is not empty\n"+
				"• Props field contains valid JSON if provided\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "401") || strings.Contains(errStr, "Unauthorized"):
		return "Authentication Failed",
			fmt.Sprintf("Authentication failed for %s operation on resource '%s'. Please check that:\n"+
				"• Your OAuth2 credentials are valid and not expired\n"+
				"• Your client has the necessary permissions\n"+
				"• The tenant ID is correct\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "403") || strings.Contains(errStr, "Forbidden"):
		return "Permission Denied",
			fmt.Sprintf("You don't have permission to %s resource '%s'. Please check that:\n"+
				"• Your OAuth2 token includes the required scopes (iam:read, iam:write)\n"+
				"• Your account has the necessary IAM permissions\n"+
				"• You're accessing the correct tenant\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "404") || strings.Contains(errStr, "Not Found"):
		return "Resource Not Found",
			fmt.Sprintf("Resource '%s' was not found during %s operation. This could mean:\n"+
				"• The resource doesn't exist in the specified tenant\n"+
				"• The resource ID is incorrect\n"+
				"• The resource was deleted by another process\n"+
				"Error details: %s", resourceId, operation, errStr)

	case strings.Contains(errStr, "409") || strings.Contains(errStr, "Conflict"):
		return "Resource Conflict",
			fmt.Sprintf("A conflict occurred during %s operation on resource '%s'. This usually means:\n"+
				"• A resource with this ID already exists (for create operations)\n"+
				"• The resource was modified by another process (for update operations)\n"+
				"• There are conflicting constraints or dependencies\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "429") || strings.Contains(errStr, "Too Many Requests"):
		return "Rate Limit Exceeded",
			fmt.Sprintf("Rate limit exceeded for %s operation on resource '%s'. Please:\n"+
				"• Wait before retrying the operation\n"+
				"• Reduce the frequency of API calls\n"+
				"• Contact support if the problem persists\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "500") || strings.Contains(errStr, "Internal Server Error"):
		return "Server Error",
			fmt.Sprintf("An internal server error occurred during %s operation on resource '%s'. Please:\n"+
				"• Retry the operation after a short delay\n"+
				"• Check the HiiRetail service status\n"+
				"• Contact support if the problem persists\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "502") || strings.Contains(errStr, "Bad Gateway"):
		return "Service Unavailable",
			fmt.Sprintf("The IAM service is temporarily unavailable for %s operation on resource '%s'. Please:\n"+
				"• Retry the operation after a short delay\n"+
				"• Check network connectivity\n"+
				"• Verify the service endpoint URL\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "503") || strings.Contains(errStr, "Service Unavailable"):
		return "Service Maintenance",
			fmt.Sprintf("The IAM service is under maintenance during %s operation on resource '%s'. Please:\n"+
				"• Retry the operation later\n"+
				"• Check the service status page\n"+
				"• Plan operations during maintenance windows\n"+
				"Error details: %s", operation, resourceId, errStr)

	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "context deadline exceeded"):
		return "Request Timeout",
			fmt.Sprintf("The %s operation on resource '%s' timed out. Please:\n"+
				"• Check network connectivity\n"+
				"• Retry with a longer timeout if possible\n"+
				"• Verify the service is responsive\n"+
				"Error details: %s", operation, resourceId, errStr)

	default:
		return "Unexpected Error",
			fmt.Sprintf("An unexpected error occurred during %s operation on resource '%s'. Please:\n"+
				"• Check the error details below\n"+
				"• Verify your configuration\n"+
				"• Contact support if the problem persists\n"+
				"Error details: %s", operation, resourceId, errStr)
	}
}

// ValidateJSONString validates that a string contains valid JSON (exported for testing)
func ValidateJSONString(jsonStr string) error {
	if jsonStr == "" {
		return nil // Empty string is valid (represents null/omitted props)
	}

	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js)
}

// jsonValidator implements validator.String for JSON validation
type jsonValidator struct{}

func (v jsonValidator) Description(ctx context.Context) string {
	return "value must be valid JSON"
}

func (v jsonValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be valid JSON"
}

func (v jsonValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if err := ValidateJSONString(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON",
			fmt.Sprintf("Props must be valid JSON: %s", err.Error()),
		)
	}
}

// mapAPIErrorToDiagnostic converts API errors to user-friendly diagnostic messages
func mapAPIErrorToDiagnostic(err error, operation, resourceId string) (string, string) {
	if clientErr, ok := err.(*client.Error); ok {
		switch clientErr.StatusCode {
		case 400:
			return fmt.Sprintf("Invalid Request - %s", operation),
				fmt.Sprintf("The request for resource '%s' was invalid: %s", resourceId, clientErr.Message)
		case 401:
			return fmt.Sprintf("Authentication Failed - %s", operation),
				fmt.Sprintf("Authentication failed for resource '%s': %s", resourceId, clientErr.Message)
		case 403:
			return fmt.Sprintf("Access Denied - %s", operation),
				fmt.Sprintf("Insufficient permissions for resource '%s': %s", resourceId, clientErr.Message)
		case 404:
			return fmt.Sprintf("Resource Not Found - %s", operation),
				fmt.Sprintf("Resource '%s' was not found: %s", resourceId, clientErr.Message)
		case 409:
			return fmt.Sprintf("Resource Conflict - %s", operation),
				fmt.Sprintf("Resource '%s' already exists or has been modified: %s", resourceId, clientErr.Message)
		case 429:
			return fmt.Sprintf("Rate Limited - %s", operation),
				fmt.Sprintf("Too many requests for resource '%s': %s", resourceId, clientErr.Message)
		case 500, 502, 503, 504:
			return fmt.Sprintf("Server Error - %s", operation),
				fmt.Sprintf("Server error for resource '%s': %s", resourceId, clientErr.Message)
		default:
			return fmt.Sprintf("API Error - %s", operation),
				fmt.Sprintf("Unexpected API error for resource '%s' (status %d): %s", resourceId, clientErr.StatusCode, clientErr.Message)
		}
	}

	// Non-client errors (network, parsing, etc.)
	return fmt.Sprintf("Error - %s", operation),
		fmt.Sprintf("Unexpected error for resource '%s': %s", resourceId, err.Error())
}
