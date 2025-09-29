package acceptance_tests

import (
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccProvider_Minimal tests just provider instantiation without terraform operations
func TestAccProvider_Minimal(t *testing.T) {
	// Set TF_ACC for acceptance tests
	t.Setenv("TF_ACC", "1")

	// Set up test environment
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Test just the provider factory creation
	factories := testAccProtoV6ProviderFactoriesWithMockServer()

	// Verify we can create a provider instance
	if len(factories) == 0 {
		t.Fatal("No provider factories created")
	}

	// Try to get a provider instance
	provider, err := factories["hiiretail-iam"]()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}

	t.Log("✅ Provider factory works correctly")
}

// TestAccProvider_ConfigValidation tests just config validation without plan/apply
func TestAccProvider_ConfigValidation(t *testing.T) {
	// Set TF_ACC for acceptance tests
	t.Setenv("TF_ACC", "1")

	// Set up test environment
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfigWithMockServer(env),
				// No PlanOnly - just test provider instantiation
				Check: func(s *terraform.State) error {
					t.Log("✅ Provider configuration validated successfully")
					return nil
				},
			},
		},
	})
}

// TestAccProvider_PlanOnly tests terraform plan without apply
func TestAccProvider_PlanOnly(t *testing.T) {
	// Set TF_ACC for acceptance tests
	t.Setenv("TF_ACC", "1")

	// Set up test environment
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfigWithMockServer(env) + `
resource "hiiretail-iam_custom_role" "test" {
  name         = "test-role"
  display_name = "Test Role"
  description  = "Test role for acceptance testing"
  permissions  = ["iam.roles.list", "iam.roles.get"]
}
`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // We expect a plan for resource creation
				Check: func(s *terraform.State) error {
					t.Log("✅ Plan generated successfully")
					return nil
				},
			},
		},
	})
}
