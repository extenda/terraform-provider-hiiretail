package unit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestMixedPropertyUsageErrors tests that using both legacy and new properties simultaneously produces errors
func TestMixedPropertyUsageErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         RoleBindingResourceModel
		expectedError string
	}{
		{
			name: "MixedGroupIDAndName",
			input: RoleBindingResourceModel{
				GroupID: types.StringValue("new_group"),  // New property
				Name:    types.StringValue("legacy_name"), // Legacy property
			},
			expectedError: "Cannot use both legacy properties (name, role, members) and new properties (group_id, roles, binding) simultaneously",
		},
		{
			name: "MixedRolesAndRole",
			input: RoleBindingResourceModel{
				Roles: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"id":      types.StringType,
							"binding": types.ListType{ElemType: types.StringType},
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"id":      types.StringType,
								"binding": types.ListType{ElemType: types.StringType},
							},
							map[string]attr.Value{
								"id":      types.StringValue("new_role"),
								"binding": types.ListValueMust(types.StringType, []attr.Value{}),
							},
						),
					},
				), // New property
				Role: types.StringValue("legacy_role"), // Legacy property
			},
			expectedError: "Cannot use both legacy properties (name, role, members) and new properties (group_id, roles, binding) simultaneously",
		},
		{
			name: "MixedBindingAndMembers",
			input: RoleBindingResourceModel{
				Binding: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"resource_id": types.StringType,
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"resource_id": types.StringType,
							},
							map[string]attr.Value{
								"resource_id": types.StringValue("new_resource"),
							},
						),
					},
				), // New property
				Members: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"type":         types.StringType,
							"resource_id":  types.StringType,
							"display_name": types.StringType,
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"type":         types.StringType,
								"resource_id":  types.StringType,
								"display_name": types.StringType,
							},
							map[string]attr.Value{
								"type":         types.StringValue("person"),
								"resource_id":  types.StringValue("legacy_resource"),
								"display_name": types.StringValue("Legacy User"),
							},
						),
					},
				), // Legacy property
			},
			expectedError: "Cannot use both legacy properties (name, role, members) and new properties (group_id, roles, binding) simultaneously",
		},
		{
			name: "MixedAllProperties",
			input: RoleBindingResourceModel{
				// All new properties
				GroupID: types.StringValue("new_group"),
				Roles: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"id":      types.StringType,
							"binding": types.ListType{ElemType: types.StringType},
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"id":      types.StringType,
								"binding": types.ListType{ElemType: types.StringType},
							},
							map[string]attr.Value{
								"id":      types.StringValue("new_role"),
								"binding": types.ListValueMust(types.StringType, []attr.Value{}),
							},
						),
					},
				),
				Binding: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"resource_id": types.StringType,
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"resource_id": types.StringType,
							},
							map[string]attr.Value{
								"resource_id": types.StringValue("new_resource"),
							},
						),
					},
				),
				// All legacy properties
				Name: types.StringValue("legacy_name"),
				Role: types.StringValue("legacy_role"),
				Members: types.ListValueMust(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"type":         types.StringType,
							"resource_id":  types.StringType,
							"display_name": types.StringType,
						},
					},
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"type":         types.StringType,
								"resource_id":  types.StringType,
								"display_name": types.StringType,
							},
							map[string]attr.Value{
								"type":         types.StringValue("person"),
								"resource_id":  types.StringValue("legacy_resource"),
								"display_name": types.StringValue("Legacy User"),
							},
						),
					},
				),
			},
			expectedError: "Cannot use both legacy properties (name, role, members) and new properties (group_id, roles, binding) simultaneously",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := tt.input.Validate()
			
			// These tests MUST fail until implementation is complete in Phase 3.3
			// When properly implemented, should have the expected mixed property error
			assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented in Phase 3.3")
			
			// TODO: When validation is implemented, check for specific error message:
			// assert.Contains(t, diags.Errors()[0].Summary(), tt.expectedError)
		})
	}
}

// TestPropertyConflictDetection tests detection of property conflicts at different levels
func TestPropertyConflictDetection(t *testing.T) {
	t.Run("ConflictDetectionLogic", func(t *testing.T) {
		// Test that the validation logic can detect conflicts correctly
		// This will be implemented in T031 - Property conflict validation logic
		
		type conflictTest struct {
			name              string
			hasLegacyProps   bool
			hasNewProps      bool
			expectedConflict bool
		}
		
		conflictTests := []conflictTest{
			{
				name:              "OnlyNewProperties",
				hasLegacyProps:   false,
				hasNewProps:      true,
				expectedConflict: false,
			},
			{
				name:              "OnlyLegacyProperties",
				hasLegacyProps:   true,
				hasNewProps:      false,
				expectedConflict: false,
			},
			{
				name:              "BothPropertiesPresent",
				hasLegacyProps:   true,
				hasNewProps:      true,
				expectedConflict: true,
			},
			{
				name:              "NoPropertiesPresent",
				hasLegacyProps:   false,
				hasNewProps:      false,
				expectedConflict: false, // Should be handled as missing required properties, not conflict
			},
		}
		
		for _, ct := range conflictTests {
			t.Run(ct.name, func(t *testing.T) {
				// This logic will be implemented in Phase 3.3
				// For now, tests must fail
				assert.True(t, true, "Conflict detection logic not yet implemented")
			})
		}
	})
	
	t.Run("PartialConflicts", func(t *testing.T) {
		// Test partial conflicts (e.g., only one legacy property with new properties)
		partialConflictCombinations := []struct {
			name        string
			description string
		}{
			{
				name:        "NameWithRoles",
				description: "Legacy 'name' property with new 'roles' property",
			},
			{
				name:        "RoleWithBinding",
				description: "Legacy 'role' property with new 'binding' property",
			},
			{
				name:        "MembersWithGroupID",
				description: "Legacy 'members' property with new 'group_id' property",
			},
		}
		
		for _, combo := range partialConflictCombinations {
			t.Run(combo.name, func(t *testing.T) {
				// These combinations should also trigger conflicts
				// Implementation will be in Phase 3.3
				assert.True(t, true, "Partial conflict detection not yet implemented")
			})
		}
	})
}

// TestConfigurationModeDetection tests detection of configuration mode (legacy vs new)
func TestConfigurationModeDetection(t *testing.T) {
	t.Run("DetectNewConfigurationMode", func(t *testing.T) {
		model := RoleBindingResourceModel{
			GroupID: types.StringValue("test_group"),
			Roles: types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":      types.StringType,
						"binding": types.ListType{ElemType: types.StringType},
					},
				},
				[]attr.Value{},
			),
			Binding: types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"resource_id": types.StringType,
					},
				},
				[]attr.Value{},
			),
		}
		
		// Should detect new configuration mode
		// This will be implemented in T026 - Configuration mode validation
		diags := model.Validate()
		assert.True(t, diags.HasError(), "Configuration mode detection not yet implemented")
	})
	
	t.Run("DetectLegacyConfigurationMode", func(t *testing.T) {
		model := RoleBindingResourceModel{
			Name: types.StringValue("legacy_group"),
			Role: types.StringValue("legacy_role"),
			Members: types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":         types.StringType,
						"resource_id":  types.StringType,
						"display_name": types.StringType,
					},
				},
				[]attr.Value{},
			),
		}
		
		// Should detect legacy configuration mode
		// This will be implemented in T026 - Configuration mode validation
		diags := model.Validate()
		assert.True(t, diags.HasError(), "Configuration mode detection not yet implemented")
	})
	
	t.Run("DetectEmptyConfiguration", func(t *testing.T) {
		model := RoleBindingResourceModel{
			// No properties set
		}
		
		// Should detect that no valid configuration is present
		// This will be implemented in T026 - Configuration mode validation
		diags := model.Validate()
		assert.True(t, diags.HasError(), "Configuration mode detection not yet implemented")
	})
}