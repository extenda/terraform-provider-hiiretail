package resource_iam_role_binding

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestValidatePropertyStructure tests the new property validation logic
func TestValidatePropertyStructure(t *testing.T) {
	ctx := context.Background()

	t.Run("EmptyModel", func(t *testing.T) {
		model := &RoleBindingResourceModel{}
		result := ValidatePropertyStructure(ctx, model)

		// This should fail validation since no properties are set
		assert.False(t, result.IsValid, "Empty model should fail validation")
		assert.Equal(t, "none", result.PropertyMix, "Should detect no property structure")
		assert.Contains(t, result.Errors, "no valid property structure found")
	})

	t.Run("LegacyPropertiesOnly", func(t *testing.T) {
		model := &RoleBindingResourceModel{
			Name:    types.StringValue("test-group"),
			Role:    types.StringValue("test-role"),
			Members: types.ListValueMust(types.ObjectType{}, []attr.Value{}),
		}
		result := ValidatePropertyStructure(ctx, model)

		// This should detect legacy properties and show deprecation warning
		assert.True(t, result.IsValid, "Legacy properties should be valid")
		assert.Equal(t, "legacy", result.PropertyMix, "Should detect legacy property structure")
		assert.Contains(t, result.Warnings, "using deprecated properties - consider migrating to new structure")
	})

	t.Run("NewPropertiesOnly", func(t *testing.T) {
		model := &RoleBindingResourceModel{
			GroupId:  types.StringValue("test-group"),
			Roles:    types.ListValueMust(types.ObjectType{}, []attr.Value{}),
			Bindings: types.ListValueMust(types.ObjectType{}, []attr.Value{}),
		}
		result := ValidatePropertyStructure(ctx, model)

		// This should detect new properties
		assert.True(t, result.IsValid, "New properties should be valid")
		assert.Equal(t, "new", result.PropertyMix, "Should detect new property structure")
		assert.Empty(t, result.Warnings, "New properties should not have warnings")
	})

	t.Run("MixedProperties", func(t *testing.T) {
		model := &RoleBindingResourceModel{
			Name:    types.StringValue("test-group"),                         // Legacy
			GroupId: types.StringValue("test-group"),                         // New
			Role:    types.StringValue("test-role"),                          // Legacy
			Roles:   types.ListValueMust(types.ObjectType{}, []attr.Value{}), // New
		}
		result := ValidatePropertyStructure(ctx, model)

		// This should fail validation due to mixed properties
		assert.False(t, result.IsValid, "Mixed properties should fail validation")
		assert.Equal(t, "mixed", result.PropertyMix, "Should detect mixed property structure")
		assert.Contains(t, result.Errors, "cannot use both legacy and new property structures simultaneously")
	})
}

// TestConvertLegacyToNew tests the conversion utility
func TestConvertLegacyToNew(t *testing.T) {
	ctx := context.Background()

	t.Run("ConvertBasicLegacyModel", func(t *testing.T) {
		legacyModel := &RoleBindingResourceModel{
			Id:       types.StringValue("test-id"),
			TenantId: types.StringValue("test-tenant"),
			Name:     types.StringValue("test-group"),
			Role:     types.StringValue("test-role"),
			IsCustom: types.BoolValue(true),
		}

		newModel, err := ConvertLegacyToNew(ctx, legacyModel)

		assert.NoError(t, err, "Conversion should succeed")
		assert.NotNil(t, newModel, "New model should be returned")
		assert.Equal(t, legacyModel.Name.ValueString(), newModel.GroupId.ValueString(), "name should become group_id")
		assert.True(t, newModel.Name.IsNull(), "legacy name should be cleared")
		assert.True(t, newModel.Role.IsNull(), "legacy role should be cleared")
	})

	t.Run("ConvertModelWithoutLegacyProperties", func(t *testing.T) {
		newModel := &RoleBindingResourceModel{
			GroupId: types.StringValue("test-group"),
		}

		_, err := ConvertLegacyToNew(ctx, newModel)

		assert.Error(t, err, "Should fail when no legacy properties found")
		assert.Contains(t, err.Error(), "no legacy properties found")
	})
}
