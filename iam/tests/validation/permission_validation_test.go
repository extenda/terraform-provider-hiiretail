package validation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPermissionValidationContract tests the permission validation contracts
// Based on simple_test.tf which uses dot notation: "iam.groups.read", "iam.groups.write", "iam.roles.read"
func TestPermissionValidationContract(t *testing.T) {

	t.Run("ValidPermissionFormats", func(t *testing.T) {
		testCases := []PermissionTestCase{
			{
				Name:               "ValidIAMGroupsRead",
				Permission:         "iam.groups.read",
				Service:            "iam",
				Context:            "custom_role",
				ExpectedValid:      true,
				ExpectedNormalized: "iam.groups.read",
				ExpectedCategory:   "read",
			},
			{
				Name:               "ValidIAMGroupsWrite",
				Permission:         "iam.groups.write",
				Service:            "iam",
				Context:            "custom_role",
				ExpectedValid:      true,
				ExpectedNormalized: "iam.groups.write",
				ExpectedCategory:   "write",
			},
			{
				Name:               "ValidIAMRolesRead",
				Permission:         "iam.roles.read",
				Service:            "iam",
				Context:            "custom_role",
				ExpectedValid:      true,
				ExpectedNormalized: "iam.roles.read",
				ExpectedCategory:   "read",
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Expected permission validator not found - this is expected for TDD")
			assert.Fail(t, "Permission validator not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestPermissionValidationContract(validator, testCases)
		}
	})

	t.Run("InvalidPermissionFormats", func(t *testing.T) {
		testCases := []PermissionTestCase{
			{
				Name:          "ColonNotationInsteadOfDots",
				Permission:    "iam:groups:read", // wrong format - should be dots
				Service:       "iam",
				Context:       "custom_role",
				ExpectedValid: false,
			},
			{
				Name:          "MissingAction",
				Permission:    "iam.groups", // missing action part
				Service:       "iam",
				Context:       "custom_role",
				ExpectedValid: false,
			},
			{
				Name:          "InvalidService",
				Permission:    "invalid.groups.read", // unknown service
				Service:       "invalid",
				Context:       "custom_role",
				ExpectedValid: false,
			},
			{
				Name:          "InvalidResource",
				Permission:    "iam.invalid.read", // unknown resource
				Service:       "iam",
				Context:       "custom_role",
				ExpectedValid: false,
			},
			{
				Name:          "InvalidAction",
				Permission:    "iam.groups.invalid", // unknown action
				Service:       "iam",
				Context:       "custom_role",
				ExpectedValid: false,
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Expected permission validator not found - this is expected for TDD")
			assert.Fail(t, "Permission validator not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestPermissionValidationContract(validator, testCases)
		}
	})

	t.Run("PermissionSuggestions", func(t *testing.T) {
		testCases := []PermissionTestCase{
			{
				Name:            "TypoInGroups",
				Permission:      "iam.grup.read", // typo: grup instead of groups
				Service:         "iam",
				Context:         "custom_role",
				ExpectedValid:   false,
				ExpectedRelated: []string{"iam.groups.read", "iam.groups.write"},
			},
			{
				Name:            "TypoInRoles",
				Permission:      "iam.role.read", // missing 's' - should be roles
				Service:         "iam",
				Context:         "custom_role",
				ExpectedValid:   false,
				ExpectedRelated: []string{"iam.roles.read", "iam.roles.write"},
			},
			{
				Name:            "ColonFormatSuggestion",
				Permission:      "iam:groups:read", // should suggest dot format
				Service:         "iam",
				Context:         "custom_role",
				ExpectedValid:   false,
				ExpectedRelated: []string{"iam.groups.read"},
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Expected permission validator not found - this is expected for TDD")
			assert.Fail(t, "Permission validator suggestions not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestPermissionValidationContract(validator, testCases)
		}
	})
}

// TestPermissionNormalization tests permission normalization
func TestPermissionNormalization(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()

	t.Run("NormalizeValidPermissions", func(t *testing.T) {
		testCases := map[string]string{
			"iam.groups.read":  "iam.groups.read",
			"iam.groups.write": "iam.groups.write",
			"iam.roles.read":   "iam.roles.read",
		}

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission normalization not implemented yet")
		} else {
			for input, expected := range testCases {
				normalized := validator.NormalizePermission(input)
				assert.Equal(t, expected, normalized, "Permission %s should normalize to %s", input, expected)
			}
		}
	})

	t.Run("NormalizeInvalidPermissions", func(t *testing.T) {
		// Invalid permissions should return empty string or the original
		testCases := []string{
			"invalid.format",
			"iam:groups:read", // colon format
			"",
		}

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission normalization not implemented yet")
		} else {
			for _, input := range testCases {
				normalized := validator.NormalizePermission(input)
				// For invalid permissions, normalization should either return empty or original
				assert.NotEqual(t, "iam.groups.read", normalized, "Invalid permission should not normalize to valid format")
			}
		}
	})
}

// TestPermissionCategories tests permission categorization
func TestPermissionCategories(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()

	t.Run("CategorizeReadPermissions", func(t *testing.T) {
		readPermissions := []string{
			"iam.groups.read",
			"iam.roles.read",
		}

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission categorization not implemented yet")
		} else {
			for _, permission := range readPermissions {
				category := validator.GetPermissionCategory(permission)
				assert.Equal(t, "read", category, "Permission %s should be categorized as read", permission)
			}
		}
	})

	t.Run("CategorizeWritePermissions", func(t *testing.T) {
		writePermissions := []string{
			"iam.groups.write",
			"iam.roles.write",
		}

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission categorization not implemented yet")
		} else {
			for _, permission := range writePermissions {
				category := validator.GetPermissionCategory(permission)
				assert.Equal(t, "write", category, "Permission %s should be categorized as write", permission)
			}
		}
	})
}

// TestPermissionArrayValidation tests validation of permission arrays as used in simple_test.tf
func TestPermissionArrayValidation(t *testing.T) {
	t.Run("ValidPermissionArray", func(t *testing.T) {
		// Exact permissions from simple_test.tf
		permissions := []string{
			"iam.groups.read",
			"iam.groups.write",
			"iam.roles.read",
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission array validation not implemented yet")
		} else {
			for _, permission := range permissions {
				result := validator.ValidatePermission(ctx, permission, "iam", "custom_role")
				assert.True(t, result.Valid, "Permission %s should be valid", permission)
			}
		}
	})

	t.Run("MixedValidInvalidPermissionArray", func(t *testing.T) {
		permissions := []string{
			"iam.groups.read",    // valid
			"iam:groups:write",   // invalid format (colon notation)
			"iam.roles.read",     // valid
			"invalid.permission", // invalid
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the permission validator doesn't exist yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Permission array validation not implemented yet")
		} else {
			validCount := 0
			invalidCount := 0

			for _, permission := range permissions {
				result := validator.ValidatePermission(ctx, permission, "iam", "custom_role")
				if result.Valid {
					validCount++
				} else {
					invalidCount++
				}
			}

			assert.Equal(t, 2, validCount, "Should have 2 valid permissions")
			assert.Equal(t, 2, invalidCount, "Should have 2 invalid permissions")
		}
	})
}

// TestRelatedPermissions tests the suggestion system for related permissions
func TestRelatedPermissions(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()

	t.Run("GroupsPermissionRelated", func(t *testing.T) {
		// This should fail since no validator is implemented yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Related permissions not implemented yet")
		} else {
			related := validator.GetRelatedPermissions("iam.groups.read")

			// Should include other group-related permissions
			expectedRelated := []string{"iam.groups.write"}
			assert.Contains(t, related, expectedRelated[0])
		}
	})

	t.Run("RolesPermissionRelated", func(t *testing.T) {
		// This should fail since no validator is implemented yet
		if validator := registry.GetPermissionValidator(); validator == nil {
			t.Logf("Permission validator not found - expected for TDD")
			assert.Fail(t, "Related permissions not implemented yet")
		} else {
			related := validator.GetRelatedPermissions("iam.roles.read")

			// Should include other role-related permissions
			expectedRelated := []string{"iam.roles.write"}
			assert.Contains(t, related, expectedRelated[0])
		}
	})
}
