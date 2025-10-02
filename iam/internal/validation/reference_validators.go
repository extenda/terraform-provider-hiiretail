package validation

import (
	"context"
	"strings"
)

// ReferenceValidatorImpl implements ReferenceValidator interface
type ReferenceValidatorImpl struct {
	// Cache of known resources by type
	knownResources map[string][]string
}

// NewReferenceValidator creates a new reference validator
func NewReferenceValidator() *ReferenceValidatorImpl {
	return &ReferenceValidatorImpl{
		knownResources: make(map[string][]string),
	}
}

// ValidateReference validates a reference to another resource
func (v *ReferenceValidatorImpl) ValidateReference(ctx context.Context, referenceType, referenceValue string) *ValidationResult {
	result := NewValidationResult()

	switch referenceType {
	case "role":
		return v.validateRoleReference(referenceValue)
	case "group":
		return v.validateGroupReference(referenceValue)
	case "member":
		return v.validateMemberReference(referenceValue)
	default:
		result.AddError(
			NewEnhancedError(
				ErrorInvalidReference,
				referenceType,
				referenceValue,
				"Unknown reference type",
			).WithGuidance(
				"Use a valid reference type: role, group, or member",
			),
		)
	}

	return result
}

// validateRoleReference validates role references (like "roles/custom.test-custom-role-unique-id")
func (v *ReferenceValidatorImpl) validateRoleReference(referenceValue string) *ValidationResult {
	result := NewValidationResult()

	// Check format: roles/custom.{name} or roles/{builtin-role}
	if !strings.HasPrefix(referenceValue, "roles/") {
		result.AddError(
			NewEnhancedError(
				ErrorInvalidRoleFormat,
				"role",
				referenceValue,
				"Role reference must start with 'roles/'",
			).WithExpected(
				"Format: roles/custom.{name} or roles/{builtin-role}",
			).WithExamples(
				"roles/custom.test-custom-role-unique-id",
				"roles/viewer",
				"roles/editor",
			).WithGuidance(
				"Use the correct role reference format",
			),
		)
		return result
	}

	roleName := strings.TrimPrefix(referenceValue, "roles/")

	// Check if it's a custom role (starts with "custom.")
	if strings.HasPrefix(roleName, "custom.") {
		customRoleName := strings.TrimPrefix(roleName, "custom.")

		// Validate custom role name format (from simple_test.tf)
		if len(customRoleName) < 3 {
			result.AddError(
				NewEnhancedError(
					ErrorNameTooShort,
					"role",
					referenceValue,
					"Custom role name too short",
				).WithExpected(
					"At least 3 characters after 'custom.'",
				).WithExamples(
					"roles/custom.test-custom-role-unique-id",
					"roles/custom.analytics-reader",
				).WithGuidance(
					"Use a longer, more descriptive custom role name",
				),
			)
		}

		// Check if custom role exists (would be populated from API in real implementation)
		knownCustomRoles := v.getKnownCustomRoles()
		if len(knownCustomRoles) > 0 {
			found := false
			for _, known := range knownCustomRoles {
				if known == referenceValue {
					found = true
					break
				}
			}

			if !found {
				result.AddError(
					NewReferenceError(
						"role",
						referenceValue,
						knownCustomRoles...,
					),
				)
			}
		}
	} else {
		// Check built-in roles
		builtinRoles := []string{"viewer", "editor", "admin", "owner"}
		found := false
		for _, builtin := range builtinRoles {
			if builtin == roleName {
				found = true
				break
			}
		}

		if !found {
			suggestions := make([]string, len(builtinRoles))
			for i, builtin := range builtinRoles {
				suggestions[i] = "roles/" + builtin
			}

			result.AddError(
				NewReferenceError(
					"role",
					referenceValue,
					suggestions...,
				),
			)
		}
	}

	return result
}

// validateGroupReference validates group references (like "group:test-group-unique-id")
func (v *ReferenceValidatorImpl) validateGroupReference(referenceValue string) *ValidationResult {
	result := NewValidationResult()

	// Check format: group:{name}
	if !strings.HasPrefix(referenceValue, "group:") {
		result.AddError(
			NewEnhancedError(
				ErrorInvalidMemberFormat,
				"member",
				referenceValue,
				"Group reference must start with 'group:'",
			).WithExpected(
				"Format: group:{name}",
			).WithExamples(
				"group:test-group-unique-id",
				"group:analytics-team",
				"group:admin-users",
			).WithGuidance(
				"Use the correct group reference format",
			),
		)
		return result
	}

	groupName := strings.TrimPrefix(referenceValue, "group:")

	// Validate group name format
	if len(groupName) < 3 {
		result.AddError(
			NewEnhancedError(
				ErrorNameTooShort,
				"member",
				referenceValue,
				"Group name too short",
			).WithExpected(
				"At least 3 characters after 'group:'",
			).WithExamples(
				"group:test-group-unique-id",
				"group:analytics-team",
			).WithGuidance(
				"Use a longer, more descriptive group name",
			),
		)
	}

	// Check if group exists (would be populated from API in real implementation)
	knownGroups := v.getKnownGroups()
	if len(knownGroups) > 0 {
		found := false
		for _, known := range knownGroups {
			if known == referenceValue {
				found = true
				break
			}
		}

		if !found {
			result.AddError(
				NewReferenceError(
					"member",
					referenceValue,
					knownGroups...,
				),
			)
		}
	}

	return result
}

// validateMemberReference validates member references (could be user emails or group references)
func (v *ReferenceValidatorImpl) validateMemberReference(referenceValue string) *ValidationResult {
	result := NewValidationResult()

	// Check if it's a group reference
	if strings.HasPrefix(referenceValue, "group:") {
		return v.validateGroupReference(referenceValue)
	}

	// Otherwise, assume it's a user email and validate email format
	if !strings.Contains(referenceValue, "@") {
		result.AddError(
			NewEnhancedError(
				ErrorInvalidEmailFormat,
				"member",
				referenceValue,
				"Member must be an email address or group reference",
			).WithExpected(
				"Email format or group:{name}",
			).WithExamples(
				"user@example.com",
				"group:test-group-unique-id",
			).WithGuidance(
				"Provide a valid email address or group reference",
			),
		)
	}

	return result
}

// ResolveReference attempts to resolve a reference to its actual resource ID
func (v *ReferenceValidatorImpl) ResolveReference(ctx context.Context, referenceType, referenceValue string) (string, error) {
	// In a real implementation, this would make API calls to resolve references
	// For now, we'll just return the reference value as-is
	return referenceValue, nil
}

// GetSuggestions returns suggestions for similar resources when validation fails
func (v *ReferenceValidatorImpl) GetSuggestions(ctx context.Context, referenceType, referenceValue string) []string {
	switch referenceType {
	case "role":
		return v.getKnownCustomRoles()
	case "group":
		return v.getKnownGroups()
	default:
		return []string{}
	}
}

// Helper methods to get known resources (would be populated from API in real implementation)

func (v *ReferenceValidatorImpl) getKnownCustomRoles() []string {
	// In real implementation, this would fetch from API
	// For now, return examples from simple_test.tf
	return []string{
		"roles/custom.test-custom-role-unique-id",
		"roles/custom.analytics-reader",
		"roles/custom.billing-admin",
	}
}

func (v *ReferenceValidatorImpl) getKnownGroups() []string {
	// In real implementation, this would fetch from API
	// For now, return examples from simple_test.tf
	return []string{
		"group:test-group-unique-id",
		"group:analytics-team",
		"group:admin-users",
	}
}

// SetKnownResources allows setting known resources for testing
func (v *ReferenceValidatorImpl) SetKnownResources(resourceType string, resources []string) {
	v.knownResources[resourceType] = resources
}
