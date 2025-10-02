package validation

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	validation "github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/validation"
)

// TestCase represents a single validation test case
type TestCase struct {
	Name          string
	ResourceType  string
	FieldPath     string
	Value         interface{}
	ExpectedValid bool
	ExpectedCode  validation.ErrorCode
	ExpectedError string
	Suggestions   []string
	Examples      []string
}

// ValidationTestSuite provides utilities for testing validators
type ValidationTestSuite struct {
	t        *testing.T
	registry *ValidatorRegistry
}

// NewValidationTestSuite creates a new validation test suite
func NewValidationTestSuite(t *testing.T) *ValidationTestSuite {
	return &ValidationTestSuite{
		t:        t,
		registry: NewValidatorRegistry(),
	}
}

// RunTestCases executes a series of validation test cases
func (suite *ValidationTestSuite) RunTestCases(testCases []TestCase) {
	for _, tc := range testCases {
		suite.t.Run(tc.Name, func(t *testing.T) {
			suite.runSingleTestCase(t, tc)
		})
	}
}

// runSingleTestCase executes a single test case
func (suite *ValidationTestSuite) runSingleTestCase(t *testing.T, tc TestCase) {
	ctx := context.Background()

	// Create validation request
	req := &ValidationRequest{
		ResourceType: tc.ResourceType,
		FieldPath:    tc.FieldPath,
		Value:        tc.Value,
		Context: &ValidationContext{
			ResourceConfig: map[string]interface{}{
				tc.FieldPath: tc.Value,
			},
			PlanningPhase: true,
		},
	}

	// Run validation
	result := suite.validateRequest(ctx, req)

	// Assert results
	if tc.ExpectedValid {
		assert.True(t, result.Valid, "Expected validation to pass but it failed")
		assert.Empty(t, result.Errors, "Expected no errors but got: %v", result.Errors)
	} else {
		assert.False(t, result.Valid, "Expected validation to fail but it passed")
		assert.NotEmpty(t, result.Errors, "Expected validation errors but got none")

		if len(result.Errors) > 0 {
			err := result.Errors[0]

			if tc.ExpectedCode != "" {
				assert.Equal(t, tc.ExpectedCode, err.Code, "Expected error code %s but got %s", tc.ExpectedCode, err.Code)
			}

			if tc.ExpectedError != "" {
				assert.Contains(t, err.Message, tc.ExpectedError, "Expected error message to contain '%s'", tc.ExpectedError)
			}

			if len(tc.Suggestions) > 0 {
				assert.ElementsMatch(t, tc.Suggestions, err.Examples, "Expected suggestions to match")
			}

			if len(tc.Examples) > 0 {
				assert.NotEmpty(t, err.Examples, "Expected examples to be provided")
			}
		}
	}
}

// validateRequest performs validation for a request
func (suite *ValidationTestSuite) validateRequest(ctx context.Context, req *ValidationRequest) *ValidationResult {
	result := NewValidationResult()

	// Try field-level validation first
	if validator, exists := suite.registry.GetFieldValidator(req.ResourceType, req.FieldPath); exists {
		// Create a mock string request for testing
		stringReq := MockStringRequest{
			ConfigValue: types.StringValue(req.Value.(string)),
			Path:        nil,
		}
		stringResp := &MockStringResponse{}

		validator.ValidateString(ctx, stringReq, stringResp)

		// Convert diagnostics to validation result
		for _, diag := range stringResp.Diagnostics {
			if diag.Severity() == diag.Error {
				result.AddError(NewEnhancedError(
					ErrorInvalidNameFormat, // Default error code for testing
					req.FieldPath,
					req.Value,
					diag.Summary(),
				))
			}
		}
	}

	// Try resource-level validation
	if resourceValidator, exists := suite.registry.GetResourceValidator(req.ResourceType); exists {
		resourceResult := resourceValidator.ValidateResource(ctx, req.Context.ResourceConfig)

		// Merge results
		result.Errors = append(result.Errors, resourceResult.Errors...)
		result.Warnings = append(result.Warnings, resourceResult.Warnings...)
		if len(resourceResult.Errors) > 0 {
			result.Valid = false
		}
	}

	return result
}

// Mock types for testing since we can't import the actual Terraform types yet

// MockStringRequest implements a basic string validation request for testing
type MockStringRequest struct {
	ConfigValue types.String
	Path        interface{} // Would be path.Path in real implementation
}

// MockStringResponse implements a basic string validation response for testing
type MockStringResponse struct {
	Diagnostics diag.Diagnostics
}

// ContractTestRunner runs contract tests against validation implementations
type ContractTestRunner struct {
	t *testing.T
}

// NewContractTestRunner creates a new contract test runner
func NewContractTestRunner(t *testing.T) *ContractTestRunner {
	return &ContractTestRunner{t: t}
}

// TestFieldValidationContract tests the field validation contract
func (runner *ContractTestRunner) TestFieldValidationContract(validator FieldValidator, testCases []TestCase) {
	for _, tc := range testCases {
		runner.t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			// Create mock request
			req := MockStringRequest{
				ConfigValue: types.StringValue(tc.Value.(string)),
				Path:        nil,
			}
			resp := &MockStringResponse{}

			// Run validation
			validator.ValidateString(ctx, req, resp)

			// Check results
			hasErrors := len(resp.Diagnostics) > 0

			if tc.ExpectedValid {
				assert.False(t, hasErrors, "Expected validation to pass but got errors: %v", resp.Diagnostics)
			} else {
				assert.True(t, hasErrors, "Expected validation to fail but it passed")
			}

			// Check description methods
			desc := validator.Description(ctx)
			assert.NotEmpty(t, desc, "Description should not be empty")

			markdownDesc := validator.MarkdownDescription(ctx)
			assert.NotEmpty(t, markdownDesc, "MarkdownDescription should not be empty")
		})
	}
}

// TestReferenceValidationContract tests the reference validation contract
func (runner *ContractTestRunner) TestReferenceValidationContract(validator ReferenceValidator, testCases []ReferenceTestCase) {
	for _, tc := range testCases {
		runner.t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			// Test ValidateReference
			result := validator.ValidateReference(ctx, tc.ReferenceType, tc.ReferenceValue)

			if tc.ExpectedValid {
				assert.True(t, result.Valid, "Expected reference validation to pass")
			} else {
				assert.False(t, result.Valid, "Expected reference validation to fail")
			}

			// Test ResolveReference if expected to be valid
			if tc.ExpectedValid {
				resolvedID, err := validator.ResolveReference(ctx, tc.ReferenceType, tc.ReferenceValue)
				assert.NoError(t, err, "Expected reference resolution to succeed")
				assert.NotEmpty(t, resolvedID, "Expected resolved ID to be non-empty")
			}

			// Test GetSuggestions
			suggestions := validator.GetSuggestions(ctx, tc.ReferenceType, tc.ReferenceValue)
			if tc.ExpectedSuggestions != nil {
				assert.ElementsMatch(t, tc.ExpectedSuggestions, suggestions)
			}
		})
	}
}

// ReferenceTestCase represents a reference validation test case
type ReferenceTestCase struct {
	Name                string
	ReferenceType       string
	ReferenceValue      string
	ExpectedValid       bool
	ExpectedSuggestions []string
}

// TestPermissionValidationContract tests the permission validation contract
func (runner *ContractTestRunner) TestPermissionValidationContract(validator PermissionValidator, testCases []PermissionTestCase) {
	for _, tc := range testCases {
		runner.t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			// Test ValidatePermission
			result := validator.ValidatePermission(ctx, tc.Permission, tc.Service, tc.Context)

			if tc.ExpectedValid {
				assert.True(t, result.Valid, "Expected permission validation to pass")
			} else {
				assert.False(t, result.Valid, "Expected permission validation to fail")
			}

			// Test NormalizePermission
			normalized := validator.NormalizePermission(tc.Permission)
			if tc.ExpectedNormalized != "" {
				assert.Equal(t, tc.ExpectedNormalized, normalized)
			}

			// Test GetPermissionCategory
			category := validator.GetPermissionCategory(tc.Permission)
			if tc.ExpectedCategory != "" {
				assert.Equal(t, tc.ExpectedCategory, category)
			}

			// Test GetRelatedPermissions
			related := validator.GetRelatedPermissions(tc.Permission)
			if tc.ExpectedRelated != nil {
				assert.ElementsMatch(t, tc.ExpectedRelated, related)
			}
		})
	}
}

// PermissionTestCase represents a permission validation test case
type PermissionTestCase struct {
	Name               string
	Permission         string
	Service            string
	Context            string
	ExpectedValid      bool
	ExpectedNormalized string
	ExpectedCategory   string
	ExpectedRelated    []string
}

// AssertValidationError is a helper to assert validation errors in tests
func AssertValidationError(t *testing.T, result *ValidationResult, expectedCode ErrorCode, expectedMessage string) {
	require.False(t, result.Valid, "Expected validation to fail")
	require.NotEmpty(t, result.Errors, "Expected validation errors")

	err := result.Errors[0]
	assert.Equal(t, expectedCode, err.Code, "Expected error code to match")
	assert.Contains(t, err.Message, expectedMessage, "Expected error message to contain expected text")
}

// AssertValidationSuccess is a helper to assert successful validation
func AssertValidationSuccess(t *testing.T, result *ValidationResult) {
	assert.True(t, result.Valid, "Expected validation to succeed")
	assert.Empty(t, result.Errors, "Expected no validation errors")
}

// AssertValidationWarning is a helper to assert validation warnings
func AssertValidationWarning(t *testing.T, result *ValidationResult, expectedMessage string) {
	assert.NotEmpty(t, result.Warnings, "Expected validation warnings")

	warning := result.Warnings[0]
	assert.Equal(t, SeverityWarning, warning.Severity, "Expected warning severity")
	assert.Contains(t, warning.Message, expectedMessage, "Expected warning message to contain expected text")
}
