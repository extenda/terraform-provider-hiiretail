package resource_iam_role_binding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SimpleIamRoleBindingResourceSchema provides a simple 1:1 Group-to-Role relationship schema
func SimpleIamRoleBindingResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages IAM role bindings with a simple 1:1 relationship between Group and Role.\n\n" +
			"**Properties:**\n" +
			"- `group_id`: Group identifier for role binding (required)\n" +
			"- `role_id`: Role identifier to bind to the group (required)\n" +
			"- `is_custom`: Whether the role is a custom role or built-in role (required)\n" +
			"- `bindings`: Array of resource IDs that should receive this role (optional, defaults to empty array)\n\n" +
			"**Note:** For multiple roles on the same group, create multiple `iam_role_binding` resources.",

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

			// Required Properties
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The group identifier for the role binding",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The role identifier to bind to the group",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"is_custom": schema.BoolAttribute{
				MarkdownDescription: "Whether this role is a custom role (true) or built-in role (false)",
				Required:            true,
			},

			// Optional Properties
			"bindings": schema.ListAttribute{
				MarkdownDescription: "Array of resource IDs that should receive this role. If empty, no specific resource bindings are applied.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(20),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the role binding",
				Optional:            true,
			},
			"condition": schema.StringAttribute{
				MarkdownDescription: "Optional condition expression for conditional role binding",
				Optional:            true,
			},
		},
	}
}
