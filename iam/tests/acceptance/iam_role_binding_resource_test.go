package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIAMRoleBinding tests IAM Role Binding resource CRUD operations
func TestAccIAMRoleBinding(t *testing.T) {
	// This test should FAIL until IAM Role Binding resource is implemented
	bindingName := "test-binding-" + randomString(8)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIAMRoleBindingConfig(bindingName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "name", bindingName),
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "role", "roles/iam.viewer"),
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "members.#", "2"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_role_binding.test", "id"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_role_binding.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hiiretail_iam_role_binding.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIAMRoleBindingConfigUpdated(bindingName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "name", bindingName),
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "role", "roles/iam.editor"),
					resource.TestCheckResourceAttr("hiiretail_iam_role_binding.test", "members.#", "3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIAMRoleBindingConfig(name string) string {
	return testAccProviderConfig + fmt.Sprintf(`
resource "hiiretail_iam_role_binding" "test" {
  name = "%s"
  role = "roles/iam.viewer"
  members = [
    "user:test-user1@example.com",
    "group:test-group1"
  ]
}
`, name)
}

func testAccIAMRoleBindingConfigUpdated(name string) string {
	return testAccProviderConfig + fmt.Sprintf(`
resource "hiiretail_iam_role_binding" "test" {
  name = "%s"
  role = "roles/iam.editor"
  members = [
    "user:test-user1@example.com",
    "user:test-user2@example.com", 
    "group:test-group1"
  ]
}
`, name)
}
