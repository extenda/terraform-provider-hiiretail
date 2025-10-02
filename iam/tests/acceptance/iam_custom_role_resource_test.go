package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIAMCustomRole tests IAM Custom Role resource CRUD operations
func TestAccIAMCustomRole(t *testing.T) {
	// This test should FAIL until IAM Custom Role resource is implemented
	roleName := "test-role-" + randomString(8)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIAMCustomRoleConfig(roleName, "Test Custom Role", "Test role description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "name", roleName),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "title", "Test Custom Role"),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "description", "Test role description"),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "permissions.#", "2"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_custom_role.test", "id"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_custom_role.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hiiretail_iam_custom_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIAMCustomRoleConfigUpdated(roleName, "Updated Test Custom Role", "Updated test role description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "name", roleName),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "title", "Updated Test Custom Role"),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "description", "Updated test role description"),
					resource.TestCheckResourceAttr("hiiretail_iam_custom_role.test", "permissions.#", "3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIAMCustomRoleConfig(name, title, description string) string {
	return testAccProviderConfig + fmt.Sprintf(`
resource "hiiretail_iam_custom_role" "test" {
  name        = "%s"
  title       = "%s"
  description = "%s"
  permissions = [
    "iam.groups.list",
    "iam.groups.get"
  ]
  stage = "ALPHA"
}
`, name, title, description)
}

func testAccIAMCustomRoleConfigUpdated(name, title, description string) string {
	return testAccProviderConfig + fmt.Sprintf(`
resource "hiiretail_iam_custom_role" "test" {
  name        = "%s"
  title       = "%s"
  description = "%s"
  permissions = [
    "iam.groups.list",
    "iam.groups.get",
    "iam.groups.create"
  ]
  stage = "BETA"
}
`, name, title, description)
}
