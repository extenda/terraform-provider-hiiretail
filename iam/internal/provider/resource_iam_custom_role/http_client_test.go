package resource_iam_custom_role

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T017: HTTP Client Unit Tests
// These tests verify HTTP client operations with mock servers

func setupMockServer(t *testing.T, statusCode int, responseBody interface{}) (*httptest.Server, *IamCustomRoleResource) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if responseBody != nil {
			err := json.NewEncoder(w).Encode(responseBody)
			require.NoError(t, err, "Should encode response body")
		}
	}))

	resource := &IamCustomRoleResource{
		client:   server.Client(),
		baseURL:  server.URL,
		tenantID: "test-tenant-123",
	}

	return server, resource
}

func TestCreateCustomRole_Success(t *testing.T) {
	// Setup mock server with successful response
	expectedResponse := CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "Test Custom Role",
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{ID: "pos.payment.create", Alias: "Create Payment"},
		},
		CreatedAt: "2025-09-28T15:30:00Z",
		UpdatedAt: "2025-09-28T15:30:00Z",
	}

	server, resource := setupMockServer(t, http.StatusCreated, expectedResponse)
	defer server.Close()

	// Create test request
	apiReq := &CustomRoleRequest{
		ID:   "test-role-001",
		Name: "Test Custom Role",
		Permissions: []Permission{
			{ID: "pos.payment.create", Alias: "Create Payment"},
		},
	}

	// Execute create
	response, err := resource.createCustomRole(context.Background(), apiReq)

	// Assert success
	assert.NoError(t, err, "createCustomRole should not error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.Equal(t, "test-role-001", response.ID)
	assert.Equal(t, "Test Custom Role", response.Name)
	assert.Equal(t, "test-tenant-123", response.TenantID)
	assert.Len(t, response.Permissions, 1)
}

func TestCreateCustomRole_ValidationError(t *testing.T) {
	// Setup mock server with validation error response
	errorResponse := ErrorResponse{
		Message: "Validation failed",
		Code:    "VALIDATION_ERROR",
		Details: []ValidationError{
			{
				Field:   "id",
				Message: "ID is required",
				Code:    "REQUIRED_FIELD",
			},
		},
	}

	server, resource := setupMockServer(t, http.StatusBadRequest, errorResponse)
	defer server.Close()

	// Create invalid request
	apiReq := &CustomRoleRequest{
		ID:          "", // Invalid empty ID
		Name:        "Test Role",
		Permissions: []Permission{},
	}

	// Execute create
	response, err := resource.createCustomRole(context.Background(), apiReq)

	// Assert error
	assert.Error(t, err, "createCustomRole should error for validation failure")
	assert.Nil(t, response, "Response should be nil on error")
	assert.Contains(t, err.Error(), "Validation failed")
	assert.Contains(t, err.Error(), "VALIDATION_ERROR")
}

func TestCreateCustomRole_ConflictError(t *testing.T) {
	// Setup mock server with conflict error
	errorResponse := map[string]string{
		"message": "Custom role with this ID already exists",
		"code":    "CONFLICT",
	}

	server, resource := setupMockServer(t, http.StatusConflict, errorResponse)
	defer server.Close()

	// Create request with duplicate ID
	apiReq := &CustomRoleRequest{
		ID:   "duplicate-role",
		Name: "Duplicate Role",
		Permissions: []Permission{
			{ID: "pos.payment.create"},
		},
	}

	// Execute create
	response, err := resource.createCustomRole(context.Background(), apiReq)

	// Assert error
	assert.Error(t, err, "createCustomRole should error for conflict")
	assert.Nil(t, response, "Response should be nil on error")
	assert.Contains(t, err.Error(), "Custom role with this ID already exists")
	assert.Contains(t, err.Error(), "CONFLICT")
}

func TestReadCustomRole_Success(t *testing.T) {
	// Setup mock server with successful response
	expectedResponse := CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "Test Custom Role",
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{ID: "pos.payment.create", Alias: "Create Payment"},
			{ID: "pos.payment.read", Alias: "Read Payment"},
		},
	}

	server, resource := setupMockServer(t, http.StatusOK, expectedResponse)
	defer server.Close()

	// Execute read
	response, err := resource.readCustomRole(context.Background(), "test-role-001")

	// Assert success
	assert.NoError(t, err, "readCustomRole should not error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.Equal(t, "test-role-001", response.ID)
	assert.Equal(t, "Test Custom Role", response.Name)
	assert.Equal(t, "test-tenant-123", response.TenantID)
	assert.Len(t, response.Permissions, 2)
}

func TestReadCustomRole_NotFound(t *testing.T) {
	// Setup mock server with not found response
	server, resource := setupMockServer(t, http.StatusNotFound, nil)
	defer server.Close()

	// Execute read
	response, err := resource.readCustomRole(context.Background(), "non-existent-role")

	// Assert error
	assert.Error(t, err, "readCustomRole should error for not found")
	assert.Nil(t, response, "Response should be nil on error")
	assert.Contains(t, err.Error(), "custom role not found")
}

func TestUpdateCustomRole_Success(t *testing.T) {
	// Setup mock server with successful response
	expectedResponse := CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "Updated Custom Role",
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{ID: "pos.payment.create", Alias: "Create Payment"},
			{ID: "pos.payment.update", Alias: "Update Payment"},
		},
		UpdatedAt: "2025-09-28T16:30:00Z",
	}

	server, resource := setupMockServer(t, http.StatusOK, expectedResponse)
	defer server.Close()

	// Create update request
	apiReq := &CustomRoleRequest{
		ID:   "test-role-001",
		Name: "Updated Custom Role",
		Permissions: []Permission{
			{ID: "pos.payment.create", Alias: "Create Payment"},
			{ID: "pos.payment.update", Alias: "Update Payment"},
		},
	}

	// Execute update
	response, err := resource.updateCustomRole(context.Background(), "test-role-001", apiReq)

	// Assert success
	assert.NoError(t, err, "updateCustomRole should not error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.Equal(t, "test-role-001", response.ID)
	assert.Equal(t, "Updated Custom Role", response.Name)
	assert.Len(t, response.Permissions, 2)
}

func TestUpdateCustomRole_NotFound(t *testing.T) {
	// Setup mock server with not found response
	server, resource := setupMockServer(t, http.StatusNotFound, nil)
	defer server.Close()

	// Create update request
	apiReq := &CustomRoleRequest{
		ID:   "non-existent-role",
		Name: "Updated Role",
		Permissions: []Permission{
			{ID: "pos.payment.create"},
		},
	}

	// Execute update
	response, err := resource.updateCustomRole(context.Background(), "non-existent-role", apiReq)

	// Assert error
	assert.Error(t, err, "updateCustomRole should error for not found")
	assert.Nil(t, response, "Response should be nil on error")
	assert.Contains(t, err.Error(), "custom role not found")
}

func TestDeleteCustomRole_Success(t *testing.T) {
	// Setup mock server with successful delete (no content)
	server, resource := setupMockServer(t, http.StatusNoContent, nil)
	defer server.Close()

	// Execute delete
	err := resource.deleteCustomRole(context.Background(), "test-role-001")

	// Assert success
	assert.NoError(t, err, "deleteCustomRole should not error")
}

func TestDeleteCustomRole_SuccessWithOKStatus(t *testing.T) {
	// Setup mock server with successful delete (200 OK)
	server, resource := setupMockServer(t, http.StatusOK, nil)
	defer server.Close()

	// Execute delete
	err := resource.deleteCustomRole(context.Background(), "test-role-001")

	// Assert success
	assert.NoError(t, err, "deleteCustomRole should not error with 200 OK")
}

func TestDeleteCustomRole_NotFound(t *testing.T) {
	// Setup mock server with not found response
	server, resource := setupMockServer(t, http.StatusNotFound, nil)
	defer server.Close()

	// Execute delete
	err := resource.deleteCustomRole(context.Background(), "non-existent-role")

	// Assert error
	assert.Error(t, err, "deleteCustomRole should error for not found")
	assert.Contains(t, err.Error(), "custom role not found")
}
