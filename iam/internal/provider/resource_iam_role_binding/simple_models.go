package resource_iam_role_binding

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SimpleRoleBindingResourceModel represents the simple 1:1 Group-to-Role binding model
type SimpleRoleBindingResourceModel struct {
	// Core Properties
	ID       types.String `tfsdk:"id"`
	TenantID types.String `tfsdk:"tenant_id"`

	// Required Properties
	GroupID  types.String `tfsdk:"group_id"`
	RoleID   types.String `tfsdk:"role_id"`
	IsCustom types.Bool   `tfsdk:"is_custom"`

	// Optional Properties
	Bindings    types.List   `tfsdk:"bindings"`
	Description types.String `tfsdk:"description"`
	Condition   types.String `tfsdk:"condition"`
}
