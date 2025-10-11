<!--
Sync Impact Report:
Version change: Initial → 1.0.0
Added principles:
- I. Provider Interface Standards
- II. Resource Implementation Completeness  
- III. Authentication & Security (NON-NEGOTIABLE)
- IV. Testing & Validation
- V. Documentation & Examples
Added sections:
- Technical Standards
- Quality Assurance
Templates requiring updates: ✅ All templates aligned with Terraform provider development standards
Follow-up TODOs: None - all placeholders filled
-->

# HiiRetail Terraform Providers Constitution

## Core Principles

### I. Provider Interface Standards
Every Terraform provider MUST use the official [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) for all development. The provider MUST implement complete CRUD operations for all managed resources, following framework conventions for schema definitions, validation, and state management. Provider configuration must support secure authentication methods and proper error handling with meaningful error messages.

### II. Resource Implementation Completeness
All resources declared in the provider specification MUST have full implementations including Create, Read, Update, Delete operations where applicable. Each resource must handle API errors gracefully, implement proper retry logic, and maintain accurate Terraform state. Partial implementations or schema-only resources are not acceptable for production use.

### III. Authentication & Security (NON-NEGOTIABLE)
OAuth2 client credentials flow MUST be implemented securely with proper token management, refresh handling, and credential validation. Sensitive configuration values must be marked as sensitive in schemas. API credentials must never be logged or exposed in debug output. All HTTP communications must use TLS encryption.

### IV. Testing & Validation
Unit tests are mandatory for all provider functions and resource operations. Integration tests must validate actual API interactions with test environments. Acceptance tests following Terraform testing conventions are required before any release. Mock tests alone are insufficient - real API validation is essential.

### V. Documentation & Examples
Complete provider documentation including resource schemas, configuration examples, and usage patterns must be maintained. Each resource requires working examples showing typical usage scenarios. API compatibility and version requirements must be clearly documented. Migration guides required for breaking changes.

## Technical Standards

All code must follow Go best practices and HashiCorp's provider development guidelines. The provider MUST implement [Terraform Plugin Protocol Version 6](https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol#protocol-version-6) and declare protocol version 6 compatibility in all plugin manifests and binaries. OpenAPI specifications must accurately reflect the target API endpoints and schemas. Generated code must be validated and enhanced with proper business logic. Error handling must be comprehensive with appropriate HTTP status code mapping. Generated code is using [Hashicorp's Terraform OpenAPI Spec Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator) which produces a Provider Code Spec which is used to create the Terraform Provider code. This is documented [here](https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator).

## Quality Assurance

Code reviews must verify compliance with all constitutional principles. Provider releases require successful test suite execution including unit, integration, and acceptance tests. Breaking changes must be documented with migration paths. Performance impact must be evaluated for large-scale resource management scenarios.

## Governance

This constitution supersedes all other development practices for the HiiRetail Terraform Providers project. Amendments require documentation of rationale, approval from project maintainers, and update of dependent templates and documentation. All development decisions must align with these principles - complexity that violates these standards must be justified or refactored.

**Version**: 1.1.0 | **Ratified**: 2025-09-28 | **Last Amended**: 2025-10-11