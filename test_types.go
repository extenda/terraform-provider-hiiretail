package main

import "github.com/hashicorp/terraform-plugin-framework/resource"

// Keep references to framework types to ensure the module dependency is present.
var _ resource.MetadataRequest
var _ resource.MetadataResponse
