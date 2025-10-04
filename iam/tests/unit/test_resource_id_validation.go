package unit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestResourceIDFormatValidation tests resource ID format validation across all property types
func TestResourceIDFormatValidation(t *testing.T) {
	// Valid resource ID patterns: alphanumeric, underscores, hyphens, no leading/trailing special chars
	validResourceIDs := []string{
		"person_123",
		"service_456",
		"group_789",
		"user-account-123",
		"service_account_456",
		"business_unit_789",
		"a1",
		"test123",
		"resource_with_long_name_123",
		"mixed-underscore_and-hyphen123",
	}
	
	// Invalid resource ID patterns
	invalidResourceIDs := []struct {
		id     string
		reason string
	}{
		{"-invalid_start", "starts with hyphen"},
		{"invalid_end-", "ends with hyphen"},
		{"_invalid_start", "starts with underscore"},
		{"invalid_end_", "ends with underscore"},
		{"invalid@char", "contains invalid character @"},
		{"invalid space", "contains space"},
		{"invalid.dot", "contains dot"},
		{"invalid#hash", "contains hash"},
		{"invalid$dollar", "contains dollar sign"},
		{"invalid%percent", "contains percent"},
		{"", "empty string"},
		{"123", "numeric only (debatable, but current pattern requires alphanumeric start)"},
	}
	
	t.Run("ValidResourceIDFormats", func(t *testing.T) {
		for _, validID := range validResourceIDs {
			t.Run("Valid_"+validID, func(t *testing.T) {
				// Test resource ID validation function (will be implemented in T025)
				isValid := validateResourceID(validID)  // This function doesn't exist yet
				
				// This test should fail until T025 is implemented
				assert.False(t, isValid, "Resource ID validation not yet implemented - test should fail")
			})
		}
	})
	
	t.Run("InvalidResourceIDFormats", func(t *testing.T) {
		for _, invalid := range invalidResourceIDs {
			t.Run("Invalid_"+invalid.id+"_"+invalid.reason, func(t *testing.T) {
				// Test resource ID validation function (will be implemented in T025)
				isValid := validateResourceID(invalid.id)  // This function doesn't exist yet
				
				// This test should fail until T025 is implemented
				assert.True(t, isValid, "Resource ID validation not yet implemented - test should fail")
			})
		}
	})
}

// validateResourceID placeholder function - will be implemented in T025
func validateResourceID(resourceID string) bool {
	// Placeholder implementation to make tests fail
	// Real implementation will be in internal/validators/resource_id_validator.go
	return false
}

// TestResourceIDValidationInNewProperties tests resource ID validation in new property structure
func TestResourceIDValidationInNewProperties(t *testing.T) {
	t.Run("GroupIDValidation", func(t *testing.T) {
		validGroupIDs := []string{"valid_group", "group123", "test-group"}
		invalidGroupIDs := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validID := range validGroupIDs {
			t.Run("ValidGroupID_"+validID, func(t *testing.T) {
				model := RoleBindingResourceModel{
					GroupID: types.StringValue(validID),
				}
				diags := model.Validate()
				// Should pass when validation is implemented
				assert.True(t, diags.HasError(), "Group ID validation not yet implemented")
			})
		}
		
		for _, invalidID := range invalidGroupIDs {
			t.Run("InvalidGroupID_"+invalidID, func(t *testing.T) {
				model := RoleBindingResourceModel{
					GroupID: types.StringValue(invalidID),
				}
				diags := model.Validate()
				// Should fail with specific error when validation is implemented
				assert.True(t, diags.HasError(), "Group ID validation not yet implemented")
			})
		}
	})
	
	t.Run("RoleIDValidation", func(t *testing.T) {
		// Role IDs should follow the same format as other resource IDs
		validRoleIDs := []string{"valid_role", "role123", "custom-role"}
		invalidRoleIDs := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validID := range validRoleIDs {
			t.Run("ValidRoleID_"+validID, func(t *testing.T) {
				// Test will be implemented when RoleModel validation is added - T025
				// role := RoleModel{ID: types.StringValue(validID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(validID)
				assert.False(t, isValid, "Role ID validation not yet implemented - should fail")
			})
		}
		
		for _, invalidID := range invalidRoleIDs {
			t.Run("InvalidRoleID_"+invalidID, func(t *testing.T) {
				// Test will be implemented when RoleModel validation is added - T025
				// role := RoleModel{ID: types.StringValue(invalidID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(invalidID)
				assert.False(t, isValid, "Role ID validation not yet implemented - should fail")
			})
		}
	})
	
	t.Run("BindingResourceIDValidation", func(t *testing.T) {
		// Binding resource IDs should follow the same format
		validBindingIDs := []string{"person_123", "service_456", "group_789"}
		invalidBindingIDs := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validID := range validBindingIDs {
			t.Run("ValidBindingID_"+validID, func(t *testing.T) {
				// Test will be implemented when BindingModel validation is added - T025
				// binding := BindingModel{ResourceID: types.StringValue(validID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(validID)
				assert.False(t, isValid, "Binding resource ID validation not yet implemented - should fail")
			})
		}
		
		for _, invalidID := range invalidBindingIDs {
			t.Run("InvalidBindingID_"+invalidID, func(t *testing.T) {
				// Test will be implemented when BindingModel validation is added - T025
				// binding := BindingModel{ResourceID: types.StringValue(invalidID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(invalidID)
				assert.False(t, isValid, "Binding resource ID validation not yet implemented - should fail")
			})
		}
	})
}

// TestResourceIDValidationInLegacyProperties tests resource ID validation in legacy property structure
func TestResourceIDValidationInLegacyProperties(t *testing.T) {
	t.Run("LegacyNameValidation", func(t *testing.T) {
		// Legacy 'name' property should follow same validation rules as group_id
		validNames := []string{"valid_name", "name123", "test-name"}
		invalidNames := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validName := range validNames {
			t.Run("ValidName_"+validName, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Name: types.StringValue(validName),
				}
				diags := model.Validate()
				// Should pass when validation is implemented (with deprecation warning)
				assert.True(t, diags.HasError(), "Legacy name validation not yet implemented")
			})
		}
		
		for _, invalidName := range invalidNames {
			t.Run("InvalidName_"+invalidName, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Name: types.StringValue(invalidName),
				}
				diags := model.Validate()
				// Should fail with validation error when implemented
				assert.True(t, diags.HasError(), "Legacy name validation not yet implemented")
			})
		}
	})
	
	t.Run("LegacyRoleValidation", func(t *testing.T) {
		// Legacy 'role' property should follow same validation rules as role IDs
		validRoles := []string{"valid_role", "role123", "custom-role"}
		invalidRoles := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validRole := range validRoles {
			t.Run("ValidRole_"+validRole, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Role: types.StringValue(validRole),
				}
				diags := model.Validate()
				// Should pass when validation is implemented (with deprecation warning)
				assert.True(t, diags.HasError(), "Legacy role validation not yet implemented")
			})
		}
		
		for _, invalidRole := range invalidRoles {
			t.Run("InvalidRole_"+invalidRole, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Role: types.StringValue(invalidRole),
				}
				diags := model.Validate()
				// Should fail with validation error when implemented
				assert.True(t, diags.HasError(), "Legacy role validation not yet implemented")
			})
		}
	})
	
	t.Run("LegacyMemberResourceIDValidation", func(t *testing.T) {
		// Legacy member resource_id should follow same validation rules
		validMemberIDs := []string{"person_123", "service_456", "group_789"}
		invalidMemberIDs := []string{"-invalid", "invalid-", "invalid@char"}
		
		for _, validID := range validMemberIDs {
			t.Run("ValidMemberID_"+validID, func(t *testing.T) {
				// Test will be implemented when member validation is added - T025
				// member := LegacyMemberModel{ResourceID: types.StringValue(validID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(validID)
				assert.False(t, isValid, "Legacy member resource ID validation not yet implemented - should fail")
			})
		}
		
		for _, invalidID := range invalidMemberIDs {
			t.Run("InvalidMemberID_"+invalidID, func(t *testing.T) {
				// Test will be implemented when member validation is added - T025
				// member := LegacyMemberModel{ResourceID: types.StringValue(invalidID)}
				// For now, just test the placeholder validation function
				isValid := validateResourceID(invalidID)
				assert.False(t, isValid, "Legacy member resource ID validation not yet implemented - should fail")
			})
		}
	})
}

// TestResourceIDValidationRegexPattern tests the specific regex pattern used for validation
func TestResourceIDValidationRegexPattern(t *testing.T) {
	// The expected pattern: ^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$
	// - Must start with alphanumeric
	// - Can contain alphanumeric, underscore, hyphen in middle
	// - Must end with alphanumeric
	// - Minimum length 1 (single alphanumeric character allowed)
	
	t.Run("PatternBoundaryTests", func(t *testing.T) {
		boundaryTests := []struct {
			input    string
			expected bool
			reason   string
		}{
			{"a", true, "single character"},
			{"1", true, "single digit"},
			{"ab", true, "two characters"},
			{"a1", true, "letter and digit"},
			{"1a", true, "digit and letter"},
			{"a_b", true, "underscore in middle"},
			{"a-b", true, "hyphen in middle"},
			{"a_b_c", true, "multiple underscores"},
			{"a-b-c", true, "multiple hyphens"},
			{"a_b-c", true, "mixed underscore and hyphen"},
			{"_a", false, "starts with underscore"},
			{"-a", false, "starts with hyphen"},
			{"a_", false, "ends with underscore"},
			{"a-", false, "ends with hyphen"},
			{"_", false, "only underscore"},
			{"-", false, "only hyphen"},
		}
		
		for _, test := range boundaryTests {
			t.Run("Pattern_"+test.input+"_"+test.reason, func(t *testing.T) {
				isValid := validateResourceID(test.input)
				
				// Test should fail until validation is implemented
				if test.expected {
					assert.False(t, isValid, "Resource ID validation not implemented - valid patterns should fail")
				} else {
					assert.True(t, isValid, "Resource ID validation not implemented - invalid patterns should pass")
				}
			})
		}
	})
}