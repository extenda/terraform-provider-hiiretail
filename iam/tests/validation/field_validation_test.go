package validation
package validation

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestFieldValidationContract tests the field validation contracts
func TestFieldValidationContract(t *testing.T) {
	// These tests MUST FAIL initially - they test against unimplemented validators
	
	t.Run("GroupNameValidation", func(t *testing.T) {
		testCases := []TestCase{
			{
				Name:          "ValidGroupName",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "name",
				Value:         "test-group-123",
				ExpectedValid: true,
			},
			{
				Name:          "InvalidGroupNameSpecialChars",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "name",
				Value:         "test@group!",
				ExpectedValid: false,
				ExpectedCode:  ErrorInvalidNameFormat,
				ExpectedError: "Group name contains invalid characters",
				Suggestions:   []string{"test-group", "testgroup123"},
				Examples:      []string{"test-group-dev", "analytics-team"},
			},
			{
				Name:          "GroupNameTooShort",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "name",
				Value:         "ab",
				ExpectedValid: false,
				ExpectedCode:  ErrorNameTooShort,
				ExpectedError: "Group name too short",
			},
			{
				Name:          "GroupNameTooLong",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "name",
				Value:         "this-is-a-very-long-group-name-that-exceeds-the-maximum-allowed-length-for-group-names",
				ExpectedValid: false,
				ExpectedCode:  ErrorNameTooLong,
				ExpectedError: "Group name too long",
			},
			{
				Name:          "GroupNameReserved",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "name",
				Value:         "admin",
				ExpectedValid: false,
				ExpectedCode:  ErrorReservedName,
				ExpectedError: "Name is reserved",
				Examples:      []string{"admin-team", "admin-users"},
			},
		}
		
		suite := NewValidationTestSuite(t)
		
		// Register a mock validator that should exist but doesn't yet
		// This will cause the tests to fail as expected
		if validator, exists := suite.registry.GetFieldValidator("hiiretail_iam_group", "name"); !exists {
			t.Logf("Expected group name validator not found - this is expected for TDD")
			// The test should fail here since the validator doesn't exist yet
			assert.Fail(t, "Group name validator not implemented yet - implement in Phase 3.3")
		} else {
			suite.RunTestCases(testCases)
		}
	})
	
	t.Run("CustomRoleNameValidation", func(t *testing.T) {
		testCases := []TestCase{
			{
				Name:          "ValidCustomRoleName",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "name",
				Value:         "analytics-reader",
				ExpectedValid: true,
			},
			{
				Name:          "InvalidCustomRoleNameFormat",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "name",
				Value:         "Analytics Reader!",
				ExpectedValid: false,
				ExpectedCode:  ErrorInvalidNameFormat,
				ExpectedError: "Invalid role name format",
				Examples:      []string{"analytics-reader", "data-analyst"},
			},
			{
				Name:          "CustomRoleNameEmpty",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "name",
				Value:         "",
				ExpectedValid: false,
				ExpectedCode:  ErrorNameTooShort,
				ExpectedError: "Role name cannot be empty",
			},
		}
		
		suite := NewValidationTestSuite(t)
		
		// This should fail since the validator doesn't exist yet
		if validator, exists := suite.registry.GetFieldValidator("hiiretail_iam_custom_role", "name"); !exists {
			t.Logf("Expected custom role name validator not found - this is expected for TDD")
			assert.Fail(t, "Custom role name validator not implemented yet - implement in Phase 3.3")
		} else {
			suite.RunTestCases(testCases)
		}
	})
	
	t.Run("RoleBindingNameValidation", func(t *testing.T) {
		testCases := []TestCase{
			{
				Name:          "ValidRoleBindingName",
				ResourceType:  "hiiretail_iam_role_binding",
				FieldPath:     "name",
				Value:         "analytics-team-binding",
				ExpectedValid: true,
			},
			{
				Name:          "InvalidRoleBindingNameFormat",
				ResourceType:  "hiiretail_iam_role_binding",
				FieldPath:     "name",
				Value:         "Analytics Team Binding!",
				ExpectedValid: false,
				ExpectedCode:  ErrorInvalidNameFormat,
				ExpectedError: "Invalid role binding name format",
				Examples:      []string{"analytics-team-binding", "data-access-binding"},
			},
		}
		
		suite := NewValidationTestSuite(t)
		
		// This should fail since the validator doesn't exist yet
		if validator, exists := suite.registry.GetFieldValidator("hiiretail_iam_role_binding", "name"); !exists {
			t.Logf("Expected role binding name validator not found - this is expected for TDD")
			assert.Fail(t, "Role binding name validator not implemented yet - implement in Phase 3.3")
		} else {
			suite.RunTestCases(testCases)
		}
	})
	
	t.Run("DescriptionValidation", func(t *testing.T) {
		testCases := []TestCase{
			{
				Name:          "ValidDescription",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "description",
				Value:         "Test IAM group created via Terraform",
				ExpectedValid: true,
			},
			{
				Name:          "DescriptionTooShort",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "description",
				Value:         "x",
				ExpectedValid: false,
				ExpectedCode:  ErrorDescriptionTooShort,
				ExpectedError: "Description too short",
			},
			{
				Name:          "DescriptionTooLong",
				ResourceType:  "hiiretail_iam_group",
				FieldPath:     "description",
				Value:         generateLongString(1001), // Assuming 1000 char limit
				ExpectedValid: false,
				ExpectedCode:  ErrorDescriptionTooLong,
				ExpectedError: "Description too long",
			},
		}
		
		suite := NewValidationTestSuite(t)
		
		// This should fail since the validator doesn't exist yet
		if validator, exists := suite.registry.GetFieldValidator("hiiretail_iam_group", "description"); !exists {
			t.Logf("Expected description validator not found - this is expected for TDD")
			assert.Fail(t, "Description validator not implemented yet - implement in Phase 3.3")
		} else {
			suite.RunTestCases(testCases)
		}
	})
	
	t.Run("StageValidation", func(t *testing.T) {
		testCases := []TestCase{
			{
				Name:          "ValidStageGA",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "stage",
				Value:         "GA",
				ExpectedValid: true,
			},
			{
				Name:          "ValidStageBETA",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "stage",
				Value:         "BETA",
				ExpectedValid: true,
			},
			{
				Name:          "InvalidStage",
				ResourceType:  "hiiretail_iam_custom_role",
				FieldPath:     "stage",
				Value:         "INVALID_STAGE",
				ExpectedValid: false,
				ExpectedCode:  ErrorInvalidStage,
				ExpectedError: "Invalid stage value",
				Examples:      []string{"GA", "BETA", "ALPHA", "DEPRECATED"},
			},
		}
		
		suite := NewValidationTestSuite(t)
		
		// This should fail since the validator doesn't exist yet
		if validator, exists := suite.registry.GetFieldValidator("hiiretail_iam_custom_role", "stage"); !exists {
			t.Logf("Expected stage validator not found - this is expected for TDD")
			assert.Fail(t, "Stage validator not implemented yet - implement in Phase 3.3")
		} else {
			suite.RunTestCases(testCases)
		}
	})
}

// TestValidationRequestStructure tests the validation request/response structure
func TestValidationRequestStructure(t *testing.T) {
	ctx := context.Background()
	
	t.Run("ValidationRequestCreation", func(t *testing.T) {
		req := &ValidationRequest{
			ResourceType: "hiiretail_iam_group",
			FieldPath:    "name",
			Value:        "test-group",
			Context: &ValidationContext{
				ResourceConfig: map[string]interface{}{
					"name": "test-group",
				},
				PlanningPhase: true,
				TenantID:      "test-tenant",
			},
		}
		
		assert.Equal(t, "hiiretail_iam_group", req.ResourceType)
		assert.Equal(t, "name", req.FieldPath)
		assert.Equal(t, "test-group", req.Value)
		assert.NotNil(t, req.Context)
		assert.True(t, req.Context.PlanningPhase)
		assert.Equal(t, "test-tenant", req.Context.TenantID)
	})
	
	t.Run("ValidationResponseCreation", func(t *testing.T) {
		resp := &ValidationResponse{
			Valid:       false,
			ErrorCode:   string(ErrorInvalidNameFormat),
			Message:     "Invalid name format",
			Suggestions: []string{"test-group", "test-name"},
			Examples:    []string{"analytics-team", "data-team"},
			Severity:    string(SeverityError),
		}
		
		assert.False(t, resp.Valid)
		assert.Equal(t, string(ErrorInvalidNameFormat), resp.ErrorCode)
		assert.Equal(t, "Invalid name format", resp.Message)
		assert.Len(t, resp.Suggestions, 2)
		assert.Len(t, resp.Examples, 2)
		assert.Equal(t, string(SeverityError), resp.Severity)
	})
}

// TestValidationResultStructure tests the validation result structure
func TestValidationResultStructure(t *testing.T) {
	t.Run("ValidationResultCreation", func(t *testing.T) {
		result := NewValidationResult()
		
		assert.True(t, result.Valid)
		assert.Empty(t, result.Errors)
		assert.Empty(t, result.Warnings)
		assert.False(t, result.HasErrors())
		assert.False(t, result.HasWarnings())
	})
	
	t.Run("AddErrorToResult", func(t *testing.T) {
		result := NewValidationResult()
		
		err := NewEnhancedError(ErrorInvalidNameFormat, "name", "invalid@name", "Invalid name format")
		result.AddError(err)
		
		assert.False(t, result.Valid)
		assert.Len(t, result.Errors, 1)
		assert.True(t, result.HasErrors())
		assert.Equal(t, ErrorInvalidNameFormat, result.Errors[0].Code)
	})
	
	t.Run("AddWarningToResult", func(t *testing.T) {
		result := NewValidationResult()
		
		warning := NewEnhancedError(ErrorInvalidNameFormat, "name", "questionable-name", "Questionable name format")
		result.AddWarning(warning)
		
		assert.True(t, result.Valid) // Warnings don't invalidate
		assert.Len(t, result.Warnings, 1)
		assert.True(t, result.HasWarnings())
		assert.Equal(t, SeverityWarning, result.Warnings[0].Severity)
	})
}

// Helper function to generate long strings for testing
func generateLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

// TestFieldValidationIntegration tests integration between field validators and the registry
func TestFieldValidationIntegration(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()
	
	t.Run("RegisterAndRetrieveFieldValidator", func(t *testing.T) {
		// This should fail since we haven't implemented any validators yet
		validator, exists := registry.GetFieldValidator("hiiretail_iam_group", "name")
		assert.False(t, exists, "No validators should be registered yet")
		assert.Nil(t, validator, "Validator should be nil when not found")
		
		// Try to validate a field without a registered validator
		diags := registry.ValidateField(ctx, "hiiretail_iam_group", "name", "test-group")
		assert.Empty(t, diags, "Should return no diagnostics when no validator is registered")
		
		// Mark as expected failure for TDD
		t.Logf("Field validator integration test failed as expected - validators not implemented yet")
	})
	
	t.Run("ValidateResourceWithoutValidator", func(t *testing.T) {
		config := map[string]interface{}{
			"name":        "test-group",
			"description": "Test group",
		}
		
		result := registry.ValidateResource(ctx, "hiiretail_iam_group", config)
		
		// Should return valid result when no validator is registered
		assert.True(t, result.Valid)
		assert.Empty(t, result.Errors)
		
		t.Logf("Resource validation without validator passed - this is expected default behavior")
	})
}