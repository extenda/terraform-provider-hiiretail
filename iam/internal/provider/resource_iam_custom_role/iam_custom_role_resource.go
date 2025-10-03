package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IamCustomRoleResource{}
var _ resource.ResourceWithImportState = &IamCustomRoleResource{}

func NewIamCustomRoleResource() resource.Resource {
	return &IamCustomRoleResource{}
}

// IamCustomRoleResource defines the resource implementation.
type IamCustomRoleResource struct {
	client   *http.Client
	baseURL  string
	tenantID string
}

// APIClient represents the configuration for making API calls (matches provider)
type APIClient struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}

// API request/response structures
type CustomRoleRequest struct {
	ID          string       `json:"id"`
	Name        string       `json:"name,omitempty"`
	Permissions []Permission `json:"permissions"`
}

type CustomRoleResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	TenantID    string       `json:"tenant_id"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   string       `json:"created_at,omitempty"`
	UpdatedAt   string       `json:"updated_at,omitempty"`
}

type Permission struct {
	ID         string                 `json:"id"`
	Alias      string                 `json:"alias,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Details []ValidationError `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r *IamCustomRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_role"
}

func (r *IamCustomRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Use the generated schema from the generated file
	resp.Schema = IamCustomRoleResourceSchema(ctx)
}

func (r *IamCustomRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Handle provider.APIClient type (due to import cycle, we can't import the provider package)
	// Use reflection to extract the fields we need
	switch client := req.ProviderData.(type) {
	case *APIClient:
		// Direct type match (shouldn't happen in practice due to different packages)
		r.client = client.HTTPClient
		r.baseURL = client.BaseURL
		r.tenantID = client.TenantID
	default:
		// Use reflection to extract fields from provider.APIClient
		if apiClient := extractAPIClientFields(req.ProviderData); apiClient != nil {
			r.client = apiClient.HTTPClient
			r.baseURL = apiClient.BaseURL
			r.tenantID = apiClient.TenantID
		} else {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected provider data with BaseURL, TenantID, and HTTPClient fields, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
			return
		}
	}
}

// extractAPIClientFields uses reflection to extract APIClient fields from provider data
func extractAPIClientFields(providerData interface{}) *APIClient {
	if providerData == nil {
		return nil
	}

	v := reflect.ValueOf(providerData)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	// Extract the fields we need
	baseURLField := v.FieldByName("BaseURL")
	tenantIDField := v.FieldByName("TenantID")
	httpClientField := v.FieldByName("HTTPClient")

	if !baseURLField.IsValid() || !tenantIDField.IsValid() || !httpClientField.IsValid() {
		return nil
	}

	if baseURLField.Type().Kind() != reflect.String ||
		tenantIDField.Type().Kind() != reflect.String ||
		httpClientField.Type() != reflect.TypeOf((*http.Client)(nil)) {
		return nil
	}

	return &APIClient{
		BaseURL:    baseURLField.String(),
		TenantID:   tenantIDField.String(),
		HTTPClient: httpClientField.Interface().(*http.Client),
	}
}

func (r *IamCustomRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IamCustomRoleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	apiReq, err := r.modelToAPIRequest(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Data Conversion Error", err.Error())
		return
	}

	// Make API call to create custom role
	apiResp, err := r.createCustomRole(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create custom role: %s", err.Error()))
		return
	}

	// Convert API response back to Terraform model
	err = r.apiResponseToModel(ctx, apiResp, &data)
	if err != nil {
		resp.Diagnostics.AddError("Data Conversion Error", err.Error())
		return
	}

	// Log the successful creation
	tflog.Debug(ctx, "About to save state in Create", map[string]interface{}{
		"id":                     data.Id.ValueString(),
		"name":                   data.Name.ValueString(),
		"permissions_is_null":    data.Permissions.IsNull(),
		"permissions_is_unknown": data.Permissions.IsUnknown(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IamCustomRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IamCustomRoleModel

	// Read Terraform current state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ID from state
	id := data.Id.ValueString()
	if id == "" {
		// For new resources during planning, the state is empty and we should not make API calls
		// This is expected behavior during terraform plan for new resources
		tflog.Debug(ctx, "Read called for new resource (no ID in state), skipping API call")
		return
	}

	// Make API call to read custom role
	tflog.Debug(ctx, "Calling readCustomRole", map[string]interface{}{
		"id": id,
	})

	apiResp, err := r.readCustomRole(ctx, id)
	if err != nil {
		// If not found, remove from state
		if err.Error() == "custom role not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read custom role: %s", err.Error()))
		return
	}

	tflog.Debug(ctx, "Read API response received", map[string]interface{}{
		"apiResp_ID":        apiResp.ID,
		"apiResp_Name":      apiResp.Name,
		"permissions_count": len(apiResp.Permissions),
	})

	// Convert API response back to Terraform model
	err = r.apiResponseToModel(ctx, apiResp, &data)
	if err != nil {
		resp.Diagnostics.AddError("Data Conversion Error", err.Error())
		return
	}

	tflog.Debug(ctx, "apiResponseToModel completed in Read", map[string]interface{}{
		"data_id":            data.Id.ValueString(),
		"permissions_length": data.Permissions.IsNull(),
	})

	// Log the successful read
	tflog.Trace(ctx, "read custom role resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IamCustomRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IamCustomRoleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ID from state
	id := data.Id.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Invalid State", "Custom role ID is missing from state")
		return
	}

	// Convert Terraform model to API request
	apiReq, err := r.modelToAPIRequest(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Data Conversion Error", err.Error())
		return
	}

	// Make API call to update custom role
	apiResp, err := r.updateCustomRole(ctx, id, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update custom role: %s", err.Error()))
		return
	}

	// Convert API response back to Terraform model
	err = r.apiResponseToModel(ctx, apiResp, &data)
	if err != nil {
		resp.Diagnostics.AddError("Data Conversion Error", err.Error())
		return
	}

	// Log the successful update
	tflog.Trace(ctx, "updated custom role resource", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IamCustomRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IamCustomRoleModel

	// Read Terraform current state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ID from state
	id := data.Id.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Invalid State", "Custom role ID is missing from state")
		return
	}

	// Make API call to delete custom role
	err := r.deleteCustomRole(ctx, id)
	if err != nil {
		// If already deleted (not found), that's OK
		if err.Error() != "custom role not found" {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete custom role: %s", err.Error()))
			return
		}
	}

	// Log the successful deletion
	tflog.Trace(ctx, "deleted custom role resource", map[string]interface{}{
		"id": id,
	})

	// Resource is automatically removed from state when Delete returns without error
}

func (r *IamCustomRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID should be the custom role ID
	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Custom role ID is required for import. Use: terraform import hiiretail_iam_custom_role.example role-id",
		)
		return
	}

	// Set the ID in the state using path.Root
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)

	// Log the successful import
	tflog.Trace(ctx, "imported custom role resource", map[string]interface{}{
		"id": id,
	})

	// The Read operation will be called automatically after ImportState to populate the full state
}

// Helper methods for API communication and data conversion

// modelToAPIRequest converts Terraform model to API request format
func (r *IamCustomRoleResource) modelToAPIRequest(ctx context.Context, data IamCustomRoleModel) (*CustomRoleRequest, error) {
	req := &CustomRoleRequest{
		ID: data.Id.ValueString(),
	}

	// Add name if provided
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		req.Name = data.Name.ValueString()
	}

	// Convert permissions
	var permissions []PermissionsValue
	diags := data.Permissions.ElementsAs(ctx, &permissions, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert permissions: %s", diags[0].Summary())
	}

	req.Permissions = make([]Permission, len(permissions))
	for i, perm := range permissions {
		req.Permissions[i] = Permission{
			ID: perm.Id.ValueString(),
		}

		// Add alias if provided
		if !perm.Alias.IsNull() && !perm.Alias.IsUnknown() {
			req.Permissions[i].Alias = perm.Alias.ValueString()
		}

		// Add attributes if provided
		if !perm.Attributes.IsNull() && !perm.Attributes.IsUnknown() {
			attrs := make(map[string]interface{})
			for key, value := range perm.Attributes.Attributes() {
				if strVal, ok := value.(types.String); ok && !strVal.IsNull() {
					attrs[key] = strVal.ValueString()
				}
			}
			if len(attrs) > 0 {
				req.Permissions[i].Attributes = attrs
			}
		}
	}

	return req, nil
}

// createCustomRole makes API call to create custom role
func (r *IamCustomRoleResource) createCustomRole(ctx context.Context, req *CustomRoleRequest) (*CustomRoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tenants/%s/roles", r.baseURL, r.tenantID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-ID", r.tenantID)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		var errorResp ErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("API error: %s (code: %s)", errorResp.Message, errorResp.Code)
	}

	var resp CustomRoleResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// readCustomRole makes API call to read custom role
func (r *IamCustomRoleResource) readCustomRole(ctx context.Context, id string) (*CustomRoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tenants/%s/roles/%s", r.baseURL, r.tenantID, id)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("X-Tenant-ID", r.tenantID)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("custom role not found")
	}

	if httpResp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("API error: %s (code: %s)", errorResp.Message, errorResp.Code)
	}

	var resp CustomRoleResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// updateCustomRole makes API call to update custom role
func (r *IamCustomRoleResource) updateCustomRole(ctx context.Context, id string, req *CustomRoleRequest) (*CustomRoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tenants/%s/roles/%s", r.baseURL, r.tenantID, id)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-ID", r.tenantID)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("custom role not found")
	}

	if httpResp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("API error: %s (code: %s)", errorResp.Message, errorResp.Code)
	}

	var resp CustomRoleResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// deleteCustomRole makes API call to delete custom role
func (r *IamCustomRoleResource) deleteCustomRole(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/api/v1/tenants/%s/roles/%s", r.baseURL, r.tenantID, id)

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("X-Tenant-ID", r.tenantID)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("custom role not found")
	}

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("API request failed with status %d", httpResp.StatusCode)
		}
		return fmt.Errorf("API error: %s (code: %s)", errorResp.Message, errorResp.Code)
	}

	return nil
}

// apiResponseToModel converts API response to Terraform model
func (r *IamCustomRoleResource) apiResponseToModel(ctx context.Context, apiResp *CustomRoleResponse, data *IamCustomRoleModel) error {
	tflog.Debug(ctx, "apiResponseToModel called", map[string]interface{}{
		"apiResp.ID":        apiResp.ID,
		"apiResp.Name":      apiResp.Name,
		"apiResp.TenantID":  apiResp.TenantID,
		"permissions_count": len(apiResp.Permissions),
	})

	data.Id = types.StringValue(apiResp.ID)
	data.Name = types.StringValue(apiResp.Name)
	data.TenantId = types.StringValue(apiResp.TenantID)

	// Convert permissions back to Terraform format
	permissionsList := make([]PermissionsValue, len(apiResp.Permissions))
	tflog.Debug(ctx, "Converting permissions", map[string]interface{}{
		"permissions_count": len(apiResp.Permissions),
	})

	for i, perm := range apiResp.Permissions {
		tflog.Debug(ctx, "Processing permission", map[string]interface{}{
			"index":           i,
			"perm.ID":         perm.ID,
			"perm.Alias":      perm.Alias,
			"perm.Attributes": perm.Attributes,
		})

		// Always set alias to null since it's computed and not provided in config
		// This ensures consistency between planned and actual state
		aliasValue := types.StringNull()

		permissionsList[i] = PermissionsValue{
			Id:         types.StringValue(perm.ID),
			Alias:      aliasValue,
			Attributes: types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}), // Create empty object to match {} in config
			state:      attr.ValueStateKnown,
		}
	}

	// Convert to types.List
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
		},
	}

	tflog.Debug(ctx, "About to convert to types.List", map[string]interface{}{
		"permissionsList_length": len(permissionsList),
		"permissionType":         fmt.Sprintf("%T", permissionType),
	})

	permissionsListValue, diags := types.ListValueFrom(ctx, permissionType, permissionsList)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert permissions to list", map[string]interface{}{
			"error":  diags[0].Summary(),
			"detail": diags[0].Detail(),
		})
		return fmt.Errorf("failed to convert permissions to list: %s", diags[0].Summary())
	}

	tflog.Debug(ctx, "Successfully converted to types.List", map[string]interface{}{
		"list_is_null":    permissionsListValue.IsNull(),
		"list_is_unknown": permissionsListValue.IsUnknown(),
	})

	data.Permissions = permissionsListValue
	return nil
}
