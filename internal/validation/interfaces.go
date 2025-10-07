package validation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FieldValidator defines the interface for field-level validation
type FieldValidator interface {
	// ValidateString validates a string field value
	ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse)

	// Description returns a human-readable description of the validator
	Description(ctx context.Context) string

	// MarkdownDescription returns a markdown-formatted description of the validator
	MarkdownDescription(ctx context.Context) string
}

// ResourceValidator defines the interface for resource-level validation
type ResourceValidator interface {
	// ValidateResource validates an entire resource configuration
	ValidateResource(ctx context.Context, config map[string]interface{}) *ValidationResult

	// GetResourceType returns the resource type this validator handles
	GetResourceType() string
}

// ReferenceValidator defines the interface for cross-resource reference validation
type ReferenceValidator interface {
	// ValidateReference validates a reference to another resource
	ValidateReference(ctx context.Context, referenceType, referenceValue string) *ValidationResult

	// ResolveReference attempts to resolve a reference to its actual resource ID
	ResolveReference(ctx context.Context, referenceType, referenceValue string) (string, error)

	// GetSuggestions returns suggestions for similar resources when validation fails
	GetSuggestions(ctx context.Context, referenceType, referenceValue string) []string
}

// PermissionValidator defines the interface for permission string validation
type PermissionValidator interface {
	// ValidatePermission validates a permission string format and existence
	ValidatePermission(ctx context.Context, permission, service, context string) *ValidationResult

	// NormalizePermission returns the normalized form of a permission
	NormalizePermission(permission string) string

	// GetPermissionCategory returns the category of a permission (read, write, admin, etc.)
	GetPermissionCategory(permission string) string

	// GetRelatedPermissions returns related permissions for suggestions
	GetRelatedPermissions(permission string) []string
}

// APIValidator defines the interface for API connectivity validation
type APIValidator interface {
	// ValidateConnectivity checks if the API is accessible
	ValidateConnectivity(ctx context.Context) error

	// ValidateCredentials checks if the credentials are valid
	ValidateCredentials(ctx context.Context) error

	// GetAPIEndpoint returns the API endpoint being validated
	GetAPIEndpoint() string
}

// ValidationContext provides context for validation operations
type ValidationContext struct {
	// ResourceConfig contains the current resource configuration
	ResourceConfig map[string]interface{}

	// ProviderConfig contains the provider configuration
	ProviderConfig map[string]interface{}

	// ExistingState contains the existing resource state (if any)
	ExistingState map[string]interface{}

	// PlanningPhase indicates if validation is happening during planning
	PlanningPhase bool

	// TenantID is the current tenant identifier
	TenantID string

	// APIClient can be used for API calls during validation
	APIClient interface{}
}

// ValidationRequest represents a validation request
type ValidationRequest struct {
	ResourceType string
	FieldPath    string
	Value        interface{}
	Context      *ValidationContext
}

// ValidationResponse represents a validation response
type ValidationResponse struct {
	Valid       bool
	ErrorCode   string
	Message     string
	Suggestions []string
	Examples    []string
	Severity    string
}

// ValidatorRegistry manages registered validators
type ValidatorRegistry struct {
	fieldValidators     map[string]map[string]FieldValidator
	resourceValidators  map[string]ResourceValidator
	referenceValidators map[string]ReferenceValidator
	permissionValidator PermissionValidator
	apiValidator        APIValidator
}

// NewValidatorRegistry creates a new validator registry
func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		fieldValidators:     make(map[string]map[string]FieldValidator),
		resourceValidators:  make(map[string]ResourceValidator),
		referenceValidators: make(map[string]ReferenceValidator),
	}
}

// RegisterFieldValidator registers a field validator for a specific resource type and field
func (r *ValidatorRegistry) RegisterFieldValidator(resourceType, fieldName string, validator FieldValidator) {
	if r.fieldValidators[resourceType] == nil {
		r.fieldValidators[resourceType] = make(map[string]FieldValidator)
	}
	r.fieldValidators[resourceType][fieldName] = validator
}

// RegisterResourceValidator registers a resource validator
func (r *ValidatorRegistry) RegisterResourceValidator(validator ResourceValidator) {
	r.resourceValidators[validator.GetResourceType()] = validator
}

// RegisterReferenceValidator registers a reference validator
func (r *ValidatorRegistry) RegisterReferenceValidator(referenceType string, validator ReferenceValidator) {
	r.referenceValidators[referenceType] = validator
}

// RegisterPermissionValidator registers the permission validator
func (r *ValidatorRegistry) RegisterPermissionValidator(validator PermissionValidator) {
	r.permissionValidator = validator
}

// RegisterAPIValidator registers the API validator
func (r *ValidatorRegistry) RegisterAPIValidator(validator APIValidator) {
	r.apiValidator = validator
}

// GetFieldValidator retrieves a field validator
func (r *ValidatorRegistry) GetFieldValidator(resourceType, fieldName string) (FieldValidator, bool) {
	if fields, ok := r.fieldValidators[resourceType]; ok {
		validator, exists := fields[fieldName]
		return validator, exists
	}
	return nil, false
}

// GetResourceValidator retrieves a resource validator
func (r *ValidatorRegistry) GetResourceValidator(resourceType string) (ResourceValidator, bool) {
	validator, exists := r.resourceValidators[resourceType]
	return validator, exists
}

// GetReferenceValidator retrieves a reference validator
func (r *ValidatorRegistry) GetReferenceValidator(referenceType string) (ReferenceValidator, bool) {
	validator, exists := r.referenceValidators[referenceType]
	return validator, exists
}

// GetPermissionValidator retrieves the permission validator
func (r *ValidatorRegistry) GetPermissionValidator() PermissionValidator {
	return r.permissionValidator
}

// GetAPIValidator retrieves the API validator
func (r *ValidatorRegistry) GetAPIValidator() APIValidator {
	return r.apiValidator
}

// ValidateField validates a single field using registered validators
func (r *ValidatorRegistry) ValidateField(ctx context.Context, resourceType, fieldName string, value interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if _, exists := r.GetFieldValidator(resourceType, fieldName); exists {
		// Create a string request for the validator
		// Note: This is simplified - real implementation would handle different value types
		// TODO: Fix validator interface compatibility
		// req := validator.StringRequest{
		//	ConfigValue: value,
		//	Path:        nil, // Would be set appropriately
		// }

		// resp := &validator.StringResponse{}
		// validator.ValidateString(ctx, req, resp)

		// diags.Append(resp.Diagnostics...)
	}

	return diags
}

// ValidateResource validates an entire resource using registered validators
func (r *ValidatorRegistry) ValidateResource(ctx context.Context, resourceType string, config map[string]interface{}) *ValidationResult {
	if validator, exists := r.GetResourceValidator(resourceType); exists {
		return validator.ValidateResource(ctx, config)
	}

	// Return valid result if no validator is registered
	return NewValidationResult()
}

// DefaultRegistry is the global validator registry instance
var DefaultRegistry = NewValidatorRegistry()
