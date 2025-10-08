# HiiRetail IAM Resources Example

This example demonstrates the use of Hii Retail IAM Resources to manage Groups, Roles, Resources and Role Bindings

## Usage

1. Configure your provider credentials
2. Set HIIRETAIL environment variables or update terraform.tfvars with your values
3. Run: `terraform init && terraform plan && terraform apply`

## What This Creates

- A single IAM Group with Name `Store-tf01-Financial-Managers`
- A single Custom Role called `ReconcilliationApprover`
- A single IAM Resource representating a Store Business Unit `bu:tf01`
- A Role Binding on Group `Store-tf01-Financial-Managers` with Custom Role `ReconcilliationApprover` for Resource `bu:tf01`
- A Role Bindig on Group `Store-tf01-Financial-Managers` with Built-in Role `rec.manager` for Resource `bu:tf01`

## Files

- `main.tf` - Main configuration
- `variables.tf` - Variable definitions  
- `terraform.tfvars.example` - Example variable values
- `README.md` - This file