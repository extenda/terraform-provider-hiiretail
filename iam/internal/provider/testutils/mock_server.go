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

// SimulateError configures the mock server to return specific error responses
func (env *TestEnvironment) SimulateError(endpoint string, statusCode int, errorResponse ErrorResponse) {
	// This could be extended to configure specific error responses for testing
	// For now, errors are hardcoded in the handlers above
}
