package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/testutils"
)

// TestAccIamRoleBindingResource_basic tests the basic CRUD lifecycle
func TestAccIamRoleBindingResource_basic(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	resourceName := "hiiretail_iam_role_binding.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		CheckDestroy:             testAccCheckIamRoleBindingDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIamRoleBindingConfig_basic("test-role-123", "user:user-456"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIamRoleBindingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_id", "test-role-123"),
					resource.TestCheckResourceAttr(resourceName, "is_custom", "true"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0", "user:user-456"),
				),
			},
			// Update testing - change binding
			{
				Config: testAccIamRoleBindingConfig_basic("test-role-123", "group:test-group"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIamRoleBindingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_id", "test-role-123"),
					resource.TestCheckResourceAttr(resourceName, "is_custom", "true"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0", "group:test-group"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIamRoleBindingResource_maxBindings tests the maximum bindings limit
func TestAccIamRoleBindingResource_maxBindings(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Create 10 bindings (maximum allowed)
	bindings := make([]string, 10)
	for i := 0; i < 10; i++ {
		bindings[i] = fmt.Sprintf("user:user-%d", i)
	}

	resourceName := "hiiretail_iam_role_binding.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		CheckDestroy:             testAccCheckIamRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIamRoleBindingConfig_maxBindings("test-role-max", bindings),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIamRoleBindingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_id", "test-role-max"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "10"),
				),
			},
		},
	})
}

// TestAccIamRoleBindingResource_exceedsMaxBindings tests error handling for too many bindings
func TestAccIamRoleBindingResource_exceedsMaxBindings(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Create 11 bindings (exceeds maximum allowed)
	bindings := make([]string, 11)
	for i := 0; i < 11; i++ {
		bindings[i] = fmt.Sprintf("user:user-%d", i)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config:      testAccIamRoleBindingConfig_maxBindings("test-role-exceed", bindings),
				ExpectError: regexp.MustCompile("exceeds maximum allowed bindings \\(10\\)"),
			},
		},
	})
}

// TestAccIamRoleBindingResource_invalidBinding tests error handling for invalid binding formats
func TestAccIamRoleBindingResource_invalidBinding(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config:      testAccIamRoleBindingConfig_basic("test-role-invalid", "invalid-binding-format"),
				ExpectError: regexp.MustCompile("invalid binding format"),
			},
		},
	})
}

// TestAccIamRoleBindingResource_multipleBindings tests multiple bindings with different formats
func TestAccIamRoleBindingResource_multipleBindings(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	resourceName := "hiiretail_iam_role_binding.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		CheckDestroy:             testAccCheckIamRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIamRoleBindingConfig_multipleBindings(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIamRoleBindingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_id", "test-role-multi"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0", "user:user-1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.1", "group:group-1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.2", "serviceAccount:service-1"),
				),
			},
		},
	})
}

// Helper functions for tests

func testAccCheckIamRoleBindingDestroy(s *terraform.State) error {
	// Note: This will be implemented in Phase 3.3 - Core Implementation
	// For now, this is a placeholder that ensures the test structure is correct
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hiiretail_iam_role_binding" {
			continue
		}

		// TODO: Check that the role binding has been destroyed
		// This will involve calling the actual API to verify deletion
	}
	return nil
}

func testAccCheckIamRoleBindingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Note: This will be implemented in Phase 3.3 - Core Implementation
		// For now, this is a placeholder that ensures the test structure is correct
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		// TODO: Check that the role binding exists in the API
		// This will involve calling the actual API to verify existence
		return nil
	}
}

// Terraform configuration templates

func testAccIamRoleBindingConfig_basic(roleId, binding string) string {
	return fmt.Sprintf(`
resource "hiiretail_iam_role_binding" "test" {
  role_id    = "%s"
  is_custom  = true
  bindings   = ["%s"]
}
`, roleId, binding)
}

func testAccIamRoleBindingConfig_maxBindings(roleId string, bindings []string) string {
	bindingsStr := ""
	for i, binding := range bindings {
		if i > 0 {
			bindingsStr += ", "
		}
		bindingsStr += fmt.Sprintf(`"%s"`, binding)
	}

	return fmt.Sprintf(`
resource "hiiretail_iam_role_binding" "test" {
  role_id    = "%s"
  is_custom  = true
  bindings   = [%s]
}
`, roleId, bindingsStr)
}

func testAccIamRoleBindingConfig_multipleBindings() string {
	return `
resource "hiiretail_iam_role_binding" "test" {
  role_id    = "test-role-multi"
  is_custom  = true
  bindings   = [
    "user:user-1",
    "group:group-1", 
    "serviceAccount:service-1"
  ]
}
`
}
