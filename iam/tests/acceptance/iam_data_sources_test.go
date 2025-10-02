package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIAMGroupsDataSource tests IAM Groups data source
func TestAccIAMGroupsDataSource(t *testing.T) {
	// This test should FAIL until IAM Groups data source is implemented
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMGroupsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hiiretail_iam_groups.test", "groups.#"),
					resource.TestCheckResourceAttrSet("data.hiiretail_iam_groups.test", "id"),
				),
			},
		},
	})
}

func testAccIAMGroupsDataSourceConfig() string {
	return testAccProviderConfig + `
data "hiiretail_iam_groups" "test" {
  filter = "name:test-*"
}
`
}

// TestAccIAMRolesDataSource tests IAM Roles data source
func TestAccIAMRolesDataSource(t *testing.T) {
	// This test should FAIL until IAM Roles data source is implemented
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMRolesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hiiretail_iam_roles.test", "roles.#"),
					resource.TestCheckResourceAttrSet("data.hiiretail_iam_roles.test", "id"),
				),
			},
		},
	})
}

func testAccIAMRolesDataSourceConfig() string {
	return testAccProviderConfig + `
data "hiiretail_iam_roles" "test" {
  filter = "type:custom"
}
`
}
