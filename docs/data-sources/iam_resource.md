---
page_title: "hiiretail_iam_resource Data Source - terraform-provider-hiiretail"
subcategory: "IAM"
description: |-
  Fetches an IAM resource by ID from HiiRetail.
---

# hiiretail_iam_resource (Data Source)

Fetches an IAM resource by ID from HiiRetail. This is useful for referencing resources that are managed outside of Terraform.

## Example Usage

```terraform
# Fetch a store resource
data "hiiretail_iam_resource" "my_store" {
  id = "store:001"
}

# Use the resource in a role binding
resource "hiiretail_iam_role_binding" "store_access" {
  group_id = "sales-team"
  roles = [
    {
      role_id = "roles/store.viewer"
      bindings = [
        {
          type       = "resource"
          subject_id = data.hiiretail_iam_resource.my_store.id
        }
      ]
    }
  ]
}

# Output resource details
output "store_name" {
  value = data.hiiretail_iam_resource.my_store.name
}

output "store_properties" {
  value = data.hiiretail_iam_resource.my_store.properties
}
```

## Schema

### Required

- `id` (String) The unique identifier of the resource in the format `type:name` (e.g., `store:001`, `product:abc123`).

### Read-Only

- `name` (String) The name of the resource.
- `properties` (String) JSON string containing additional properties of the resource.

## API Documentation

This data source uses the [GetResource API endpoint](https://developer.hiiretail.com/api/iam-api#tag/Resources/operation/getResource).
