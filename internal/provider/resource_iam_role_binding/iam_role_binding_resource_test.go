package resource_iam_role_binding

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// local helper for substring checks in tests
func containsSubstring(s, substr string) bool {
	if substr == "" {
		return true
	}
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}

// TestRoleBindingModelValidation tests the validation logic for role binding models
func TestRoleBindingModelValidation(t *testing.T) {
	tests := []struct {
		name          string
		roleId        string
		isCustom      bool
		bindings      []string
		expectedValid bool
		expectedError string
	}{
		{
			name:          "valid custom role binding",
			roleId:        "custom-role-123",
			isCustom:      true,
			bindings:      []string{"user:user-456"},
			expectedValid: true,
		},
		{
			name:          "valid system role binding",
			roleId:        "system-role-789",
			isCustom:      false,
			bindings:      []string{"group:group-101"},
			expectedValid: true,
		},
		{
			name:          "empty role_id should be invalid",
			roleId:        "",
			isCustom:      true,
			bindings:      []string{"user:user-456"},
			expectedValid: false,
			expectedError: "role[0]: role id cannot be empty",
		},
		{
			name:          "empty bindings should be invalid",
			roleId:        "custom-role-123",
			isCustom:      true,
			bindings:      []string{},
			expectedValid: false,
			expectedError: "role[0]: bindings cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := RoleModel{
				Id:       types.StringValue(tt.roleId),
				IsCustom: types.BoolValue(tt.isCustom),
				Bindings: types.ListValueMust(types.StringType, stringSliceToAttrValues(tt.bindings)),
			}
			err := ValidateRoleBindingModel(context.Background(), []RoleModel{role})
			if tt.expectedValid {
				if err != nil {
					t.Errorf("Expected validation to pass for %s, got error: %v", tt.name, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected validation to fail for %s, but got no error", tt.name)
				} else if tt.expectedError != "" {
					if !containsSubstring(err.Error(), tt.expectedError) {
						t.Errorf("Expected error message containing '%s', got '%s'", tt.expectedError, err.Error())
					}
				}
			}
		})
	}
}

// Helper function to convert []string to []attr.Value
func stringSliceToAttrValues(slice []string) []attr.Value {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return values
}

// TestMaxBindingsValidation tests the maximum bindings limit validation
func TestMaxBindingsValidation(t *testing.T) {
	tests := []struct {
		name         string
		bindingCount int
		expectError  bool
		errorMessage string
	}{
		{
			name:         "1 binding should be valid",
			bindingCount: 1,
			expectError:  false,
		},
		{
			name:         "10 bindings should be valid (max allowed)",
			bindingCount: 10,
			expectError:  false,
		},
		{
			name:         "11 bindings should be invalid",
			bindingCount: 11,
			expectError:  true,
			errorMessage: "exceeds maximum allowed bindings (10)",
		},
		{
			name:         "0 bindings should be invalid",
			bindingCount: 0,
			expectError:  true,
			errorMessage: "bindings cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test bindings
			bindings := make([]string, tt.bindingCount)
			for i := 0; i < tt.bindingCount; i++ {
				bindings[i] = fmt.Sprintf("user:user-%d", i)
			}

			// Note: This function will be implemented in Phase 3.3 - Core Implementation
			// For now, this test will fail because ValidateMaxBindings doesn't exist yet
			err := ValidateMaxBindings(bindings)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got: %v", tt.name, err)
				}
			}
		})
	}
}

// TestTenantIsolationLogic tests the tenant isolation logic
func TestTenantIsolationLogic(t *testing.T) {
	tests := []struct {
		name         string
		tenantId     string
		roleId       string
		bindings     []string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid tenant isolation",
			tenantId:    "tenant-123",
			roleId:      "custom-role-456",
			bindings:    []string{"user:user-789"},
			expectError: false,
		},
		{
			name:         "empty tenant_id should be invalid",
			tenantId:     "",
			roleId:       "custom-role-456",
			bindings:     []string{"user:user-789"},
			expectError:  true,
			errorMessage: "tenant_id cannot be empty",
		},
		{
			name:         "mismatched tenant context should be invalid",
			tenantId:     "tenant-123",
			roleId:       "other-tenant-role-456",
			bindings:     []string{"user:user-789"},
			expectError:  true,
			errorMessage: "role does not belong to the specified tenant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This function will be implemented in Phase 3.3 - Core Implementation
			// For now, this test will fail because ValidateTenantIsolation doesn't exist yet
			err := ValidateTenantIsolation(context.Background(), tt.tenantId, tt.roleId, tt.bindings)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got: %v", tt.name, err)
				}
			}
		})
	}
}

// TestBindingFormatValidation tests the binding format validation
func TestBindingFormatValidation(t *testing.T) {
	tests := []struct {
		name         string
		bindings     []string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid user binding",
			bindings:    []string{"user:user-456"},
			expectError: false,
		},
		{
			name:        "valid group binding",
			bindings:    []string{"group:group-789"},
			expectError: false,
		},
		{
			name:        "valid service account binding",
			bindings:    []string{"serviceAccount:service-123"},
			expectError: false,
		},
		{
			name:        "valid multiple bindings",
			bindings:    []string{"user:user-456", "group:group-789"},
			expectError: false,
		},
		{
			name:         "invalid binding format",
			bindings:     []string{"invalid-format"},
			expectError:  true,
			errorMessage: "invalid binding format",
		},
		{
			name:         "empty binding should be invalid",
			bindings:     []string{""},
			expectError:  true,
			errorMessage: "binding cannot be empty",
		},
		{
			name:         "invalid binding type",
			bindings:     []string{"invalidtype:user-456"},
			expectError:  true,
			errorMessage: "invalid binding type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This function will be implemented in Phase 3.3 - Core Implementation
			// For now, this test will fail because ValidateBindingFormat doesn't exist yet
			err := ValidateBindingFormat(tt.bindings)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got: %v", tt.name, err)
				}
			}
		})
	}
}

// Placeholder functions that will be implemented in Phase 3.3 - Core Implementation
// These functions don't exist yet, so tests will fail as expected in TDD approach

// Validation functions are now in validation.go
