output "basic_store" {
  description = "Details of the basic store resource"
  value = {
    id        = hiiretail_iam_resource.basic_store.id
    name      = hiiretail_iam_resource.basic_store.name
    tenant_id = hiiretail_iam_resource.basic_store.tenant_id
  }
}

output "basic_department" {
  description = "Details of the basic department resource"
  value = {
    id        = hiiretail_iam_resource.basic_department.id
    name      = hiiretail_iam_resource.basic_department.name
    tenant_id = hiiretail_iam_resource.basic_department.tenant_id
  }
}

output "resource_ids" {
  description = "List of all resource IDs created"
  value = [
    hiiretail_iam_resource.basic_store.id,
    hiiretail_iam_resource.basic_department.id
  ]
}