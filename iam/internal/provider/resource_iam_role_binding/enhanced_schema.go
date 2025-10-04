package resource_iam_role_binding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// EnhancedIamRoleBindingResourceSchema provides the enhanced schema supporting both legacy and new properties
// T030-T032: Enhanced schema with dual property support and deprecation warnings
func EnhancedIamRoleBindingResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages IAM role bindings with support for both legacy and enhanced property structures.\n\n" +
			"**Enhanced Structure (Recommended):**\n" +
			"- `group_id`: Group identifier for role binding\n" +
			"- `roles`: Array of role objects with `role_id` and `is_custom` fields\n" +
			"- `bindings`: Array of binding objects with `type` and `id` fields\n\n" +
			"**Legacy Structure (Deprecated):**\n" +
			"- `name`: Group name (deprecated, use `group_id` instead)\n" +
			"- `role`: Single role ID (deprecated, use `roles` array instead)\n" +
			"- `members`: Array of member strings (deprecated, use `bindings` array instead)\n\n" +
			"**Note:** Cannot mix legacy and enhanced properties in the same resource configuration.",

		Attributes: map[string]schema.Attribute{
			// Core Properties
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the role binding resource",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID for the role binding",
				Optional:            true,
				Computed:            true,
			},

			// Legacy Properties (Deprecated but supported)
			"name": schema.StringAttribute{
				MarkdownDescription: "**Deprecated:** Use `group_id` instead. The name/identifier of the group for the role binding.",
				Optional:            true,
				DeprecationMessage:  "The 'name' attribute is deprecated. Use 'group_id' instead for the enhanced property structure.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("group_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "**Deprecated:** Use `roles` array instead. The single role ID to bind.",
				Optional:            true,
				DeprecationMessage:  "The 'role' attribute is deprecated. Use the 'roles' array instead for the enhanced property structure.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("roles")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"members": schema.ListAttribute{
				MarkdownDescription: "**Deprecated:** Use `bindings` array instead. List of member identifiers in format 'type:id'.",
				ElementType:         types.StringType,
				Optional:            true,
				DeprecationMessage:  "The 'members' attribute is deprecated. Use the 'bindings' array instead for the enhanced property structure.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("bindings")),
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(10),
				},
			},

			// Enhanced Properties (New structure)
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The group identifier for the role binding. Use this instead of the deprecated 'name' attribute.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"roles": schema.ListNestedAttribute{
				MarkdownDescription: "Array of roles with their specific bindings. Each role contains the role_id and its associated bindings.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("role")),
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(10),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The identifier of the role to bind",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"bindings": schema.ListAttribute{
							MarkdownDescription: "Array of resource IDs that should receive this role",
							ElementType:         types.StringType,
							Required:            true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(20),
							},
						},
					},
				},
			},

			// Optional Properties (compatible with both structures)
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the role binding",
				Optional:            true,
			},
			"condition": schema.StringAttribute{
				MarkdownDescription: "Optional condition expression for conditional role binding",
				Optional:            true,
			},

			// Backward compatibility support
			"is_custom": schema.BoolAttribute{
				MarkdownDescription: "**Legacy compatibility:** Whether the role is custom. This is automatically determined from the roles array in the enhanced structure.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},

			// Legacy bindings compatibility (the old simple string array)
			"bindings_legacy": schema.ListAttribute{
				MarkdownDescription: "**Internal use only:** Legacy bindings format for backward compatibility",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},

			// Enhanced role ID for backward compatibility
			"role_id": schema.StringAttribute{
				MarkdownDescription: "**Legacy compatibility:** The role ID. This is automatically determined from the roles array in the enhanced structure.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// GetRoleModelObjectType returns the object type for RoleModel
func GetRoleModelObjectType() basetypes.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_id":   types.StringType,
			"is_custom": types.BoolType,
		},
	}
}

// GetBindingModelObjectType returns the object type for BindingModel
func GetBindingModelObjectType() basetypes.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
		},
	}
}

// GetLegacyMemberModelObjectType returns the object type for LegacyMemberModel
func GetLegacyMemberModelObjectType() basetypes.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
		},
	}
}
