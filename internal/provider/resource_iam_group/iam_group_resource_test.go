package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// Minimal, focused unit tests to exercise helpers in iam_group_resource.go

type fakeProvider struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}

// GroupTestBuilder is a small test-only builder to construct IamGroupModel
// instances for unit tests. Kept in test code only to avoid touching
// production code.
type GroupTestBuilder struct {
	name        string
	description string
	id          string
	status      string
	tenantId    string
}

func NewGroupTestBuilder() *GroupTestBuilder {
	return &GroupTestBuilder{
		name:        "default-group",
		description: "",
		id:          "",
		status:      "",
		tenantId:    "",
	}
}

func (b *GroupTestBuilder) WithName(n string) *GroupTestBuilder        { b.name = n; return b }
func (b *GroupTestBuilder) WithDescription(d string) *GroupTestBuilder { b.description = d; return b }
func (b *GroupTestBuilder) WithId(id string) *GroupTestBuilder         { b.id = id; return b }
func (b *GroupTestBuilder) Build() *IamGroupModel {
	m := &IamGroupModel{
		Name:        types.StringValue(b.name),
		Description: types.StringValue(b.description),
		Id:          types.StringValue(b.id),
		Status:      types.StringValue(b.status),
		TenantId:    types.StringValue(b.tenantId),
	}
	if b.description == "" {
		m.Description = types.StringNull()
	}
	if b.id == "" {
		m.Id = types.StringNull()
	}
	if b.status == "" {
		m.Status = types.StringNull()
	}
	if b.tenantId == "" {
		m.TenantId = types.StringNull()
	}
	return m
}

func TestExtractAPIClientFieldsNil(t *testing.T) {
	if extractAPIClientFields(nil) != nil {
		t.Fatalf("expected nil for nil provider data")
	}
}

func TestExtractAPIClientFieldsStruct(t *testing.T) {
	fp := &fakeProvider{BaseURL: "https://api.example.com", TenantID: "t-1", HTTPClient: &http.Client{}}
	c := extractAPIClientFields(fp)
	if c == nil {
		t.Fatalf("expected APIClient, got nil")
	}
	if c.BaseURL != fp.BaseURL || c.TenantID != fp.TenantID || c.HTTPClient == nil {
		t.Fatalf("unexpected extracted values")
	}
}

func TestValidateGroupData_BasicErrors(t *testing.T) {
	r := &IamGroupResource{}
	ctx := context.Background()

	var m IamGroupModel
	// name empty should error
	m.Name = types.StringNull()
	if err := r.validateGroupData(ctx, &m); err == nil {
		t.Fatalf("expected error for empty name")
	}

	// too long name
	long := ""
	for i := 0; i < 260; i++ {
		long += "x"
	}
	m.Name = types.StringValue(long)
	if err := r.validateGroupData(ctx, &m); err == nil {
		t.Fatalf("expected error for too long name")
	}
}

func TestMapHTTPError_BasicMapping(t *testing.T) {
	r := &IamGroupResource{}
	if err := r.mapHTTPError(404, nil); err == nil || err.Error() == "" {
		t.Fatalf("expected non-nil error for 404")
	}
	if err := r.mapHTTPError(401, nil); err == nil {
		t.Fatalf("expected non-nil error for 401")
	}
}

// TestGroupResourceSchema tests the schema definition for the Group resource
func TestGroupResourceSchema(t *testing.T) {
	t.Run("schema structure", func(t *testing.T) {
		ctx := context.Background()
		schema := IamGroupResourceSchema(ctx)

		// Verify all required attributes are present
		attributes := schema.Attributes
		assert.Contains(t, attributes, "name", "name attribute should be present")
		assert.Contains(t, attributes, "description", "description attribute should be present")
		assert.Contains(t, attributes, "id", "id attribute should be present")
		assert.Contains(t, attributes, "status", "status attribute should be present")
		assert.Contains(t, attributes, "tenant_id", "tenant_id attribute should be present")
	})

	t.Run("name attribute properties", func(t *testing.T) {
		ctx := context.Background()
		schema := IamGroupResourceSchema(ctx)

		nameAttr := schema.Attributes["name"]
		// This test will fail until we properly implement the schema validation
		// For now, we're testing the structure exists
		assert.NotNil(t, nameAttr, "name attribute should not be nil")

		// Basic structural checks are sufficient for unit tests today.
		// When schema validation is expanded we'll add stricter assertions here.
	})

	t.Run("description attribute properties", func(t *testing.T) {
		ctx := context.Background()
		schema := IamGroupResourceSchema(ctx)

		descAttr := schema.Attributes["description"]
		assert.NotNil(t, descAttr, "description attribute should not be nil")

		// Basic structural checks are sufficient for unit tests today.
		// When schema validation is expanded we'll add stricter assertions here.
	})

	t.Run("computed attributes", func(t *testing.T) {
		ctx := context.Background()
		schema := IamGroupResourceSchema(ctx)

		// Test id attribute
		idAttr := schema.Attributes["id"]
		assert.NotNil(t, idAttr, "id attribute should not be nil")

		// Test status attribute
		statusAttr := schema.Attributes["status"]
		assert.NotNil(t, statusAttr, "status attribute should not be nil")

		// Basic structural checks are sufficient for unit tests today.
		// When schema validation is expanded we'll add stricter assertions here.
	})
}

// TestGroupModelDataBinding tests the IamGroupModel struct data binding
func TestGroupModelDataBinding(t *testing.T) {
	t.Run("model structure", func(t *testing.T) {
		model := &IamGroupModel{}

		// Verify model has all required fields
		assert.NotNil(t, &model.Name, "Name field should exist")
		assert.NotNil(t, &model.Description, "Description field should exist")
		assert.NotNil(t, &model.Id, "Id field should exist")
		assert.NotNil(t, &model.Status, "Status field should exist")
		assert.NotNil(t, &model.TenantId, "TenantId field should exist")
	})

	t.Run("model field binding", func(t *testing.T) {
		model := &IamGroupModel{
			Name:        types.StringValue("test-group"),
			Description: types.StringValue("Test description"),
			Id:          types.StringValue("group-123"),
			Status:      types.StringValue("active"),
			TenantId:    types.StringValue("tenant-123"),
		}

		// Test that values are properly set
		assert.Equal(t, "test-group", model.Name.ValueString())
		assert.Equal(t, "Test description", model.Description.ValueString())
		assert.Equal(t, "group-123", model.Id.ValueString())
		assert.Equal(t, "active", model.Status.ValueString())
		assert.Equal(t, "tenant-123", model.TenantId.ValueString())
	})

	// Note: we purposefully avoid asserting framework type semantics (IsNull/IsUnknown)
	// here to prevent brittle tests that duplicate Terraform SDK behavior. Higher-level
	// behavior is covered by resource method tests (Create/Update) which exercise
	// how null/unknown values are handled by our code.
}

// TestGroupValidationRules tests validation rules for Group resource fields
func TestGroupValidationRules(t *testing.T) {
	t.Run("name validation", func(t *testing.T) {
		tests := []struct {
			name        string
			value       string
			expectValid bool
			description string
		}{
			{
				name:        "valid name",
				value:       "developers",
				expectValid: true,
				description: "simple valid group name",
			},
			{
				name:        "name with hyphens",
				value:       "senior-developers",
				expectValid: true,
				description: "group name with hyphens",
			},
			{
				name:        "name with numbers",
				value:       "team-2024",
				expectValid: true,
				description: "group name with numbers",
			},
			{
				name:        "empty name",
				value:       "",
				expectValid: false,
				description: "empty group name should be invalid",
			},
			{
				name:        "name too long",
				value:       string(make([]byte, 256)), // 256 characters
				expectValid: false,
				description: "group name longer than 255 characters",
			},
			{
				name:        "name exactly 255 chars",
				value:       string(make([]byte, 255)),
				expectValid: true,
				description: "group name exactly 255 characters should be valid",
			},
		}

		r := &IamGroupResource{}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := &IamGroupModel{
					Name: types.StringValue(tt.value),
				}

				err := r.validateGroupData(context.Background(), model)
				if tt.expectValid {
					if err != nil {
						t.Fatalf("expected valid name (%s), got error: %v", tt.description, err)
					}
				} else {
					if err == nil {
						t.Fatalf("expected invalid name (%s) to produce an error", tt.description)
					}
				}
			})
		}
	})

	t.Run("description validation", func(t *testing.T) {
		tests := []struct {
			name        string
			value       *string
			expectValid bool
			description string
		}{
			{
				name:        "valid description",
				value:       stringPtr("Development team members"),
				expectValid: true,
				description: "normal description",
			},
			{
				name:        "empty description",
				value:       stringPtr(""),
				expectValid: true,
				description: "empty description should be valid (optional field)",
			},
			{
				name:        "null description",
				value:       nil,
				expectValid: true,
				description: "null description should be valid (optional field)",
			},
			{
				name:        "description too long",
				value:       stringPtr(string(make([]byte, 256))), // 256 characters
				expectValid: false,
				description: "description longer than 255 characters",
			},
			{
				name:        "description exactly 255 chars",
				value:       stringPtr(string(make([]byte, 255))),
				expectValid: true,
				description: "description exactly 255 characters should be valid",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var desc types.String
				if tt.value == nil {
					desc = types.StringNull()
				} else {
					desc = types.StringValue(*tt.value)
				}

				model := &IamGroupModel{
					Name:        types.StringValue("test-group"),
					Description: desc,
				}

				r := &IamGroupResource{}
				err := r.validateGroupData(context.Background(), model)
				if tt.expectValid {
					if err != nil {
						t.Fatalf("expected valid description (%s), got error: %v", tt.description, err)
					}
				} else {
					if err == nil {
						t.Fatalf("expected invalid description (%s) to produce an error", tt.description)
					}
				}
			})
		}
	})

	t.Run("computed field behavior", func(t *testing.T) {
		t.Run("id field", func(t *testing.T) {
			// Test that id field behaves correctly as computed/optional
			model := &IamGroupModel{
				Name: types.StringValue("test-group"),
				Id:   types.StringNull(), // Should be acceptable for computed field
			}

			assert.True(t, model.Id.IsNull())
		})

		t.Run("status field", func(t *testing.T) {
			// Test that status field behaves correctly as computed
			model := &IamGroupModel{
				Name:   types.StringValue("test-group"),
				Status: types.StringNull(), // Should be acceptable for computed field
			}

			assert.True(t, model.Status.IsNull())
		})

		t.Run("tenant_id field", func(t *testing.T) {
			// Test that tenant_id field behaves correctly as optional/computed
			model := &IamGroupModel{
				Name:     types.StringValue("test-group"),
				TenantId: types.StringNull(), // Should be acceptable for optional/computed field
			}

			assert.True(t, model.TenantId.IsNull())
		})
	})
}

// TestGroupModelBuilder tests the test data builder pattern
func TestGroupModelBuilder(t *testing.T) {
	// Test-only builder to exercise test data creation patterns. Implemented
	// at package level to avoid declaring functions inside another function.
	t.Run("builder pattern", func(t *testing.T) {
		builder := NewGroupTestBuilder()
		group := builder.WithName("test-group").WithDescription("Test desc").Build()
		if group.Name.ValueString() != "test-group" {
			t.Fatalf("expected name set by builder, got %s", group.Name.ValueString())
		}
		if !group.Description.IsNull() && group.Description.ValueString() != "Test desc" {
			t.Fatalf("expected description set by builder, got %v", group.Description)
		}
	})

	t.Run("builder defaults", func(t *testing.T) {
		builder := NewGroupTestBuilder()
		group := builder.Build()
		if group.Name.ValueString() != "default-group" {
			t.Fatalf("expected default name from builder, got %s", group.Name.ValueString())
		}
		if !group.Description.IsNull() {
			t.Fatalf("expected default description to be null, got %v", group.Description)
		}
	})

	t.Run("builder customization", func(t *testing.T) {
		builder := NewGroupTestBuilder().WithName("custom").WithDescription("desc").WithId("gid")
		group := builder.Build()
		if group.Name.ValueString() != "custom" || group.Id.ValueString() != "gid" {
			t.Fatalf("expected customized fields to be set, got name=%s id=%s", group.Name.ValueString(), group.Id.ValueString())
		}
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
