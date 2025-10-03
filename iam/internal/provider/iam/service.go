package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/client"
)

// Service provides IAM API operations
type Service struct {
	client   *client.ServiceClient
	tenantID string
}

// NewService creates a new IAM service client
func NewService(apiClient *client.Client, tenantID string) *Service {
	return &Service{
		client:   apiClient.IAMClient(),
		tenantID: tenantID,
	}
}

// Group represents an IAM group
type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Members     []string `json:"members,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// CustomRole represents an IAM custom role
type CustomRole struct {
	ID          string       `json:"id"`
	Name        string       `json:"name,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions"`
	Stage       string       `json:"stage,omitempty"`
	CreatedAt   string       `json:"created_at,omitempty"`
	UpdatedAt   string       `json:"updated_at,omitempty"`
}

// Permission represents a permission in a custom role
type Permission struct {
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// RoleBinding represents an IAM role binding
type RoleBinding struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	Members   []string `json:"members"`
	Condition string   `json:"condition,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// Role represents a basic IAM role (for data sources)
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Stage       string `json:"stage,omitempty"`
	Type        string `json:"type"` // "basic" or "custom"
}

// ListGroupsRequest represents a request to list groups
type ListGroupsRequest struct {
	Filter   string `json:"filter,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Page     int    `json:"page,omitempty"`
}

// ListGroupsResponse represents a response from listing groups
type ListGroupsResponse struct {
	Groups   []Group `json:"groups"`
	NextPage int     `json:"next_page,omitempty"`
	Total    int     `json:"total"`
}

// ListGroups retrieves a list of IAM groups
func (s *Service) ListGroups(ctx context.Context, req *ListGroupsRequest) (*ListGroupsResponse, error) {
	query := make(map[string]string)
	if req.Filter != "" {
		query["filter"] = req.Filter
	}
	if req.PageSize > 0 {
		query["page_size"] = fmt.Sprintf("%d", req.PageSize)
	}
	if req.Page > 0 {
		query["page"] = fmt.Sprintf("%d", req.Page)
	}

	path := fmt.Sprintf("tenants/%s/groups", s.tenantID)
	resp, err := s.client.Get(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var groups []Group
	if err := json.Unmarshal(resp.Body, &groups); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Wrap the groups array in the expected response structure
	result := &ListGroupsResponse{
		Groups: groups,
		Total:  len(groups),
	}

	return result, nil
}

// GetGroup retrieves a specific IAM group by ID
func (s *Service) GetGroup(ctx context.Context, id string) (*Group, error) {
	path := fmt.Sprintf("tenants/%s/groups/%s", s.tenantID, id)
	resp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get group %s: %w", id, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var group Group
	if err := json.Unmarshal(resp.Body, &group); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &group, nil
}

// CreateGroup creates a new IAM group
func (s *Service) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	path := fmt.Sprintf("tenants/%s/groups", s.tenantID)

	// Create a simplified request body without the ID field
	requestBody := map[string]interface{}{
		"name": group.Name,
	}

	// Only add optional fields if they have values
	if group.Description != "" {
		requestBody["description"] = group.Description
	}
	if len(group.Members) > 0 {
		requestBody["members"] = group.Members
	}

	resp, err := s.client.Post(ctx, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var result Group
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// UpdateGroup updates an existing IAM group
func (s *Service) UpdateGroup(ctx context.Context, id string, group *Group) (*Group, error) {
	path := fmt.Sprintf("tenants/%s/groups/%s", s.tenantID, id)

	// Create a simplified request body without the ID field (same as CreateGroup)
	requestBody := map[string]interface{}{
		"name": group.Name,
	}

	// Only add optional fields if they have values
	if group.Description != "" {
		requestBody["description"] = group.Description
	}
	if len(group.Members) > 0 {
		requestBody["members"] = group.Members
	}

	resp, err := s.client.Put(ctx, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to update group %s: %w", id, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	// Handle 204 No Content response (common for successful updates)
	if resp.StatusCode == 204 || len(resp.Body) == 0 {
		// For 204 responses, fetch the updated group data separately
		return s.GetGroup(ctx, id)
	}

	var result Group
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteGroup deletes an IAM group
func (s *Service) DeleteGroup(ctx context.Context, id string) error {
	path := fmt.Sprintf("tenants/%s/groups/%s", s.tenantID, id)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete group %s: %w", id, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

// ListRoles retrieves a list of IAM roles (both basic and custom)
func (s *Service) ListRoles(ctx context.Context, filter string) ([]Role, error) {
	query := make(map[string]string)
	if filter != "" {
		query["filter"] = filter
	}

	resp, err := s.client.Get(ctx, "roles", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var result struct {
		Roles []Role `json:"roles"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Roles, nil
}

// GetRole retrieves a specific IAM role by name
func (s *Service) GetRole(ctx context.Context, name string) (*Role, error) {
	path := fmt.Sprintf("roles/%s", name)
	resp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get role %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var role Role
	if err := json.Unmarshal(resp.Body, &role); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &role, nil
}

// CreateCustomRole creates a new IAM custom role
func (s *Service) CreateCustomRole(ctx context.Context, role *CustomRole) (*CustomRole, error) {
	path := fmt.Sprintf("tenants/%s/roles", s.tenantID)

	// Create a request body that matches the API specification
	requestBody := map[string]interface{}{
		"id":          role.ID,
		"permissions": role.Permissions,
	}

	// Only add optional fields if they have values
	if role.Name != "" {
		requestBody["name"] = role.Name
	}

	resp, err := s.client.Post(ctx, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom role: %w", err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var result CustomRole
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetCustomRole retrieves a specific IAM custom role by name
func (s *Service) GetCustomRole(ctx context.Context, name string) (*CustomRole, error) {
	path := fmt.Sprintf("tenants/%s/roles/%s", s.tenantID, name)
	resp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom role %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var role CustomRole
	if err := json.Unmarshal(resp.Body, &role); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &role, nil
}

// UpdateCustomRole updates an existing IAM custom role
func (s *Service) UpdateCustomRole(ctx context.Context, name string, role *CustomRole) (*CustomRole, error) {
	path := fmt.Sprintf("tenants/%s/roles/%s", s.tenantID, name)

	// Create a request body that matches the API specification
	requestBody := map[string]interface{}{
		"permissions": role.Permissions,
	}

	// Only add optional fields if they have values
	if role.Name != "" {
		requestBody["name"] = role.Name
	}

	resp, err := s.client.Put(ctx, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom role %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	// Handle 204 No Content response (common for successful updates)
	if resp.StatusCode == 204 || len(resp.Body) == 0 {
		// For 204 responses, fetch the updated role data separately
		return s.GetCustomRole(ctx, name)
	}

	var result CustomRole
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteCustomRole deletes an IAM custom role
func (s *Service) DeleteCustomRole(ctx context.Context, name string) error {
	path := fmt.Sprintf("tenants/%s/roles/%s", s.tenantID, name)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete custom role %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

// ListRoleBindings retrieves a list of IAM role bindings
func (s *Service) ListRoleBindings(ctx context.Context, filter string) ([]RoleBinding, error) {
	query := make(map[string]string)
	if filter != "" {
		query["filter"] = filter
	}

	resp, err := s.client.Get(ctx, "bindings", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list role bindings: %w", err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var result struct {
		Bindings []RoleBinding `json:"bindings"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Bindings, nil
}

// GetRoleBinding retrieves a specific IAM role binding by name
// Since role bindings are stored as group role assignments, we parse the binding ID
// (format: "groupId-roleId") to make direct API calls instead of searching all groups
func (s *Service) GetRoleBinding(ctx context.Context, name string) (*RoleBinding, error) {
	// Parse the binding ID to extract groupId and roleId
	// Expected format: "groupId-roleId" (e.g., "EYNaCiYX6WFmoPxXCGMf-custom.TerraformTestShayne")
	parts := strings.Split(name, "-")
	if len(parts) < 2 {
		return nil, &client.Error{
			StatusCode: 400,
			Message:    fmt.Sprintf("invalid role binding ID format: %s", name),
		}
	}
	
	groupID := parts[0]
	roleID := strings.Join(parts[1:], "-") // Handle role IDs that might contain hyphens
	
	// First get the group to get its name
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		if client.IsNotFoundError(err) {
			return nil, &client.Error{
				StatusCode: 404,
				Message:    fmt.Sprintf("role binding %s not found (group not found)", name),
			}
		}
		return nil, fmt.Errorf("failed to get group %s: %w", groupID, err)
	}
	
	// Get roles for this specific group
	path := fmt.Sprintf("../v2/tenants/%s/groups/%s/roles", s.tenantID, groupID)
	resp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for group %s: %w", groupID, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		if client.IsNotFoundError(err) {
			return nil, &client.Error{
				StatusCode: 404,
				Message:    fmt.Sprintf("role binding %s not found", name),
			}
		}
		return nil, err
	}

	// Parse the response as a generic slice
	var roles []map[string]interface{}
	if err := json.Unmarshal(resp.Body, &roles); err != nil {
		return nil, fmt.Errorf("failed to parse roles response: %w", err)
	}

	// Look for the specific role assignment
	for _, role := range roles {
		currentRoleID, ok := role["roleId"].(string)
		if !ok {
			continue
		}
		
		if currentRoleID == roleID {
			// Found the role assignment, reconstruct the binding
			// Determine if it's a custom role
			rolePrefix := "roles/"
			if isCustom, ok := role["isCustom"].(bool); ok && isCustom {
				rolePrefix = "roles/custom."
			}
			
			binding := &RoleBinding{
				ID:        name,
				Name:      name,
				Role:      fmt.Sprintf("%s%s", rolePrefix, roleID),
				Members:   []string{fmt.Sprintf("group:%s", group.Name)},
				Condition: "", // Role bindings don't have conditions in V2 API
				CreatedAt: group.CreatedAt, // Use group creation time as fallback
				UpdatedAt: group.UpdatedAt, // Use group update time as fallback
			}
			return binding, nil
		}
	}

	// Role assignment not found for this group
	return nil, &client.Error{
		StatusCode: 404,
		Message:    fmt.Sprintf("role binding %s not found", name),
	}
}

// CreateRoleBinding creates a new IAM role binding using V2 group role endpoints
func (s *Service) CreateRoleBinding(ctx context.Context, binding *RoleBinding) (*RoleBinding, error) {
	// Extract group ID from members array (expected format: "group:groupName")
	var groupID string
	var groupName string
	for _, member := range binding.Members {
		if strings.HasPrefix(member, "group:") {
			groupName = strings.TrimPrefix(member, "group:")
			// Find the group by name to get its ID
			groupsResp, err := s.ListGroups(ctx, &ListGroupsRequest{})
			if err != nil {
				return nil, fmt.Errorf("failed to list groups to find group '%s': %w", groupName, err)
			}
			for _, group := range groupsResp.Groups {
				if group.Name == groupName {
					groupID = group.ID
					break
				}
			}
			if groupID == "" {
				return nil, fmt.Errorf("group '%s' not found", groupName)
			}
			break
		}
	}
	
	if groupID == "" {
		return nil, fmt.Errorf("no group found in members array - role binding requires a group member")
	}

	// Parse role to extract roleId and determine if it's custom
	roleId := binding.Role
	isCustom := false
	
	// Handle "roles/custom.roleId" format from main.tf
	if strings.HasPrefix(roleId, "roles/custom.") {
		roleId = strings.TrimPrefix(roleId, "roles/custom.")
		isCustom = true
	} else if strings.HasPrefix(roleId, "roles/") {
		roleId = strings.TrimPrefix(roleId, "roles/")
		isCustom = false
	}

	// Create the proper CreateRoleBindingDto payload for V2 API
	payload := map[string]interface{}{
		"roleId":   roleId,
		"isCustom": isCustom,
		"bindings": binding.Members,
	}

	// Use V2 group role endpoint: POST /api/v2/tenants/{tenantId}/groups/{groupId}/roles
	// We need to use the full path including /api/v2 since the client defaults to /api/v1
	path := fmt.Sprintf("../v2/tenants/%s/groups/%s/roles", s.tenantID, groupID)
	resp, err := s.client.Post(ctx, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create role binding for group %s: %w", groupID, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	// The V2 API may return the role assignment, but we construct our response
	// to match the expected RoleBinding format
	result := &RoleBinding{
		ID:      fmt.Sprintf("%s-%s", groupID, roleId), // Create composite ID
		Name:    binding.Name,
		Role:    binding.Role,
		Members: binding.Members,
		// Only set Condition if it's not empty to maintain consistency with Terraform
		Condition: binding.Condition,
	}

	return result, nil
}

// UpdateRoleBinding updates an existing IAM role binding
func (s *Service) UpdateRoleBinding(ctx context.Context, name string, binding *RoleBinding) (*RoleBinding, error) {
	path := fmt.Sprintf("bindings/%s", name)
	resp, err := s.client.Put(ctx, path, binding)
	if err != nil {
		return nil, fmt.Errorf("failed to update role binding %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return nil, err
	}

	var result RoleBinding
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteRoleBinding deletes an IAM role binding
func (s *Service) DeleteRoleBinding(ctx context.Context, name string) error {
	path := fmt.Sprintf("bindings/%s", name)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete role binding %s: %w", name, err)
	}

	if err := client.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}
