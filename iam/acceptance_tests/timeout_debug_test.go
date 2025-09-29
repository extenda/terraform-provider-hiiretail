package acceptance_tests

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
)

// TestAccProvider_PlanTimeoutDebugging tests provider plan operation with detailed timing
func TestAccProvider_PlanTimeoutDebugging(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready by testing OAuth2 endpoint
	env.ValidateMockServerReady(t)

	t.Logf("Mock server running at: %s", env.BaseURL)

	// Test with a shorter timeout to identify where the hang occurs
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			// Plan-only testing with debug
			{
				Config: testAccProviderConfigWithMockServer(env) + `
resource "hiiretail-iam_custom_role" "debug" {
  id = "debug-role-timeout"
  name = "Debug Custom Role"
  permissions = [
    {
      id = "pos.payment.create"
    }
  ]
}`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("Plan completed successfully!")
						return nil
					},
				),
			},
		},
	})
}

// TestAccProvider_ConfigurationOnly tests just provider configuration without resources
func TestAccProvider_ConfigurationOnly(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	t.Logf("Testing provider configuration only at: %s", env.BaseURL)

	// Test that provider can be configured without any resources
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfigWithMockServer(env),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// TestProviderFactory_Instantiation tests that the provider factory creates a working provider
func TestProviderFactory_Instantiation(t *testing.T) {
	factory := testAccProviderFactoryWithMockEnvironment()

	server, err := factory()
	if err != nil {
		t.Fatalf("Provider factory failed: %v", err)
	}

	if server == nil {
		t.Fatal("Provider factory returned nil server")
	}

	t.Log("Provider factory creates server successfully")
}

// TestProviderCommunication_WithTimeout tests provider communication with controlled timeout
func TestProviderCommunication_WithTimeout(t *testing.T) {
	// Set TF_ACC environment variable for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)
	env.ValidateMockServerReady(t)

	// Create context with shorter timeout for debugging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use context in test setup (this would need to be supported by the testing framework)
	t.Logf("Testing provider communication with 30s timeout at: %s", env.BaseURL)

	// Simple test to validate provider can start and respond
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckWithMock(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithMockServer(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfigWithMockServer(env),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
							t.Log("Provider communication test passed")
							return nil
						}
					},
				),
			},
		},
	})
}
