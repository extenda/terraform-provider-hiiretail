package validation

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NameValidator validates resource names
type NameValidator struct {
	MinLength     int
	MaxLength     int
	AllowedChars  *regexp.Regexp
	ReservedNames []string
	ResourceType  string
}

// NewNameValidator creates a new name validator for the given resource type
func NewNameValidator(resourceType string) *NameValidator {
	// Default validation rules based on simple_test.tf patterns
	return &NameValidator{
		MinLength:     3,
		MaxLength:     64,
		AllowedChars:  regexp.MustCompile(`^[a-z0-9-]+$`), // Lowercase letters, numbers, hyphens
		ReservedNames: []string{"admin", "root", "system", "service"},
		ResourceType:  resourceType,
	}
}

// ValidateString implements FieldValidator interface
func (v *NameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Check length constraints
	if len(value) < v.MinLength {
		resp.Diagnostics.Append(
			NewEnhancedError(
				ErrorNameTooShort,
				"name",
				value,
				fmt.Sprintf("%s name too short", strings.Title(v.ResourceType)),
			).WithExpected(
				fmt.Sprintf("At least %d characters", v.MinLength),
			).WithExamples(
				"test-group", "analytics-team", "user-management",
			).WithGuidance(
				"Use a longer, more descriptive name",
			).ToDiagnostic(),
		)
		return
	}

	if len(value) > v.MaxLength {
		resp.Diagnostics.Append(
			NewEnhancedError(
				ErrorNameTooLong,
				"name",
				value,
				fmt.Sprintf("%s name too long", strings.Title(v.ResourceType)),
			).WithExpected(
				fmt.Sprintf("At most %d characters", v.MaxLength),
			).WithGuidance(
				"Use a shorter, more concise name",
			).ToDiagnostic(),
		)
		return
	}

	// Check format constraints
	if !v.AllowedChars.MatchString(value) {
		resp.Diagnostics.Append(
			NewEnhancedError(
				ErrorInvalidNameFormat,
				"name",
				value,
				fmt.Sprintf("%s name contains invalid characters", strings.Title(v.ResourceType)),
			).WithExpected(
				"Lowercase letters, numbers, and hyphens only",
			).WithExamples(
				"test-group", "analytics-team-2", "user-mgmt",
			).WithGuidance(
				"Remove special characters and use lowercase letters with hyphens as separators",
			).ToDiagnostic(),
		)
		return
	}

	// Check reserved names
	for _, reserved := range v.ReservedNames {
		if strings.EqualFold(value, reserved) {
			resp.Diagnostics.Append(
				NewEnhancedError(
					ErrorReservedName,
					"name",
					value,
					fmt.Sprintf("Name '%s' is reserved", value),
				).WithGuidance(
					"Please choose a different name that doesn't conflict with system names",
				).WithExamples(
					"my-group", "team-analytics", "project-alpha",
				).ToDiagnostic(),
			)
			return
		}
	}
}

// Description implements FieldValidator interface
func (v *NameValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Validates %s name format and constraints", v.ResourceType)
}

// MarkdownDescription implements FieldValidator interface
func (v *NameValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Validates %s name format and constraints. Names must be %d-%d characters long, contain only lowercase letters, numbers, and hyphens.",
		v.ResourceType, v.MinLength, v.MaxLength)
}

// DescriptionValidator validates resource descriptions
type DescriptionValidator struct {
	MinLength    int
	MaxLength    int
	ResourceType string
}

// NewDescriptionValidator creates a new description validator for the given resource type
func NewDescriptionValidator(resourceType string) *DescriptionValidator {
	return &DescriptionValidator{
		MinLength:    1,
		MaxLength:    512,
		ResourceType: resourceType,
	}
}

// ValidateString implements FieldValidator interface
func (v *DescriptionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Check length constraints
	if len(value) < v.MinLength {
		resp.Diagnostics.Append(
			NewEnhancedError(
				ErrorDescriptionTooShort,
				"description",
				value,
				fmt.Sprintf("%s description too short", strings.Title(v.ResourceType)),
			).WithExpected(
				fmt.Sprintf("At least %d character", v.MinLength),
			).WithExamples(
				"Group for test users", "Custom role for analytics team",
			).WithGuidance(
				"Provide a meaningful description of the resource's purpose",
			).ToDiagnostic(),
		)
		return
	}

	if len(value) > v.MaxLength {
		resp.Diagnostics.Append(
			NewEnhancedError(
				ErrorDescriptionTooLong,
				"description",
				value,
				fmt.Sprintf("%s description too long", strings.Title(v.ResourceType)),
			).WithExpected(
				fmt.Sprintf("At most %d characters", v.MaxLength),
			).WithGuidance(
				"Provide a more concise description",
			).ToDiagnostic(),
		)
		return
	}
}

// Description implements FieldValidator interface
func (v *DescriptionValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Validates %s description length", v.ResourceType)
}

// MarkdownDescription implements FieldValidator interface
func (v *DescriptionValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Validates %s description length. Descriptions must be %d-%d characters long.",
		v.ResourceType, v.MinLength, v.MaxLength)
}
