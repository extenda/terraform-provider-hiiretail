package testutils

// Package testutils provides common utilities for testing the IAM provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestEnvironment holds configuration for test execution
type TestEnvironment struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	BaseURL      string
	MockServer   *httptest.Server
}

// SetupTestEnvironment configures environment variables for testing
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	// Set default test environment variables
	env := &TestEnvironment{
		TenantID:     getEnvWithDefault("HIIRETAIL_TENANT_ID", "test-tenant-123"),
		ClientID:     getEnvWithDefault("HIIRETAIL_CLIENT_ID", "test-client-id"),
		ClientSecret: getEnvWithDefault("HIIRETAIL_CLIENT_SECRET", "test-client-secret"),
	}

	// Set environment variables for the provider
	os.Setenv("HIIRETAIL_TENANT_ID", env.TenantID)
	os.Setenv("HIIRETAIL_CLIENT_ID", env.ClientID)
	os.Setenv("HIIRETAIL_CLIENT_SECRET", env.ClientSecret)

	return env
}

// SetupMockServer creates a mock HTTP server for API testing
func (env *TestEnvironment) SetupMockServer(t *testing.T) {
	mux := http.NewServeMux()

	// Mock OAuth2 token endpoint
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Validate Content-Type
		if !strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate grant type
		if r.Form.Get("grant_type") != "client_credentials" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "unsupported_grant_type",
				"error_description": "Only client_credentials grant type is supported",
			})
			return
		}

		// Return mock token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-access-token-12345",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})

	// Mock Groups API endpoints
	mux.HandleFunc("/iam/v1/groups", env.handleGroupsCollection)
	mux.HandleFunc("/iam/v1/groups/", env.handleGroupsResource)

	// Mock Custom Roles API endpoints (T003)
	mux.HandleFunc("/iam/v1/custom-roles", env.handleCustomRolesCollection)
	mux.HandleFunc("/iam/v1/custom-roles/", env.handleCustomRolesResource)

	env.MockServer = httptest.NewServer(mux)
	env.BaseURL = env.MockServer.URL

	// Set the base URL environment variable
	os.Setenv("HIIRETAIL_BASE_URL", env.BaseURL)

	t.Cleanup(func() {
		env.MockServer.Close()
	})
}

// handleGroupsCollection handles requests to /groups (collection operations)
func (env *TestEnvironment) handleGroupsCollection(w http.ResponseWriter, r *http.Request) {
	// Validate authorization header
	if !env.validateAuth(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		env.handleListGroups(w, r)
	case http.MethodPost:
		env.handleCreateGroup(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGroupsResource handles requests to /groups/{id} (individual resource operations)
func (env *TestEnvironment) handleGroupsResource(w http.ResponseWriter, r *http.Request) {
	// Validate authorization header
	if !env.validateAuth(w, r) {
		return
	}

	// Extract group ID from path
	path := strings.TrimPrefix(r.URL.Path, "/iam/v1/groups/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		env.handleGetGroup(w, r, path)
	case http.MethodPut:
		env.handleUpdateGroup(w, r, path)
	case http.MethodDelete:
		env.handleDeleteGroup(w, r, path)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// validateAuth validates the Authorization header
func (env *TestEnvironment) validateAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Authorization header required",
			"code":    "UNAUTHORIZED",
		})
		return false
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bearer token required",
			"code":    "UNAUTHORIZED",
		})
		return false
	}

	return true
}

// handleListGroups handles GET /groups
func (env *TestEnvironment) handleListGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"groups": []map[string]interface{}{
			{
				"id":          "group-123",
				"name":        "test-group",
				"description": "Test group description",
				"status":      "active",
				"tenant_id":   env.TenantID,
				"created_at":  "2025-09-28T10:30:00Z",
				"updated_at":  "2025-09-28T10:30:00Z",
			},
		},
		"total":  1,
		"limit":  20,
		"offset": 0,
	})
}

// handleCreateGroup handles POST /groups
func (env *TestEnvironment) handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON",
			"code":    "VALIDATION_ERROR",
		})
		return
	}

	// Validate required fields
	name, ok := req["name"].(string)
	if !ok || name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Validation failed",
			"code":    "VALIDATION_ERROR",
			"details": []map[string]string{
				{
					"field":   "name",
					"message": "Name is required",
					"code":    "REQUIRED_FIELD",
				},
			},
		})
		return
	}

	// Check for duplicate name (simulate conflict)
	if name == "duplicate-group" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Group with this name already exists",
			"code":    "CONFLICT",
		})
		return
	}

	// Create successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          "group-new-123",
		"name":        name,
		"description": req["description"],
		"status":      "active",
		"tenant_id":   env.TenantID,
		"created_at":  "2025-09-28T15:30:00Z",
		"updated_at":  "2025-09-28T15:30:00Z",
	})
}

// handleGetGroup handles GET /groups/{id}
func (env *TestEnvironment) handleGetGroup(w http.ResponseWriter, r *http.Request, groupID string) {
	// Simulate not found for specific IDs
	if groupID == "nonexistent-group" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Group not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          groupID,
		"name":        "test-group",
		"description": "Test group description",
		"status":      "active",
		"tenant_id":   env.TenantID,
		"created_at":  "2025-09-28T10:30:00Z",
		"updated_at":  "2025-09-28T10:30:00Z",
	})
}

// handleUpdateGroup handles PUT /groups/{id}
func (env *TestEnvironment) handleUpdateGroup(w http.ResponseWriter, r *http.Request, groupID string) {
	// Simulate not found for specific IDs
	if groupID == "nonexistent-group" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Group not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON",
			"code":    "VALIDATION_ERROR",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          groupID,
		"name":        req["name"],
		"description": req["description"],
		"status":      "active",
		"tenant_id":   env.TenantID,
		"created_at":  "2025-09-28T10:30:00Z",
		"updated_at":  "2025-09-28T15:45:00Z",
	})
}

// handleDeleteGroup handles DELETE /groups/{id}
func (env *TestEnvironment) handleDeleteGroup(w http.ResponseWriter, r *http.Request, groupID string) {
	// Simulate not found for specific IDs
	if groupID == "nonexistent-group" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Group not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Details []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ========== Custom Role Handlers (T003) ==========

// handleCustomRolesCollection handles requests to /custom-roles (collection operations)
func (env *TestEnvironment) handleCustomRolesCollection(w http.ResponseWriter, r *http.Request) {
	// Validate authorization header
	if !env.validateAuth(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		env.handleListCustomRoles(w, r)
	case http.MethodPost:
		env.handleCreateCustomRole(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleCustomRolesResource handles requests to /custom-roles/{id} (individual resource operations)
func (env *TestEnvironment) handleCustomRolesResource(w http.ResponseWriter, r *http.Request) {
	// Validate authorization header
	if !env.validateAuth(w, r) {
		return
	}

	// Extract role ID from path
	path := strings.TrimPrefix(r.URL.Path, "/iam/v1/custom-roles/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		env.handleGetCustomRole(w, r, path)
	case http.MethodPut:
		env.handleUpdateCustomRole(w, r, path)
	case http.MethodDelete:
		env.handleDeleteCustomRole(w, r, path)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleListCustomRoles handles GET /custom-roles
func (env *TestEnvironment) handleListCustomRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"custom_roles": []map[string]interface{}{
			{
				"id":        "role-123",
				"name":      "test-custom-role",
				"tenant_id": env.TenantID,
				"permissions": []map[string]interface{}{
					{
						"id":    "pos.payment.create",
						"alias": "Create Payment",
						"attributes": map[string]string{
							"department": "finance",
						},
					},
				},
				"created_at": "2025-09-28T10:30:00Z",
				"updated_at": "2025-09-28T10:30:00Z",
			},
		},
		"total":  1,
		"limit":  20,
		"offset": 0,
	})
}

// handleCreateCustomRole handles POST /custom-roles
func (env *TestEnvironment) handleCreateCustomRole(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON",
			"code":    "VALIDATION_ERROR",
		})
		return
	}

	// Validate required fields
	id, ok := req["id"].(string)
	if !ok || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Validation failed",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "id",
					Message: "ID is required",
					Code:    "REQUIRED_FIELD",
				},
			},
		})
		return
	}

	permissions, ok := req["permissions"].([]interface{})
	if !ok || len(permissions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Validation failed",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "permissions",
					Message: "At least one permission is required",
					Code:    "REQUIRED_FIELD",
				},
			},
		})
		return
	}

	// Validate permission patterns and limits
	if err := env.validatePermissions(permissions); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Check for duplicate ID (simulate conflict)
	if id == "duplicate-role" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Custom role with this ID already exists",
			"code":    "CONFLICT",
		})
		return
	}

	// Create successful response with computed fields
	processedPermissions := env.processPermissions(permissions)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          id,
		"name":        req["name"],
		"tenant_id":   env.TenantID,
		"permissions": processedPermissions,
		"created_at":  "2025-09-28T15:30:00Z",
		"updated_at":  "2025-09-28T15:30:00Z",
	})
}

// handleGetCustomRole handles GET /custom-roles/{id}
func (env *TestEnvironment) handleGetCustomRole(w http.ResponseWriter, r *http.Request, roleID string) {
	// Simulate not found for specific IDs
	if roleID == "nonexistent-role" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Custom role not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        roleID,
		"name":      "Test Custom Role", // Use consistent name with CREATE
		"tenant_id": env.TenantID,
		"permissions": []map[string]interface{}{
			{
				"id":    "pos.payment.create",
				"alias": "Create Payment",
				// No attributes - consistent with CREATE when none provided
			},
		},
		"created_at": "2025-09-28T15:30:00Z", // Use consistent timestamp
		"updated_at": "2025-09-28T15:30:00Z",
	})
}

// handleUpdateCustomRole handles PUT /custom-roles/{id}
func (env *TestEnvironment) handleUpdateCustomRole(w http.ResponseWriter, r *http.Request, roleID string) {
	// Simulate not found for specific IDs
	if roleID == "nonexistent-role" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Custom role not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON",
			"code":    "VALIDATION_ERROR",
		})
		return
	}

	// Validate permissions if provided
	if permissions, ok := req["permissions"].([]interface{}); ok {
		if err := env.validatePermissions(permissions); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(err)
			return
		}
		req["permissions"] = env.processPermissions(permissions)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          roleID,
		"name":        req["name"],
		"tenant_id":   env.TenantID,
		"permissions": req["permissions"],
		"created_at":  "2025-09-28T10:30:00Z",
		"updated_at":  "2025-09-28T15:45:00Z",
	})
}

// handleDeleteCustomRole handles DELETE /custom-roles/{id}
func (env *TestEnvironment) handleDeleteCustomRole(w http.ResponseWriter, r *http.Request, roleID string) {
	// Simulate not found for specific IDs
	if roleID == "nonexistent-role" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Custom role not found",
			"code":    "NOT_FOUND",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validatePermissions validates permission patterns and limits
func (env *TestEnvironment) validatePermissions(permissions []interface{}) *ErrorResponse {
	if len(permissions) > 500 {
		return &ErrorResponse{
			Message: "Too many permissions",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "permissions",
					Message: "Maximum 500 permissions allowed per role",
					Code:    "LIMIT_EXCEEDED",
				},
			},
		}
	}

	posCount := 0
	generalCount := 0

	for i, perm := range permissions {
		permMap, ok := perm.(map[string]interface{})
		if !ok {
			return &ErrorResponse{
				Message: "Invalid permission format",
				Code:    "VALIDATION_ERROR",
				Details: []ValidationError{
					{
						Field:   "permissions",
						Message: "Permission must be an object",
						Code:    "INVALID_FORMAT",
					},
				},
			}
		}

		permId, ok := permMap["id"].(string)
		if !ok || permId == "" {
			return &ErrorResponse{
				Message: "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: []ValidationError{
					{
						Field:   "permissions[" + string(rune(i)) + "].id",
						Message: "Permission ID is required",
						Code:    "REQUIRED_FIELD",
					},
				},
			}
		}

		// Validate permission ID pattern
		if !env.isValidPermissionPattern(permId) {
			return &ErrorResponse{
				Message: "Invalid permission ID format",
				Code:    "VALIDATION_ERROR",
				Details: []ValidationError{
					{
						Field:   "permissions[" + string(rune(i)) + "].id",
						Message: "Permission ID must follow pattern: {systemPrefix}.{resource}.{action}",
						Code:    "INVALID_PATTERN",
					},
				},
			}
		}

		// Count POS vs general permissions
		if strings.HasPrefix(permId, "pos.") {
			posCount++
		} else {
			generalCount++
		}

		// Validate attributes if present
		if attrs, ok := permMap["attributes"].(map[string]interface{}); ok {
			if err := env.validateAttributes(attrs); err != nil {
				return err
			}
		}
	}

	// Check permission limits
	if posCount > 500 {
		return &ErrorResponse{
			Message: "Too many POS permissions",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "permissions",
					Message: "Maximum 500 POS permissions allowed per role",
					Code:    "LIMIT_EXCEEDED",
				},
			},
		}
	}

	if generalCount > 100 {
		return &ErrorResponse{
			Message: "Too many general permissions",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "permissions",
					Message: "Maximum 100 general (non-POS) permissions allowed per role",
					Code:    "LIMIT_EXCEEDED",
				},
			},
		}
	}

	return nil
}

// isValidPermissionPattern validates permission ID pattern
func (env *TestEnvironment) isValidPermissionPattern(permId string) bool {
	// Pattern: ^[a-z][-a-z]{2}\.[a-z][-a-z]{1,15}\.[a-z][-a-z]{1,15}$
	parts := strings.Split(permId, ".")
	if len(parts) != 3 {
		return false
	}

	// Validate system prefix (3 chars)
	systemPrefix := parts[0]
	if len(systemPrefix) != 3 || !env.isValidPatternPart(systemPrefix) {
		return false
	}

	// Validate resource (2-16 chars)
	resource := parts[1]
	if len(resource) < 2 || len(resource) > 16 || !env.isValidPatternPart(resource) {
		return false
	}

	// Validate action (2-16 chars)
	action := parts[2]
	if len(action) < 2 || len(action) > 16 || !env.isValidPatternPart(action) {
		return false
	}

	return true
}

// isValidPatternPart validates individual part of permission pattern
func (env *TestEnvironment) isValidPatternPart(part string) bool {
	for i, char := range part {
		if i == 0 {
			// First character must be a-z
			if char < 'a' || char > 'z' {
				return false
			}
		} else {
			// Other characters can be a-z or hyphen
			if (char < 'a' || char > 'z') && char != '-' {
				return false
			}
		}
	}
	return true
}

// validateAttributes validates attribute constraints
func (env *TestEnvironment) validateAttributes(attrs map[string]interface{}) *ErrorResponse {
	if len(attrs) > 10 {
		return &ErrorResponse{
			Message: "Too many attributes",
			Code:    "VALIDATION_ERROR",
			Details: []ValidationError{
				{
					Field:   "attributes",
					Message: "Maximum 10 attributes allowed per permission",
					Code:    "LIMIT_EXCEEDED",
				},
			},
		}
	}

	for key, value := range attrs {
		if len(key) > 40 {
			return &ErrorResponse{
				Message: "Attribute key too long",
				Code:    "VALIDATION_ERROR",
				Details: []ValidationError{
					{
						Field:   "attributes." + key,
						Message: "Attribute key must be 40 characters or less",
						Code:    "LENGTH_EXCEEDED",
					},
				},
			}
		}

		if valueStr, ok := value.(string); ok {
			if len(valueStr) > 256 {
				return &ErrorResponse{
					Message: "Attribute value too long",
					Code:    "VALIDATION_ERROR",
					Details: []ValidationError{
						{
							Field:   "attributes." + key,
							Message: "Attribute value must be 256 characters or less",
							Code:    "LENGTH_EXCEEDED",
						},
					},
				}
			}
		}
	}

	return nil
}

// processPermissions adds computed fields like alias to permissions
func (env *TestEnvironment) processPermissions(permissions []interface{}) []map[string]interface{} {
	processed := make([]map[string]interface{}, len(permissions))

	for i, perm := range permissions {
		permMap := perm.(map[string]interface{})
		processedPerm := map[string]interface{}{
			"id":    permMap["id"],
			"alias": env.generateAlias(permMap["id"].(string)),
		}

		// Only include attributes if they were provided in the request
		if attrs, ok := permMap["attributes"]; ok && attrs != nil {
			processedPerm["attributes"] = attrs
		}

		processed[i] = processedPerm
	}

	return processed
}

// generateAlias generates a human-readable alias for a permission ID
func (env *TestEnvironment) generateAlias(permId string) string {
	parts := strings.Split(permId, ".")
	if len(parts) != 3 {
		return permId
	}

	// Convert to title case
	resource := strings.Title(parts[1])
	action := strings.Title(parts[2])

	return action + " " + resource
}

// ValidateMockServerReady ensures the mock server is fully operational before running tests
func (env *TestEnvironment) ValidateMockServerReady(t *testing.T) {
	if env.MockServer == nil {
		t.Fatal("Mock server not initialized")
	}

	// Test OAuth2 endpoint to ensure it's responsive
	resp, err := http.Post(env.BaseURL+"/oauth2/token", "application/x-www-form-urlencoded",
		strings.NewReader("grant_type=client_credentials&client_id=test&client_secret=test"))
	if err != nil {
		t.Fatalf("Mock server OAuth2 endpoint not responsive: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Mock server OAuth2 endpoint returned unexpected status: %d", resp.StatusCode)
	}

	t.Logf("Mock server validated and ready at %s", env.BaseURL)
}

// SimulateError configures the mock server to return specific error responses
func (env *TestEnvironment) SimulateError(endpoint string, statusCode int, errorResponse ErrorResponse) {
	// This could be extended to configure specific error responses for testing
	// For now, errors are hardcoded in the handlers above
}
