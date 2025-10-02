package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnhancedErrorMessages tests the enhanced error message functionality
// These tests should fail initially since EnhancedError may not be fully implemented
func TestEnhancedErrorMessages(t *testing.T) {

	t.Run("CreateEnhancedError", func(t *testing.T) {
		// Test creating a basic enhanced error
		err := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"invalid@name",
			"Group name contains invalid characters",
		)

		assert.Equal(t, ErrorInvalidNameFormat, err.Code)
		assert.Equal(t, "name", err.FieldPath)
		assert.Equal(t, "invalid@name", err.CurrentValue)
		assert.Equal(t, "Group name contains invalid characters", err.Message)
		assert.Equal(t, SeverityError, err.Severity)
	})

	t.Run("EnhancedErrorWithDetails", func(t *testing.T) {
		// Test creating an enhanced error with all details
		err := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"invalid@name",
			"Group name contains invalid characters",
		).WithExpected(
			"Lowercase letters, numbers, and hyphens only",
		).WithExamples(
			"test-group", "analytics-team",
		).WithGuidance(
			"Remove special characters and use lowercase letters with hyphens as separators",
		)

		assert.Equal(t, "Lowercase letters, numbers, and hyphens only", err.Expected)
		assert.Equal(t, []string{"test-group", "analytics-team"}, err.Examples)
		assert.Equal(t, "Remove special characters and use lowercase letters with hyphens as separators", err.Guidance)
	})

	t.Run("EnhancedErrorToDiagnostic", func(t *testing.T) {
		// Test converting enhanced error to Terraform diagnostic
		err := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"invalid@name",
			"Group name contains invalid characters",
		).WithExpected(
			"Lowercase letters, numbers, and hyphens only",
		).WithExamples(
			"test-group", "analytics-team",
		).WithGuidance(
			"Remove special characters and use lowercase letters with hyphens as separators",
		)

		diagnostic := err.ToDiagnostic()

		assert.NotNil(t, diagnostic)
		// The summary should contain the error code
		assert.Contains(t, diagnostic.Summary(), string(ErrorInvalidNameFormat))
		// The detail should contain all the information
		detail := diagnostic.Detail()
		assert.Contains(t, detail, "Current value: 'invalid@name'")
		assert.Contains(t, detail, "Expected: Lowercase letters, numbers, and hyphens only")
		assert.Contains(t, detail, "Examples: test-group, analytics-team")
		assert.Contains(t, detail, "Guidance: Remove special characters")
		assert.Contains(t, detail, "https://docs.hiiretail.com/terraform/validation-guide")
	})
}

// TestValidationResult tests the validation result functionality
func TestValidationResult(t *testing.T) {

	t.Run("CreateValidationResult", func(t *testing.T) {
		result := NewValidationResult()

		assert.True(t, result.Valid)
		assert.Empty(t, result.Errors)
		assert.Empty(t, result.Warnings)
		assert.False(t, result.HasErrors())
		assert.False(t, result.HasWarnings())
	})

	t.Run("AddErrorToValidationResult", func(t *testing.T) {
		result := NewValidationResult()

		err := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"invalid@name",
			"Group name contains invalid characters",
		)

		result.AddError(err)

		assert.False(t, result.Valid)
		assert.Len(t, result.Errors, 1)
		assert.True(t, result.HasErrors())
		assert.Equal(t, ErrorInvalidNameFormat, result.Errors[0].Code)
	})

	t.Run("AddWarningToValidationResult", func(t *testing.T) {
		result := NewValidationResult()

		warning := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"questionable-name",
			"Name format could be improved",
		)

		result.AddWarning(warning)

		assert.True(t, result.Valid) // Warnings don't invalidate
		assert.Len(t, result.Warnings, 1)
		assert.True(t, result.HasWarnings())
		assert.Equal(t, SeverityWarning, result.Warnings[0].Severity)
	})

	t.Run("ValidationResultToDiagnostics", func(t *testing.T) {
		result := NewValidationResult()

		// Add an error
		err := NewEnhancedError(
			ErrorInvalidNameFormat,
			"name",
			"invalid@name",
			"Group name contains invalid characters",
		)
		result.AddError(err)

		// Add a warning
		warning := NewEnhancedError(
			ErrorNameTooShort,
			"description",
			"x",
			"Description is very short",
		)
		result.AddWarning(warning)

		diagnostics := result.ToDiagnostics()

		assert.Len(t, diagnostics, 2)
		// First should be error, second should be warning
		assert.Contains(t, diagnostics[0].Summary(), "INVALID_NAME_FORMAT")
		assert.Contains(t, diagnostics[1].Summary(), "NAME_TOO_SHORT")
	})
}

// TestErrorCreationHelpers tests the helper functions for creating errors
func TestErrorCreationHelpers(t *testing.T) {

	t.Run("NewFormatError", func(t *testing.T) {
		err := NewFormatError(
			"name",
			"invalid@name",
			"Lowercase letters, numbers, and hyphens only",
			"test-group", "analytics-team",
		)

		assert.Equal(t, ErrorInvalidNameFormat, err.Code)
		assert.Equal(t, "name", err.FieldPath)
		assert.Equal(t, "invalid@name", err.CurrentValue)
		assert.Equal(t, "Lowercase letters, numbers, and hyphens only", err.Expected)
		assert.Equal(t, []string{"test-group", "analytics-team"}, err.Examples)
		assert.Contains(t, err.Guidance, "Please correct the format")
	})

	t.Run("NewConstraintError", func(t *testing.T) {
		err := NewConstraintError(
			ErrorNameTooShort,
			"name",
			"ab",
			"Must be at least 3 characters long",
		)

		assert.Equal(t, ErrorNameTooShort, err.Code)
		assert.Equal(t, "name", err.FieldPath)
		assert.Equal(t, "ab", err.CurrentValue)
		assert.Equal(t, "Must be at least 3 characters long", err.Expected)
		assert.Contains(t, err.Guidance, "Please adjust the value")
	})

	t.Run("NewReferenceError", func(t *testing.T) {
		err := NewReferenceError(
			"role",
			"roles/custom.nonexistent-role",
			"roles/custom.test-custom-role-unique-id",
		)

		assert.Equal(t, ErrorResourceNotFound, err.Code)
		assert.Equal(t, "role", err.FieldPath)
		assert.Equal(t, "roles/custom.nonexistent-role", err.CurrentValue)
		assert.Contains(t, err.Examples, "roles/custom.test-custom-role-unique-id")
		assert.Contains(t, err.Guidance, "Did you mean one of these?")
	})

	t.Run("NewBusinessLogicError", func(t *testing.T) {
		err := NewBusinessLogicError(
			ErrorReservedName,
			"name",
			"admin",
			"Name 'admin' is reserved",
			"Please choose a different name",
		)

		assert.Equal(t, ErrorReservedName, err.Code)
		assert.Equal(t, "name", err.FieldPath)
		assert.Equal(t, "admin", err.CurrentValue)
		assert.Equal(t, "Name 'admin' is reserved", err.Message)
		assert.Equal(t, "Please choose a different name", err.Guidance)
	})
}

// TestResourceSpecificErrorMessages tests error messages specific to simple_test.tf resources
func TestResourceSpecificErrorMessages(t *testing.T) {

	t.Run("GroupNameErrorMessages", func(t *testing.T) {
		// Test error messages that should be shown for group name validation failures
		testCases := []struct {
			value         string
			expectedCode  ErrorCode
			expectedError string
		}{
			{
				value:         "test@group",
				expectedCode:  ErrorInvalidNameFormat,
				expectedError: "Group name contains invalid characters",
			},
			{
				value:         "ab",
				expectedCode:  ErrorNameTooShort,
				expectedError: "Group name too short",
			},
			{
				value:         "admin",
				expectedCode:  ErrorReservedName,
				expectedError: "Name is reserved",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.value, func(t *testing.T) {
				err := NewEnhancedError(tc.expectedCode, "name", tc.value, tc.expectedError)
				assert.Equal(t, tc.expectedCode, err.Code)
				assert.Contains(t, err.Message, tc.expectedError)
			})
		}
	})

	t.Run("CustomRolePermissionErrorMessages", func(t *testing.T) {
		// Test error messages for permission validation failures
		testCases := []struct {
			permission    string
			expectedCode  ErrorCode
			expectedError string
		}{
			{
				permission:    "iam:groups:read", // wrong format - should be dots
				expectedCode:  ErrorInvalidPermissionFormat,
				expectedError: "Invalid permission format",
			},
			{
				permission:    "invalid.groups.read", // unknown service
				expectedCode:  ErrorUnknownPermission,
				expectedError: "Unknown permission",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.permission, func(t *testing.T) {
				err := NewEnhancedError(tc.expectedCode, "permissions[0]", tc.permission, tc.expectedError)
				assert.Equal(t, tc.expectedCode, err.Code)
				assert.Contains(t, err.Message, tc.expectedError)
			})
		}
	})

	t.Run("RoleBindingReferenceErrorMessages", func(t *testing.T) {
		// Test error messages for role binding reference failures
		testCases := []struct {
			reference     string
			expectedCode  ErrorCode
			expectedError string
		}{
			{
				reference:     "roles/custom.nonexistent-role",
				expectedCode:  ErrorResourceNotFound,
				expectedError: "Referenced role not found",
			},
			{
				reference:     "group:nonexistent-group",
				expectedCode:  ErrorResourceNotFound,
				expectedError: "Referenced group not found",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.reference, func(t *testing.T) {
				err := NewReferenceError("role", tc.reference)
				assert.Equal(t, tc.expectedCode, err.Code)
				assert.Equal(t, tc.reference, err.CurrentValue)
			})
		}
	})
}
