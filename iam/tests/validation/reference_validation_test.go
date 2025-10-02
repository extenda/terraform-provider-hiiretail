package validation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReferenceValidationContract tests the reference validation contracts
// These tests MUST FAIL initially - they test against unimplemented validators
func TestReferenceValidationContract(t *testing.T) {

	t.Run("GroupReferenceValidation", func(t *testing.T) {
		testCases := []ReferenceTestCase{
			{
				Name:           "ValidGroupReference",
				ReferenceType:  "group",
				ReferenceValue: "test-group-unique-id",
				ExpectedValid:  true,
			},
			{
				Name:                "NonExistentGroupReference",
				ReferenceType:       "group",
				ReferenceValue:      "nonexistent-group",
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"test-group-unique-id"},
			},
			{
				Name:                "TypoInGroupReference",
				ReferenceType:       "group",
				ReferenceValue:      "test-grup-unique-id", // typo: grup instead of group
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"test-group-unique-id"},
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the reference validator doesn't exist yet
		if validator, exists := registry.GetReferenceValidator("group"); !exists {
			t.Logf("Expected group reference validator not found - this is expected for TDD")
			assert.Fail(t, "Group reference validator not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestReferenceValidationContract(validator, testCases)
		}
	})

	t.Run("RoleReferenceValidation", func(t *testing.T) {
		testCases := []ReferenceTestCase{
			{
				Name:           "ValidCustomRoleReference",
				ReferenceType:  "role",
				ReferenceValue: "roles/custom.test-custom-role-unique-id",
				ExpectedValid:  true,
			},
			{
				Name:           "ValidBuiltinRoleReference",
				ReferenceType:  "role",
				ReferenceValue: "roles/viewer",
				ExpectedValid:  true,
			},
			{
				Name:                "InvalidRoleFormat",
				ReferenceType:       "role",
				ReferenceValue:      "custom.test-role", // missing "roles/" prefix
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"roles/custom.test-role"},
			},
			{
				Name:                "NonExistentCustomRole",
				ReferenceType:       "role",
				ReferenceValue:      "roles/custom.nonexistent-role",
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"roles/custom.test-custom-role-unique-id"},
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the reference validator doesn't exist yet
		if validator, exists := registry.GetReferenceValidator("role"); !exists {
			t.Logf("Expected role reference validator not found - this is expected for TDD")
			assert.Fail(t, "Role reference validator not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestReferenceValidationContract(validator, testCases)
		}
	})

	t.Run("MemberReferenceValidation", func(t *testing.T) {
		testCases := []ReferenceTestCase{
			{
				Name:           "ValidGroupMemberReference",
				ReferenceType:  "member",
				ReferenceValue: "group:test-group-unique-id",
				ExpectedValid:  true,
			},
			{
				Name:                "InvalidMemberFormat",
				ReferenceType:       "member",
				ReferenceValue:      "test-group-unique-id", // missing "group:" prefix
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"group:test-group-unique-id"},
			},
			{
				Name:           "UnsupportedMemberType",
				ReferenceType:  "member",
				ReferenceValue: "user:test@example.com", // users not supported according to simple_test.tf
				ExpectedValid:  false,
			},
			{
				Name:                "NonExistentGroupMember",
				ReferenceType:       "member",
				ReferenceValue:      "group:nonexistent-group",
				ExpectedValid:       false,
				ExpectedSuggestions: []string{"group:test-group-unique-id"},
			},
		}

		ctx := context.Background()
		registry := NewValidatorRegistry()

		// This should fail since the reference validator doesn't exist yet
		if validator, exists := registry.GetReferenceValidator("member"); !exists {
			t.Logf("Expected member reference validator not found - this is expected for TDD")
			assert.Fail(t, "Member reference validator not implemented yet - implement in Phase 3.3")
		} else {
			runner := NewContractTestRunner(t)
			runner.TestReferenceValidationContract(validator, testCases)
		}
	})
}

// TestReferenceResolutionIntegration tests reference resolution integration
func TestReferenceResolutionIntegration(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()

	t.Run("ResolveGroupReference", func(t *testing.T) {
		// This should fail since no resolver is implemented yet
		if validator, exists := registry.GetReferenceValidator("group"); !exists {
			t.Logf("Group reference validator not found - expected for TDD")
			assert.Fail(t, "Group reference resolution not implemented yet")
		} else {
			resolvedID, err := validator.ResolveReference(ctx, "group", "test-group-unique-id")
			if err != nil {
				t.Logf("Reference resolution failed as expected: %v", err)
				assert.Fail(t, "Group reference resolution not implemented yet")
			} else {
				assert.NotEmpty(t, resolvedID, "Resolved ID should not be empty")
			}
		}
	})

	t.Run("ResolveRoleReference", func(t *testing.T) {
		// This should fail since no resolver is implemented yet
		if validator, exists := registry.GetReferenceValidator("role"); !exists {
			t.Logf("Role reference validator not found - expected for TDD")
			assert.Fail(t, "Role reference resolution not implemented yet")
		} else {
			resolvedID, err := validator.ResolveReference(ctx, "role", "roles/custom.test-custom-role-unique-id")
			if err != nil {
				t.Logf("Reference resolution failed as expected: %v", err)
				assert.Fail(t, "Role reference resolution not implemented yet")
			} else {
				assert.NotEmpty(t, resolvedID, "Resolved ID should not be empty")
			}
		}
	})
}

// TestReferenceSuggestions tests the suggestion system for references
func TestReferenceSuggestions(t *testing.T) {
	ctx := context.Background()
	registry := NewValidatorRegistry()

	t.Run("GroupReferenceSuggestions", func(t *testing.T) {
		// This should fail since no validator is implemented yet
		if validator, exists := registry.GetReferenceValidator("group"); !exists {
			t.Logf("Group reference validator not found - expected for TDD")
			assert.Fail(t, "Group reference suggestions not implemented yet")
		} else {
			suggestions := validator.GetSuggestions(ctx, "group", "test-grup") // typo

			// Should suggest similar group names
			expectedSuggestions := []string{"test-group-unique-id"}
			assert.ElementsMatch(t, expectedSuggestions, suggestions)
		}
	})

	t.Run("RoleReferenceSuggestions", func(t *testing.T) {
		// This should fail since no validator is implemented yet
		if validator, exists := registry.GetReferenceValidator("role"); !exists {
			t.Logf("Role reference validator not found - expected for TDD")
			assert.Fail(t, "Role reference suggestions not implemented yet")
		} else {
			suggestions := validator.GetSuggestions(ctx, "role", "roles/custom.test-rol") // partial

			// Should suggest similar role names
			expectedSuggestions := []string{"roles/custom.test-custom-role-unique-id"}
			assert.ElementsMatch(t, expectedSuggestions, suggestions)
		}
	})
}

// TestCrossResourceReferences tests references between resources as used in simple_test.tf
func TestCrossResourceReferences(t *testing.T) {
	t.Run("RoleBindingReferencesCustomRole", func(t *testing.T) {
		// Tests the pattern: role = "roles/custom.${hiiretail_iam_custom_role.test_custom_role.name}"
		// This should resolve to: "roles/custom.test-custom-role-unique-id"

		ctx := context.Background()
		registry := NewValidatorRegistry()

		expectedRoleReference := "roles/custom.test-custom-role-unique-id"

		if validator, exists := registry.GetReferenceValidator("role"); !exists {
			t.Logf("Role reference validator not found - expected for TDD")
			assert.Fail(t, "Cross-resource role reference validation not implemented yet")
		} else {
			result := validator.ValidateReference(ctx, "role", expectedRoleReference)

			// This should pass if the custom role exists
			assert.True(t, result.Valid, "Role reference should be valid")
		}
	})

	t.Run("RoleBindingReferencesGroup", func(t *testing.T) {
		// Tests the pattern: members = ["group:${hiiretail_iam_group.test_group.name}"]
		// This should resolve to: "group:test-group-unique-id"

		ctx := context.Background()
		registry := NewValidatorRegistry()

		expectedMemberReference := "group:test-group-unique-id"

		if validator, exists := registry.GetReferenceValidator("member"); !exists {
			t.Logf("Member reference validator not found - expected for TDD")
			assert.Fail(t, "Cross-resource member reference validation not implemented yet")
		} else {
			result := validator.ValidateReference(ctx, "member", expectedMemberReference)

			// This should pass if the group exists
			assert.True(t, result.Valid, "Member reference should be valid")
		}
	})
}
