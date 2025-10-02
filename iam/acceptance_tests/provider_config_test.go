package acceptance_tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/testutils"
)

// testAccPreCheckWithMock validates environment for mock-based acceptance tests
func testAccPreCheckWithMock(t *testing.T) {
	// TF_ACC must be set to run acceptance tests
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("TF_ACC not set, skipping acceptance test")
	}
}

// testAccProviderConfigWithMockServer returns provider configuration for mock server tests
func testAccProviderConfigWithMockServer(env *testutils.TestEnvironment) string {
	return `
provider "hiiretail-iam" {
  base_url      = "` + env.BaseURL + `"
  tenant_id     = "test-tenant-123"
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
}
`
}

// TestAccProvider verifies the provider can be instantiated
func TestAccProvider(t *testing.T) {
	// This test just verifies the provider factory works
	p := provider.New("test")()
	if p == nil {
		t.Fatal("Failed to create provider")
	}
}

// TestAccProviderConfigure tests provider configuration with mock environment
func TestAccProviderConfigure(t *testing.T) {
	// Set TF_ACC for this test
	t.Setenv("TF_ACC", "1")

	// Set up test environment
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	// Test that provider can be configured and used in Terraform context
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

// testAccProviderFactoryWithMockEnvironment creates a provider factory that uses mock server
func testAccProviderFactoryWithMockEnvironment() func() (tfprotov6.ProviderServer, error) {
	return providerserver.NewProtocol6WithError(provider.New("test")())
}

// testAccProtoV6ProviderFactoriesWithMockServer returns provider factories configured for mock server
func testAccProtoV6ProviderFactoriesWithMockServer() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hiiretail-iam": testAccProviderFactoryWithMockEnvironment(),
	}
}
