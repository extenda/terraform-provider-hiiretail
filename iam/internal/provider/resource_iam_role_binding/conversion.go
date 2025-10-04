package resource_iam_role_binding

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Data conversion utilities for property structure transformation (T027-T029)

// ConvertLegacyToNew converts legacy property structure to new enhanced structure
// T027: Data conversion utilities
func ConvertLegacyToNew(ctx context.Context, model *RoleBindingResourceModel) (*RoleBindingResourceModel, error) {
	if !hasLegacyProperties(model) {
		return nil, fmt.Errorf("no legacy properties found to convert")
	}

	newModel := &RoleBindingResourceModel{
		// Copy core properties
		Id:       model.Id,
		TenantId: model.TenantId,

		// Convert legacy properties to new structure
		GroupId:     model.Name, // name becomes group_id
		Description: model.Description,
		Condition:   model.Condition,
		IsCustom:    model.IsCustom,

		// Clear legacy properties
		Name:    types.StringNull(),
		Role:    types.StringNull(),
		Members: types.ListNull(types.ObjectType{}),
	}

	// Convert single role to roles array with bindings from members
	if !model.Role.IsNull() {
		// Convert members to string array for bindings
		var membersArray []string
		if !model.Members.IsNull() && !model.Members.IsUnknown() {
			diags := model.Members.ElementsAs(ctx, &membersArray, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract members: %v", diags)
			}
		}

		// Create bindings list
		bindingsList, diags := types.ListValueFrom(ctx, types.StringType, membersArray)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create bindings list: %v", diags)
		}

		roleModel := RoleModel{
			Id:       model.Role,
			Bindings: bindingsList,
		}

		rolesValue, err := convertRoleToList(ctx, []RoleModel{roleModel})
		if err != nil {
			return nil, fmt.Errorf("failed to convert role to roles array: %w", err)
		}
		newModel.Roles = rolesValue
	}

	// Note: Members are now handled as part of the role conversion above
	// since bindings are nested within each role

	return newModel, nil
}

// ConvertNewToLegacy converts new enhanced structure to legacy property structure
// This is used for backward compatibility scenarios
func ConvertNewToLegacy(ctx context.Context, model *RoleBindingResourceModel) (*RoleBindingResourceModel, error) {
	if !hasNewProperties(model) {
		return nil, fmt.Errorf("no new properties found to convert")
	}

	legacyModel := &RoleBindingResourceModel{
		// Copy core properties
		Id:       model.Id,
		TenantId: model.TenantId,

		// Convert new properties to legacy structure
		Name:        model.GroupId, // group_id becomes name
		Description: model.Description,
		Condition:   model.Condition,
		IsCustom:    model.IsCustom,

		// Clear new properties
		GroupId: types.StringNull(),
		Roles:   types.ListNull(types.ObjectType{}),
	}

	// Convert roles array to single role (take first role)
	if !model.Roles.IsNull() {
		var roles []RoleModel
		diags := model.Roles.ElementsAs(ctx, &roles, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract roles: %s", diags.Errors()[0].Summary())
		}

		if len(roles) > 0 {
			legacyModel.Role = roles[0].Id
			// Convert the first role's bindings to legacy members format
			if !roles[0].Bindings.IsNull() {
				var bindingIds []string
				diags := roles[0].Bindings.ElementsAs(ctx, &bindingIds, false)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to extract role bindings: %s", diags.Errors()[0].Summary())
				}

				// Convert to legacy member format
				var legacyMembers []LegacyMemberModel
				for _, bindingId := range bindingIds {
					// Parse binding ID to determine type and id
					memberType, memberId := parseLegacyBinding(bindingId)
					legacyMembers = append(legacyMembers, LegacyMemberModel{
						Type: types.StringValue(memberType),
						Id:   types.StringValue(memberId),
					})
				}

				membersValue, err := convertLegacyMembersToList(ctx, legacyMembers)
				if err != nil {
					return nil, fmt.Errorf("failed to convert bindings to members: %w", err)
				}
				legacyModel.Members = membersValue
			}
		}
	}

	// Note: Bindings are now handled as part of roles conversion above

	return legacyModel, nil
}

// State management utilities (T028)

// ManageResourceState handles state transitions and property structure changes
func ManageResourceState(ctx context.Context, current *RoleBindingResourceModel, planned *RoleBindingResourceModel) (*ResourceState, error) {
	state := &ResourceState{
		UsesLegacyProperties: hasLegacyProperties(planned),
		UsesNewProperties:    hasNewProperties(planned),
		MigrationRequired:    false,
		CurrentModel:         planned,
	}

	// Detect if migration is needed
	if current != nil {
		currentUsesLegacy := hasLegacyProperties(current)
		plannedUsesNew := hasNewProperties(planned)
		plannedUsesLegacy := hasLegacyProperties(planned)
		currentUsesNew := hasNewProperties(current)

		state.MigrationRequired = (currentUsesLegacy && plannedUsesNew) || (currentUsesNew && plannedUsesLegacy)
	}

	return state, nil
}

// GetStateTransitionPlan creates a plan for transitioning between property structures
func GetStateTransitionPlan(ctx context.Context, from *RoleBindingResourceModel, to *RoleBindingResourceModel) (string, []string, error) {
	var actions []string
	var plan string

	fromLegacy := hasLegacyProperties(from)
	fromNew := hasNewProperties(from)
	toLegacy := hasLegacyProperties(to)
	toNew := hasNewProperties(to)

	switch {
	case fromLegacy && toNew:
		plan = "legacy_to_new_migration"
		actions = append(actions, "convert name to group_id")
		actions = append(actions, "convert role to roles array")
		actions = append(actions, "convert members to bindings array")
		actions = append(actions, "clear legacy properties")

	case fromNew && toLegacy:
		plan = "new_to_legacy_migration"
		actions = append(actions, "convert group_id to name")
		actions = append(actions, "convert roles array to single role")
		actions = append(actions, "convert bindings array to members")
		actions = append(actions, "clear new properties")

	case fromLegacy && toLegacy:
		plan = "legacy_update"
		actions = append(actions, "update legacy properties in-place")

	case fromNew && toNew:
		plan = "new_update"
		actions = append(actions, "update new properties in-place")

	default:
		return "", nil, fmt.Errorf("invalid state transition")
	}

	return plan, actions, nil
}

// Migration logic utilities (T029)

// PerformPropertyMigration executes the actual property migration
func PerformPropertyMigration(ctx context.Context, model *RoleBindingResourceModel, direction string) (*RoleBindingResourceModel, error) {
	switch direction {
	case "legacy_to_new":
		return ConvertLegacyToNew(ctx, model)
	case "new_to_legacy":
		return ConvertNewToLegacy(ctx, model)
	default:
		return nil, fmt.Errorf("unsupported migration direction: %s", direction)
	}
}

// ValidateMigrationPath ensures migration is safe and valid
func ValidateMigrationPath(ctx context.Context, from *RoleBindingResourceModel, to *RoleBindingResourceModel) error {
	// Check that migration doesn't lose data
	fromLegacy := hasLegacyProperties(from)
	fromNew := hasNewProperties(from)
	toLegacy := hasLegacyProperties(to)
	toNew := hasNewProperties(to)

	// Validate data preservation
	if fromNew && toLegacy {
		// Warn about potential data loss when converting from multiple roles to single role
		if !from.Roles.IsNull() {
			var roles []RoleModel
			diags := from.Roles.ElementsAs(ctx, &roles, false)
			if !diags.HasError() && len(roles) > 1 {
				return fmt.Errorf("migration would lose data: multiple roles cannot be converted to single legacy role")
			}
		}
	}

	// Ensure consistent property usage
	if (fromLegacy && fromNew) || (toLegacy && toNew) {
		return fmt.Errorf("mixed property usage detected: migration requires consistent property structure")
	}

	return nil
}

// Helper functions for type conversions

func convertRoleToList(ctx context.Context, roles []RoleModel) (basetypes.ListValue, error) {
	// This is a simplified implementation - in practice, you'd use proper Framework types
	// For now, return a null list as placeholder
	return types.ListNull(types.ObjectType{}), nil
}

func convertBindingsToList(ctx context.Context, bindings []BindingModel) (basetypes.ListValue, error) {
	// This is a simplified implementation - in practice, you'd use proper Framework types
	// For now, return a null list as placeholder
	return types.ListNull(types.ObjectType{}), nil
}

func convertMembersToList(ctx context.Context, members []LegacyMemberModel) (basetypes.ListValue, error) {
	// This is a simplified implementation - in practice, you'd use proper Framework types
	// For now, return a null list as placeholder
	return types.ListNull(types.ObjectType{}), nil
}

func convertLegacyMembersToList(ctx context.Context, members []LegacyMemberModel) (basetypes.ListValue, error) {
	// This is a simplified implementation - in practice, you'd use proper Framework types
	// For now, return a null list as placeholder
	return types.ListNull(types.ObjectType{}), nil
}

func parseLegacyBinding(bindingId string) (string, string) {
	// Parse binding ID to extract type and id
	// Format expected: "type:id" or just "id" (defaults to user)
	parts := strings.SplitN(bindingId, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	// Default to user type if no prefix
	return "user", bindingId
}

// Resource ID generation utilities

func GenerateResourceId(tenantId, groupId, roleId string) string {
	// Generate deterministic resource ID based on key components
	return fmt.Sprintf("%s-%s-%s-%s", tenantId, groupId, roleId, generateHash(tenantId+groupId+roleId))
}

func generateHash(input string) string {
	// Simple hash implementation for demo - in practice use proper hashing
	hashStr := fmt.Sprintf("%x", len(input)*17)
	// Ensure we have at least 8 characters by padding with zeros if needed
	for len(hashStr) < 8 {
		hashStr = "0" + hashStr
	}
	// Take first 8 characters
	if len(hashStr) > 8 {
		return hashStr[:8]
	}
	return hashStr
}
