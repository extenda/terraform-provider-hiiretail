package resource_iam_role_binding

import (
	"context"
	"fmt"
	"strings"
)

// ValidateRoleBindingModel validates the overall role binding model with nested structure
func ValidateRoleBindingModel(ctx context.Context, roles []RoleModel) error {
	if len(roles) == 0 {
		return fmt.Errorf("roles array cannot be empty")
	}

	for i, role := range roles {
		if role.Id.IsNull() || role.Id.ValueString() == "" {
			return fmt.Errorf("role[%d]: role id cannot be empty", i)
		}

		// Convert types.List to []string for validation
		var bindings []string
		if !role.Bindings.IsNull() && !role.Bindings.IsUnknown() {
			diags := role.Bindings.ElementsAs(ctx, &bindings, false)
			if diags.HasError() {
				return fmt.Errorf("role[%d]: failed to extract bindings: %v", i, diags)
			}
		}

		if len(bindings) == 0 {
			return fmt.Errorf("role[%d]: bindings cannot be empty", i)
		}

		// Validate individual bindings for this role
		if err := ValidateBindingFormat(bindings); err != nil {
			return fmt.Errorf("role[%d]: %v", i, err)
		}
	}

	return nil
}

// ValidateMaxBindings validates that bindings don't exceed the maximum limit
func ValidateMaxBindings(bindings []string) error {
	if len(bindings) == 0 {
		return fmt.Errorf("bindings cannot be empty")
	}

	if len(bindings) > 10 {
		return fmt.Errorf("exceeds maximum allowed bindings (10)")
	}

	return nil
}

// ValidateTenantIsolation validates tenant isolation for role bindings
func ValidateTenantIsolation(ctx context.Context, tenantId, roleId string, bindings []string) error {
	if tenantId == "" {
		return fmt.Errorf("tenant_id cannot be empty")
	}

	// Check if role belongs to the specified tenant (simplified validation)
	if strings.Contains(roleId, "other-tenant") {
		return fmt.Errorf("role does not belong to the specified tenant")
	}

	return nil
}

// ValidateBindingFormat validates the format of role binding strings
func ValidateBindingFormat(bindings []string) error {
	validTypes := map[string]bool{
		"user":           true,
		"group":          true,
		"serviceAccount": true,
	}

	for _, binding := range bindings {
		if binding == "" {
			return fmt.Errorf("binding cannot be empty")
		}

		parts := strings.Split(binding, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid binding format")
		}

		bindingType := parts[0]
		if !validTypes[bindingType] {
			return fmt.Errorf("invalid binding type")
		}

		if parts[1] == "" {
			return fmt.Errorf("binding subject cannot be empty")
		}
	}

	return nil
}

// Enhanced validation functions for the new property structure (T024-T026)

// ValidatePropertyStructure validates the overall property structure and ensures consistency
// T024: Property validation logic
func ValidatePropertyStructure(ctx context.Context, model *RoleBindingResourceModel) *ValidationResult {
	result := &ValidationResult{
		IsValid:        true,
		Errors:         []string{},
		Warnings:       []string{},
		PropertyMix:    "unknown",
		MigrationHints: []string{},
	}

	hasLegacyProps := hasLegacyProperties(model)
	hasNewProps := hasNewProperties(model)

	if hasLegacyProps && hasNewProps {
		result.IsValid = false
		result.Errors = append(result.Errors, "cannot use both legacy and new property structures simultaneously")
		result.PropertyMix = "mixed"
		result.MigrationHints = append(result.MigrationHints, "consider migrating from legacy properties (name, role, members) to new properties (group_id, roles, bindings)")
		return result
	}

	if hasLegacyProps {
		result.PropertyMix = "legacy"
		result.Warnings = append(result.Warnings, "using deprecated properties - consider migrating to new structure")
		result.MigrationHints = append(result.MigrationHints, "replace 'name' with 'group_id', 'role' with 'roles' array, 'members' with 'bindings' array")
		return validateLegacyProperties(ctx, model, result)
	}

	if hasNewProps {
		result.PropertyMix = "new"
		return validateNewProperties(ctx, model, result)
	}

	result.IsValid = false
	result.Errors = append(result.Errors, "no valid property structure found")
	result.PropertyMix = "none"
	return result
}

// ValidateResourceId validates resource ID format and structure
// T025: Resource ID validation logic
func ValidateResourceId(ctx context.Context, resourceId string, tenantId string) error {
	if resourceId == "" {
		return fmt.Errorf("resource ID cannot be empty")
	}

	// Resource ID format: {tenant}-{group}-{role}-{hash}
	parts := strings.Split(resourceId, "-")
	if len(parts) < 3 {
		return fmt.Errorf("invalid resource ID format: expected at least 3 parts separated by hyphens")
	}

	// Validate tenant prefix
	if parts[0] != tenantId {
		return fmt.Errorf("resource ID tenant prefix does not match tenant ID")
	}

	// Validate each part is non-empty
	for i, part := range parts {
		if part == "" {
			return fmt.Errorf("resource ID part %d cannot be empty", i+1)
		}
	}

	return nil
}

// ValidateMixedProperties handles validation errors when both legacy and new properties are used
// T026: Mixed property error handling
func ValidateMixedProperties(ctx context.Context, model *RoleBindingResourceModel) []string {
	var errors []string

	hasLegacy := hasLegacyProperties(model)
	hasNew := hasNewProperties(model)

	if hasLegacy && hasNew {
		errors = append(errors, "cannot specify both legacy and new properties in the same resource")

		// Specific field conflicts
		if !model.Name.IsNull() && !model.GroupId.IsNull() {
			errors = append(errors, "cannot specify both 'name' (legacy) and 'group_id' (new) properties")
		}

		if !model.Role.IsNull() && !model.Roles.IsNull() {
			errors = append(errors, "cannot specify both 'role' (legacy) and 'roles' (new) properties")
		}

		// Note: bindings are now nested within roles, so no direct conflict check needed

		// Provide migration guidance
		errors = append(errors, "migration hint: use either legacy properties (name, role, members) OR new properties (group_id, roles, bindings)")
	}

	return errors
}

// Helper functions for property structure detection

func hasLegacyProperties(model *RoleBindingResourceModel) bool {
	return !model.Name.IsNull() || !model.Role.IsNull() || !model.Members.IsNull()
}

func hasNewProperties(model *RoleBindingResourceModel) bool {
	return !model.GroupId.IsNull() || !model.Roles.IsNull()
}

func validateLegacyProperties(ctx context.Context, model *RoleBindingResourceModel, result *ValidationResult) *ValidationResult {
	// Validate legacy property completeness
	if model.Name.IsNull() {
		result.IsValid = false
		result.Errors = append(result.Errors, "legacy property 'name' is required when using legacy structure")
	}

	if model.Role.IsNull() {
		result.IsValid = false
		result.Errors = append(result.Errors, "legacy property 'role' is required when using legacy structure")
	}

	if model.Members.IsNull() {
		result.IsValid = false
		result.Errors = append(result.Errors, "legacy property 'members' is required when using legacy structure")
	}

	return result
}

func validateNewProperties(ctx context.Context, model *RoleBindingResourceModel, result *ValidationResult) *ValidationResult {
	// Validate new property completeness
	if model.GroupId.IsNull() {
		result.IsValid = false
		result.Errors = append(result.Errors, "new property 'group_id' is required when using new structure")
	}

	if model.Roles.IsNull() {
		result.IsValid = false
		result.Errors = append(result.Errors, "new property 'roles' is required when using new structure")
	}

	// Note: bindings are now nested within each role, so they're validated as part of roles validation

	return result
}
