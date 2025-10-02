package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIAMGroup tests IAM Group resource CRUD operations
func TestAccIAMGroup(t *testing.T) {
	// This test should FAIL until IAM Group resource is implemented
	groupName := "test-group-" + randomString(8)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIAMGroupConfig(groupName, "Test group description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_group.test", "name", groupName),
					resource.TestCheckResourceAttr("hiiretail_iam_group.test", "description", "Test group description"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_group.test", "id"),
					resource.TestCheckResourceAttrSet("hiiretail_iam_group.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hiiretail_iam_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIAMGroupConfig(groupName, "Updated test group description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hiiretail_iam_group.test", "name", groupName),
					resource.TestCheckResourceAttr("hiiretail_iam_group.test", "description", "Updated test group description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIAMGroupConfig(name, description string) string {
	return testAccProviderConfig + fmt.Sprintf(`
resource "hiiretail_iam_group" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}
