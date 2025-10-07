package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	resource_iam_resource "github.com/extenda/hiiretail-terraform-providers/internal/provider/resource_iam_resource"
)

// NewResourceResource creates a new iam_resource resource
func NewResourceResource() resource.Resource {
	return &resource_iam_resource.IAMResourceResource{}
}
