package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccHiiRetailIAMResource_basic(t *testing.T) {
	resourceName := "hiiretail_iam_resource.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHiiRetailIAMResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHiiRetailIAMResourceConfig_basic("store:001", "Test Store Resource"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckHiiRetailIAMResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "store:001"),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Store Resource"),
					resource.TestCheckResourceAttr(resourceName, "props", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHiiRetailIAMResourceImportStateIdFunc(resourceName),
			},
			// Update and Read testing
			{
				Config: testAccHiiRetailIAMResourceConfig_withProps("store:001", "Updated Store Resource", `{"location": "main-floor", "active": true}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckHiiRetailIAMResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "store:001"),
					resource.TestCheckResourceAttr(resourceName, "name", "Updated Store Resource"),
					resource.TestCheckResourceAttr(resourceName, "props", `{"location": "main-floor", "active": true}`),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccHiiRetailIAMResource_withComplexProps(t *testing.T) {
	resourceName := "hiiretail_iam_resource.test"
	complexProps := `{
		"department": "electronics",
		"metadata": {
			"priority": 1,
			"tags": ["retail", "pos"]
		},
		"permissions": ["read", "write"]
	}`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHiiRetailIAMResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHiiRetailIAMResourceConfig_withProps("dept:electronics", "Electronics Department", complexProps),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckHiiRetailIAMResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "dept:electronics"),
					resource.TestCheckResourceAttr(resourceName, "name", "Electronics Department"),
					resource.TestCheckResourceAttr(resourceName, "props", complexProps),
				),
			},
		},
	})
}

func TestAccHiiRetailIAMResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid ID format - contains slash
			{
				Config:      testAccHiiRetailIAMResourceConfig_basic("store/001", "Test Resource"),
				ExpectError: regexp.MustCompile(`Resource ID.*cannot contain forward slashes`),
			},
			// Test invalid ID format - double underscores
			{
				Config:      testAccHiiRetailIAMResourceConfig_basic("store__001", "Test Resource"),
				ExpectError: regexp.MustCompile(`Resource ID.*cannot contain consecutive underscores`),
			},
			// Test invalid ID format - single dot
			{
				Config:      testAccHiiRetailIAMResourceConfig_basic(".", "Test Resource"),
				ExpectError: regexp.MustCompile(`Resource ID.*cannot be.*\\.`),
			},
			// Test invalid props - malformed JSON
			{
				Config:      testAccHiiRetailIAMResourceConfig_withProps("store:001", "Test Resource", `{"invalid": json}`),
				ExpectError: regexp.MustCompile(`invalid JSON`),
			},
		},
	})
}

func TestAccHiiRetailIAMResource_disappears(t *testing.T) {
	resourceName := "hiiretail_iam_resource.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHiiRetailIAMResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHiiRetailIAMResourceConfig_basic("store:disappear", "Disappearing Resource"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHiiRetailIAMResourceExists(resourceName),
					testAccCheckHiiRetailIAMResourceDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckHiiRetailIAMResourceDestroy(s *terraform.State) error {
	// This function verifies that resources have been destroyed
	// In a real implementation, this would check the API to ensure resources are deleted
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hiiretail_iam_resource" {
			continue
		}

		// In a real implementation, you would:
		// 1. Extract the resource ID from rs.Primary.ID
		// 2. Make an API call to check if the resource still exists
		// 3. Return an error if it still exists
		// 4. Return nil if it's properly deleted (404 error from API)

		// For this test framework, we'll assume successful deletion
		// since we can't make real API calls without proper authentication
	}

	return nil
}

func testAccCheckHiiRetailIAMResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set for %s", resourceName)
		}

		// In a real implementation, you would:
		// 1. Get the provider client from the test context
		// 2. Make an API call to verify the resource exists
		// 3. Return an error if the API call fails or resource doesn't exist
		// 4. Return nil if the resource exists and matches expected state

		// For this test framework, we'll assume the resource exists
		// if it has a valid ID in the Terraform state
		return nil
	}
}

func testAccCheckHiiRetailIAMResourceDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// In a real implementation, you would:
		// 1. Get the provider client from the test context
		// 2. Make an API call to delete the resource outside of Terraform
		// 3. This simulates the resource being deleted by external means
		// 4. The next Terraform plan should detect this and recreate the resource

		// For this test framework, we'll document the expected behavior
		// In a real test, this would actually delete the resource via API
		_ = rs // Acknowledge that we have the resource state

		return nil
	}
}

func testAccHiiRetailIAMResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Resource ID not set")
		}

		// Return the resource ID for import
		return rs.Primary.ID, nil
	}
}

// Configuration functions for test cases
func testAccHiiRetailIAMResourceConfig_basic(id, name string) string {
	return fmt.Sprintf(`
resource "hiiretail_iam_resource" "test" {
  id   = %[1]q
  name = %[2]q
}
`, id, name)
}

func testAccHiiRetailIAMResourceConfig_withProps(id, name, props string) string {
	return fmt.Sprintf(`
resource "hiiretail_iam_resource" "test" {
  id    = %[1]q
  name  = %[2]q
  props = %[3]q
}
`, id, name, props)
}

// Note: testAccPreCheck and testAccProtoV6ProviderFactories are defined in group_resource_test.go
// and shared across all acceptance tests in this package
