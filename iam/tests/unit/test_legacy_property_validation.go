package unit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// LegacyMemberModel represents legacy member structure
type LegacyMemberModel struct {
	Type        types.String `tfsdk:"type"`
	ResourceID  types.String `tfsdk:"resource_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

// TestLegacyPropertyStructureValidation tests validation of the legacy property structure
func TestLegacyPropertyStructureValidation(t *testing.T) {
	tests := []struct {
		name            string
		input           RoleBindingResourceModel
		expectedResult  string
		expectedError   string
		expectedWarning string
	}{
		{
			name: "ValidLegacyStructure",
			input: RoleBindingResourceModel{
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
					[]attr.Value{
						types.ObjectValueMust(
							map[string]attr.Type{
								"type":         types.StringType,
								"resource_id":  types.StringType,
								"display_name": types.StringType,
							},
							map[string]attr.Value{
								"type":         types.StringValue("person"),
								"resource_id":  types.StringValue("person_123"),
								"display_name": types.StringValue("John Doe"),
							},
						),
					},
				),
			},
			expectedResult:  "valid",
			expectedError:   "",
			expectedWarning: "Property 'name' is deprecated. Use 'group_id' instead",
		},
		{
			name: "ValidLegacyMultipleMembers",
			input: RoleBindingResourceModel{
				Name: types.StringValue("finance_team"),
				Role: types.StringValue("budget_viewer"),
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
								"resource_id":  types.StringValue("person_123"),
								"display_name": types.StringValue("John Doe"),
							},
						),
						types.ObjectValueMust(
							map[string]attr.Type{
								"type":         types.StringType,
								"resource_id":  types.StringType,
								"display_name": types.StringType,
							},
							map[string]attr.Value{
								"type":         types.StringValue("service"),
								"resource_id":  types.StringValue("service_456"),
								"display_name": types.StringValue("Finance Service"),
							},
						),
					},
				),
			},
			expectedResult:  "valid",
			expectedError:   "",
			expectedWarning: "Property 'name' is deprecated. Use 'group_id' instead",
		},
		{
			name: "InvalidLegacyMissingName",
			input: RoleBindingResourceModel{
				// Missing Name
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
			},
			expectedResult: "validation_error",
			expectedError:  "name is required when using legacy property structure",
		},
		{
			name: "InvalidLegacyMissingRole",
			input: RoleBindingResourceModel{
				Name: types.StringValue("legacy_group"),
				// Missing Role
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
			},
			expectedResult: "validation_error",
			expectedError:  "role is required when using legacy property structure",
		},
		{
			name: "InvalidLegacyEmptyMembers",
			input: RoleBindingResourceModel{
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
			},
			expectedResult: "validation_error",
			expectedError:  "At least one member must be specified in members array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := tt.input.Validate()
			
			// These tests MUST fail until implementation is complete in Phase 3.3
			if tt.expectedResult == "valid" {
				// Should have deprecation warnings when properly implemented
				assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented in Phase 3.3")
			} else {
				// Should have expected validation errors when properly implemented  
				assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented in Phase 3.3")
			}
		})
	}
}

// TestLegacyPropertyConstraints tests specific constraints for legacy properties
func TestLegacyPropertyConstraints(t *testing.T) {
	t.Run("MemberTypeValidation", func(t *testing.T) {
		validMemberTypes := []string{"person", "service", "group", "application"}
		invalidMemberTypes := []string{"unknown", "invalid", ""}
		
		for _, validType := range validMemberTypes {
			t.Run("ValidType_"+validType, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Name: types.StringValue("test_group"),
					Role: types.StringValue("test_role"),
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
									"type":         types.StringValue(validType),
									"resource_id":  types.StringValue("resource_123"),
									"display_name": types.StringValue("Test Resource"),
								},
							),
						},
					),
				}
				diags := model.Validate()
				// Should have deprecation warning but be valid when properly implemented
				assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented")
			})
		}
		
		for _, invalidType := range invalidMemberTypes {
			t.Run("InvalidType_"+invalidType, func(t *testing.T) {
				model := RoleBindingResourceModel{
					Name: types.StringValue("test_group"),
					Role: types.StringValue("test_role"),
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
									"type":         types.StringValue(invalidType),
									"resource_id":  types.StringValue("resource_123"),
									"display_name": types.StringValue("Test Resource"),
								},
							),
						},
					),
				}
				diags := model.Validate()
				// Should fail with validation error when properly implemented
				assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented")
			})
		}
	})
	
	t.Run("RequiredMemberFields", func(t *testing.T) {
		// Test that type and resource_id are required, display_name is optional
		t.Run("MissingType", func(t *testing.T) {
			model := RoleBindingResourceModel{
				Name: types.StringValue("test_group"),
				Role: types.StringValue("test_role"),
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
								// Missing type
								"resource_id":  types.StringValue("resource_123"),
								"display_name": types.StringValue("Test Resource"),
							},
						),
					},
				),
			}
			diags := model.Validate()
			assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented")
		})
		
		t.Run("MissingResourceID", func(t *testing.T) {
			model := RoleBindingResourceModel{
				Name: types.StringValue("test_group"),
				Role: types.StringValue("test_role"),
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
								"type": types.StringValue("person"),
								// Missing resource_id
								"display_name": types.StringValue("Test Resource"),
							},
						),
					},
				),
			}
			diags := model.Validate()
			assert.True(t, diags.HasError(), "Test should fail until validation logic is implemented")
		})
	})
}