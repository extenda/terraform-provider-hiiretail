# Basic HiiRetail IAM Resource Example

This example demonstrates the simplest possible use of the `hiiretail_iam_resource` resource.

## Usage

1. Configure your provider credentials (see variables.tf for options)
2. Update terraform.tfvars with your values
3. Run: `terraform init && terraform plan && terraform apply`

## What This Creates

- A single IAM resource with ID `store:001`
- No additional properties (props field empty)
- Demonstrates basic resource lifecycle (create, read, update, delete)

## Files

- `main.tf` - Main configuration
- `variables.tf` - Variable definitions  
- `outputs.tf` - Output values
- `terraform.tfvars.example` - Example variable values
- `README.md` - This file