package validation

import (
	"context"
	"strings"
)

// PermissionValidatorImpl implements PermissionValidator interface
type PermissionValidatorImpl struct {
	// Known services and their available actions
	knownServices map[string][]string
}

// NewPermissionValidator creates a new permission validator
func NewPermissionValidator() *PermissionValidatorImpl {
	return &PermissionValidatorImpl{
		knownServices: map[string][]string{
			"iam":       {"groups.read", "groups.write", "groups.admin", "roles.read", "roles.write", "roles.admin", "members.read", "members.write", "members.admin"},
			"analytics": {"data.read", "data.write", "reports.read", "reports.write"},
			"billing":   {"invoices.read", "invoices.write", "payments.read", "payments.write"},
			"products":  {"catalog.read", "catalog.write", "inventory.read", "inventory.write"},
		},
	}
}

// ValidatePermission validates a permission string format and existence
func (v *PermissionValidatorImpl) ValidatePermission(ctx context.Context, permission, service, context string) *ValidationResult {
	result := NewValidationResult()

	// Check if permission uses dot notation (required format from simple_test.tf)
	if !strings.Contains(permission, ".") {
		result.AddError(
			NewEnhancedError(
				ErrorInvalidPermissionFormat,
				"permission",
				permission,
				"Invalid permission format",
			).WithExpected(
				"Dot notation format: service.resource.action",
			).WithExamples(
				"iam.groups.read", "iam.roles.write", "analytics.data.read",
			).WithGuidance(
				"Use dot notation instead of colon notation",
			),
		)
		return result
	}

	// Split permission into parts
	parts := strings.Split(permission, ".")
	if len(parts) < 3 {
		result.AddError(
			NewEnhancedError(
				ErrorInvalidPermissionFormat,
				"permission",
				permission,
				"Permission must have at least 3 parts",
			).WithExpected(
				"Format: service.resource.action",
			).WithExamples(
				"iam.groups.read", "iam.roles.write", "analytics.data.read",
			).WithGuidance(
				"Specify service, resource, and action separated by dots",
			),
		)
		return result
	}

	serviceName := parts[0]
	resourceAndAction := strings.Join(parts[1:], ".")

	// Check if service is known
	allowedActions, serviceExists := v.knownServices[serviceName]
	if !serviceExists {
		result.AddError(
			NewEnhancedError(
				ErrorUnknownPermission,
				"permission",
				permission,
				"Unknown service in permission",
			).WithExpected(
				"One of the known services",
			).WithExamples(
				v.getKnownServiceExamples()...,
			).WithGuidance(
				"Use a valid service name from the available services",
			),
		)
		return result
	}

	// Check if resource.action combination is valid for the service
	validAction := false
	for _, action := range allowedActions {
		if action == resourceAndAction {
			validAction = true
			break
		}
	}

	if !validAction {
		result.AddError(
			NewEnhancedError(
				ErrorUnknownPermission,
				"permission",
				permission,
				"Unknown resource.action for service",
			).WithExpected(
				"Valid resource.action combination for " + serviceName,
			).WithExamples(
				v.getServiceActionExamples(serviceName)...,
			).WithGuidance(
				"Use a valid resource.action combination for the " + serviceName + " service",
			),
		)
		return result
	}

	return result
}

// NormalizePermission returns the normalized form of a permission
func (v *PermissionValidatorImpl) NormalizePermission(permission string) string {
	// Convert from colon notation to dot notation if needed
	if strings.Contains(permission, ":") && !strings.Contains(permission, ".") {
		return strings.ReplaceAll(permission, ":", ".")
	}

	// Ensure lowercase
	return strings.ToLower(permission)
}

// GetPermissionCategory returns the category of a permission (read, write, admin, etc.)
func (v *PermissionValidatorImpl) GetPermissionCategory(permission string) string {
	normalized := v.NormalizePermission(permission)
	parts := strings.Split(normalized, ".")

	if len(parts) < 3 {
		return "unknown"
	}

	action := parts[len(parts)-1]

	switch action {
	case "read":
		return "read"
	case "write", "create", "update", "delete":
		return "write"
	case "admin", "manage":
		return "admin"
	default:
		return "custom"
	}
}

// GetRelatedPermissions returns related permissions for suggestions
func (v *PermissionValidatorImpl) GetRelatedPermissions(permission string) []string {
	normalized := v.NormalizePermission(permission)
	parts := strings.Split(normalized, ".")

	if len(parts) < 2 {
		return []string{}
	}

	service := parts[0]
	resource := parts[1]

	// Get all permissions for the same service and resource
	var related []string
	if actions, exists := v.knownServices[service]; exists {
		for _, action := range actions {
			actionParts := strings.Split(action, ".")
			if len(actionParts) >= 2 && actionParts[0] == resource {
				fullPermission := service + "." + action
				if fullPermission != normalized {
					related = append(related, fullPermission)
				}
			}
		}
	}

	return related
}

// Helper methods

func (v *PermissionValidatorImpl) getKnownServiceExamples() []string {
	var examples []string
	for service, actions := range v.knownServices {
		if len(actions) > 0 {
			examples = append(examples, service+"."+actions[0])
		}
	}
	return examples
}

func (v *PermissionValidatorImpl) getServiceActionExamples(service string) []string {
	var examples []string
	if actions, exists := v.knownServices[service]; exists {
		for i, action := range actions {
			if i >= 3 { // Limit to 3 examples
				break
			}
			examples = append(examples, service+"."+action)
		}
	}
	return examples
}
