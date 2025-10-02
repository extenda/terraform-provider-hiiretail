package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/testutils"
)

func TestAccIamCustomRole_planOnly(t *testing.T) {
	// Set TF_ACC for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			// Plan-only testing
			{
				Config:             testAccProviderConfigWithMockServer(env) + testAccIamCustomRoleConfig_basic("test-role-plan"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
