package resource_iam_role_binding

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleBindingResourceModel represents the main resource model for hiiretail_iam_role_binding
// This model supports both legacy properties (name, role, members) and new enhanced properties
// (group_id, roles array, bindings array) for easier IAM role binding management
type RoleBindingResourceModel struct {
	// Core Properties
	Id       types.String `tfsdk:"id"`
	TenantId types.String `tfsdk:"tenant_id"`

	// Legacy Properties (deprecated but supported for backward compatibility)
	Name    types.String `tfsdk:"name"`    // Deprecated: use group_id instead
	Role    types.String `tfsdk:"role"`    // Deprecated: use roles array instead
	Members types.List   `tfsdk:"members"` // Deprecated: use bindings array instead

	// Enhanced Properties (new easier management structure)
	GroupId types.String `tfsdk:"group_id"` // Group identifier for role binding
	Roles   types.List   `tfsdk:"roles"`    // Array of RoleModel objects with nested bindings

	// Optional Properties
	Description types.String `tfsdk:"description"`
	Condition   types.String `tfsdk:"condition"`

	// Internal Properties
	IsCustom       types.Bool   `tfsdk:"is_custom"`
	RoleId         types.String `tfsdk:"role_id"`         // Legacy compatibility field
	BindingsLegacy types.List   `tfsdk:"bindings_legacy"` // Legacy compatibility field
}

// RoleModel represents a single role in the roles array with its specific bindings
// This enables multiple roles to be assigned with role-specific bindings in a single resource
type RoleModel struct {
	Id       types.String `tfsdk:"id"`       // The role identifier
	Bindings types.List   `tfsdk:"bindings"` // Array of resource IDs that get this role
}

// BindingModel represents a single binding in the legacy bindings array (deprecated)
// This is kept for backward compatibility during the transition period
type BindingModel struct {
	Type types.String `tfsdk:"type"` // "user", "group", or "service_account"
	Id   types.String `tfsdk:"id"`   // The identifier for the entity
}

// LegacyMemberModel represents backward compatibility for the legacy 'members' property
// This ensures existing configurations continue to work during migration
type LegacyMemberModel struct {
	Type types.String `tfsdk:"type"` // "user", "group", or "service_account"
	Id   types.String `tfsdk:"id"`   // The identifier for the member
}

// ResourceState represents the internal state tracking for the resource
// This helps with state management and migration between property structures
type ResourceState struct {
	UsesLegacyProperties bool                      // Track if legacy properties are in use
	UsesNewProperties    bool                      // Track if new properties are in use
	MigrationRequired    bool                      // Whether migration is needed
	CurrentModel         *RoleBindingResourceModel // Current resource state
}

// ConversionContext provides context for property conversion operations
// This assists with converting between legacy and new property structures
type ConversionContext struct {
	TenantId       string // Tenant context for validation
	AllowMixed     bool   // Whether mixed property usage is allowed
	PreferNew      bool   // Prefer new properties during conversion
	ValidationMode bool   // Whether to run in validation-only mode
}

// ValidationResult contains the results of model validation
// This provides detailed feedback for validation operations
type ValidationResult struct {
	IsValid        bool     // Overall validation result
	Errors         []string // List of validation errors
	Warnings       []string // List of validation warnings
	PropertyMix    string   // Description of property structure used
	MigrationHints []string // Suggestions for property migration
}
