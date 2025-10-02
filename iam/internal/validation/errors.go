package validation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ErrorCode represents standardized error identifiers
type ErrorCode string

const (
	// Format Errors
	ErrorInvalidNameFormat       ErrorCode = "INVALID_NAME_FORMAT"
	ErrorInvalidEmailFormat      ErrorCode = "INVALID_EMAIL_FORMAT"
	ErrorInvalidPermissionFormat ErrorCode = "INVALID_PERMISSION_FORMAT"
	ErrorInvalidRoleFormat       ErrorCode = "INVALID_ROLE_FORMAT"
	ErrorInvalidMemberFormat     ErrorCode = "INVALID_MEMBER_FORMAT"

	// Constraint Errors
	ErrorNameTooShort        ErrorCode = "NAME_TOO_SHORT"
	ErrorNameTooLong         ErrorCode = "NAME_TOO_LONG"
	ErrorDescriptionTooShort ErrorCode = "DESCRIPTION_TOO_SHORT"
	ErrorDescriptionTooLong  ErrorCode = "DESCRIPTION_TOO_LONG"
	ErrorPermissionsTooMany  ErrorCode = "PERMISSIONS_TOO_MANY"
	ErrorPermissionsTooFew   ErrorCode = "PERMISSIONS_TOO_FEW"

	// Reference Errors
	ErrorResourceNotFound      ErrorCode = "RESOURCE_NOT_FOUND"
	ErrorResourceAlreadyExists ErrorCode = "RESOURCE_ALREADY_EXISTS"
	ErrorCircularDependency    ErrorCode = "CIRCULAR_DEPENDENCY"
	ErrorInvalidReference      ErrorCode = "INVALID_REFERENCE"

	// Business Logic Errors
	ErrorReservedName      ErrorCode = "RESERVED_NAME"
	ErrorUnknownPermission ErrorCode = "UNKNOWN_PERMISSION"
	ErrorInvalidStage      ErrorCode = "INVALID_STAGE"
	ErrorInvalidCondition  ErrorCode = "INVALID_CONDITION"
)

// Severity represents error severity levels
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Type aliases for backwards compatibility with tests
type ErrorSeverity = Severity

// EnhancedError provides rich error information for better user experience
type EnhancedError struct {
	Code         ErrorCode         `json:"code"`
	Message      string            `json:"message"`
	FieldPath    string            `json:"field_path"`
	CurrentValue interface{}       `json:"current_value"`
	Expected     string            `json:"expected"`
	Examples     []string          `json:"examples"`
	Guidance     string            `json:"guidance"`
	Severity     Severity          `json:"severity"`
	Context      map[string]string `json:"context,omitempty"`
}

// NewEnhancedError creates a new enhanced error with the specified parameters
func NewEnhancedError(code ErrorCode, fieldPath string, currentValue interface{}, message string) *EnhancedError {
	return &EnhancedError{
		Code:         code,
		Message:      message,
		FieldPath:    fieldPath,
		CurrentValue: currentValue,
		Severity:     SeverityError,
		Context:      make(map[string]string),
	}
}

// WithExpected adds expected format information to the error
func (e *EnhancedError) WithExpected(expected string) *EnhancedError {
	e.Expected = expected
	return e
}

// WithExamples adds working examples to the error
func (e *EnhancedError) WithExamples(examples ...string) *EnhancedError {
	e.Examples = examples
	return e
}

// WithGuidance adds actionable guidance for resolving the error
func (e *EnhancedError) WithGuidance(guidance string) *EnhancedError {
	e.Guidance = guidance
	return e
}

// WithSeverity sets the error severity level
func (e *EnhancedError) WithSeverity(severity Severity) *EnhancedError {
	e.Severity = severity
	return e
}

// WithContext adds additional context information
func (e *EnhancedError) WithContext(key, value string) *EnhancedError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// ToDiagnostic converts the enhanced error to a Terraform diagnostic
func (e *EnhancedError) ToDiagnostic() diag.Diagnostic {
	var detail strings.Builder

	// Build detailed error message
	detail.WriteString(e.Message)

	if e.CurrentValue != nil {
		detail.WriteString(fmt.Sprintf("\nCurrent value: '%v'", e.CurrentValue))
	}

	if e.Expected != "" {
		detail.WriteString(fmt.Sprintf("\nExpected: %s", e.Expected))
	}

	if len(e.Examples) > 0 {
		detail.WriteString(fmt.Sprintf("\nExamples: %s", strings.Join(e.Examples, ", ")))
	}

	if e.Guidance != "" {
		detail.WriteString(fmt.Sprintf("\nGuidance: %s", e.Guidance))
	}

	// Add context information
	if len(e.Context) > 0 {
		detail.WriteString("\nAdditional context:")
		for key, value := range e.Context {
			detail.WriteString(fmt.Sprintf("\n  %s: %s", key, value))
		}
	}

	// Add helpful documentation link
	detail.WriteString("\n\nFor more information, see: https://docs.hiiretail.com/terraform/validation-guide")

	// Convert severity to diagnostic severity (unused for now since we use diag.NewAttributeErrorDiagnostic)

	// Parse field path for attribute path
	attributePath := path.Empty()
	if e.FieldPath != "" {
		// Convert dot notation to path expressions
		parts := strings.Split(e.FieldPath, ".")
		for _, part := range parts {
			// Handle array indices
			if strings.Contains(part, "[") && strings.Contains(part, "]") {
				// Extract array name and index
				arrayParts := strings.Split(part, "[")
				if len(arrayParts) == 2 {
					attributePath = attributePath.AtName(arrayParts[0])
					indexStr := strings.TrimSuffix(arrayParts[1], "]")
					// For now, we'll use AtName for array elements since we don't have the actual index
					attributePath = attributePath.AtName(fmt.Sprintf("[%s]", indexStr))
				}
			} else {
				attributePath = attributePath.AtName(part)
			}
		}
	}

	return diag.NewAttributeErrorDiagnostic(
		attributePath,
		fmt.Sprintf("%s (%s)", e.Message, e.Code),
		detail.String(),
	)
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.FieldPath, e.Message)
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid    bool             `json:"valid"`
	Errors   []*EnhancedError `json:"errors,omitempty"`
	Warnings []*EnhancedError `json:"warnings,omitempty"`
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:    true,
		Errors:   make([]*EnhancedError, 0),
		Warnings: make([]*EnhancedError, 0),
	}
}

// AddError adds an error to the validation result
func (vr *ValidationResult) AddError(err *EnhancedError) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, err)
}

// AddWarning adds a warning to the validation result
func (vr *ValidationResult) AddWarning(warning *EnhancedError) {
	warning.WithSeverity(SeverityWarning)
	vr.Warnings = append(vr.Warnings, warning)
}

// HasErrors returns true if the validation result has errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if the validation result has warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// ToDiagnostics converts all errors and warnings to Terraform diagnostics
func (vr *ValidationResult) ToDiagnostics() diag.Diagnostics {
	var diags diag.Diagnostics

	for _, err := range vr.Errors {
		diags.Append(err.ToDiagnostic())
	}

	for _, warning := range vr.Warnings {
		diags.Append(warning.ToDiagnostic())
	}

	return diags
}

// Common error creation helpers

// NewFormatError creates a format validation error
func NewFormatError(fieldPath string, currentValue interface{}, expected string, examples ...string) *EnhancedError {
	return NewEnhancedError(ErrorInvalidNameFormat, fieldPath, currentValue, "Invalid format").
		WithExpected(expected).
		WithExamples(examples...).
		WithGuidance("Please correct the format and try again")
}

// NewConstraintError creates a constraint validation error
func NewConstraintError(code ErrorCode, fieldPath string, currentValue interface{}, constraint string) *EnhancedError {
	return NewEnhancedError(code, fieldPath, currentValue, "Constraint violation").
		WithExpected(constraint).
		WithGuidance("Please adjust the value to meet the constraint")
}

// NewReferenceError creates a reference validation error
func NewReferenceError(fieldPath string, currentValue interface{}, suggestions ...string) *EnhancedError {
	err := NewEnhancedError(ErrorResourceNotFound, fieldPath, currentValue, "Referenced resource not found")
	if len(suggestions) > 0 {
		err.WithExamples(suggestions...).
			WithGuidance(fmt.Sprintf("Did you mean one of these? %s", strings.Join(suggestions, ", ")))
	} else {
		err.WithGuidance("Please verify the resource name and ensure it exists")
	}
	return err
}

// NewBusinessLogicError creates a business logic validation error
func NewBusinessLogicError(code ErrorCode, fieldPath string, currentValue interface{}, message, guidance string) *EnhancedError {
	return NewEnhancedError(code, fieldPath, currentValue, message).
		WithGuidance(guidance)
}
