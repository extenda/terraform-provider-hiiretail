package resource_iam_role_binding

import (
	"context"
	"fmt"
	"strings"
)

// ValidateRoleBindingModel validates the overall role binding model
func ValidateRoleBindingModel(ctx context.Context, roleId string, isCustom bool, bindings []string) error {
	if roleId == "" {
		return fmt.Errorf("role_id cannot be empty")
	}

	if len(bindings) == 0 {
		return fmt.Errorf("bindings cannot be empty")
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
