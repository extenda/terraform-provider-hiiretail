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
## Provider Design Principles

### Focus on a Single API or Problem Domain
A Terraform provider MUST manage a single collection of components based on the underlying API or SDK, or a single problem domain. Providers that do not map to a specific API or SDK must be based on a single problem domain or industry standard. This simplifies connectivity, authentication, discovery, and enables maintainers to be experts in a single system.

### Resources Represent a Single API Object
Each Terraform resource MUST be a declarative representation of a single API object, with create, read, delete, and optionally update methods. Abstractions of multiple components or advanced behaviors should be accomplished via Terraform Modules, not provider resources.

### Resource and Attribute Schema Alignment
Resource and attribute schemas MUST closely match the underlying API, unless it degrades user experience or contradicts Terraform expectations. Dates/times should use RFC 3339. Boolean attributes should be oriented so true means to do something. Avoid recursive types. Resources MUST be importable.

### State and Versioning
Providers MUST maintain state continuity and backwards compatibility. Breaking changes require appropriate warnings and deprecation mechanisms. Providers MUST follow Semantic Versioning 2.0.0, with major version for breaking changes, minor for backwards compatible additions, and patch for bug fixes.

### Rationale
These principles maximize predictability, minimize blast radius, simplify maintenance, and enable composition and innovation for operators and maintainers.

## Core Principles

### I. Provider Interface Standards
Every Terraform provider MUST use the official [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) for all development. The provider MUST implement complete CRUD operations for all managed resources, following framework conventions for schema definitions, validation, and state management. Provider configuration must support secure authentication methods and proper error handling with meaningful error messages.

### II. Resource Implementation Completeness
All resources declared in the provider specification MUST have full implementations including Create, Read, Update, Delete operations where applicable. Each resource must handle API errors gracefully, implement proper retry logic, and maintain accurate Terraform state. Partial implementations or schema-only resources are not acceptable for production use.

### III. Authentication & Security (NON-NEGOTIABLE)
OAuth2 client credentials flow MUST be implemented securely with proper token management, refresh handling, and credential validation. Sensitive configuration values must be marked as sensitive in schemas. API credentials must never be logged or exposed in debug output. All HTTP communications must use TLS encryption.

### IV. Testing & Validation
Providers MUST implement a comprehensive testing strategy covering:

- **Unit Tests**: Validate all provider functions, resource CRUD operations, and error handling logic. Unit tests MUST isolate code from external dependencies using mocks or stubs.
- **Integration Tests**: Exercise real API interactions in a test environment, validating authentication, resource lifecycle, and error scenarios. Integration tests MUST use dedicated test accounts and avoid destructive operations on production data.
- **Acceptance Tests**: Providers MUST implement acceptance tests following [Terraform acceptance testing conventions](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests):
		- Acceptance tests MUST cover all resources and data sources, including create, update, delete, import, and error cases.
		- Each test MUST be idempotent, repeatable, and clean up all resources after execution.
		- Tests MUST use the Terraform Plugin Testing framework and be structured for parallel execution where possible.
		- Tests MUST use environment variables for credentials and configuration, never hardcoded secrets.
		- Import tests MUST verify that resources can be imported and state matches the API.
		- Error case tests MUST simulate invalid configurations, API errors, and permission issues.
		- Tests MUST validate resource state after each operation and assert expected errors for negative cases.
		- All acceptance tests MUST report results in CI and block releases on failure.
		- Providers MUST document how to run acceptance tests locally and in CI, including required environment variables and cleanup procedures.
		- Acceptance test cases MUST use the [TestCase](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/testcase) struct to define test steps, configuration, and checks. Each TestCase MUST:
			- Specify preconditions and postconditions for resource state.
			- Use `CheckFunc` to assert resource attributes and error expectations.
			- Include steps for create, update, import, and destroy operations.
			- Simulate error scenarios and validate error messages.
			- Clean up all resources after test execution.
			- Be documented with rationale and expected outcomes.
- **Test Coverage**: All provider code MUST be covered by tests. Critical paths (authentication, resource CRUD, error handling) require 100% coverage. Non-critical code should achieve high coverage and be justified if excluded.
- **Error Simulation**: Tests MUST simulate API errors, network failures, and invalid configurations to verify provider resilience and error reporting.
- **CI Integration**: All tests MUST run in CI pipelines before release. Test failures MUST block releases until resolved.

Mock tests alone are insufficient—real API validation is essential. Providers MUST document their testing strategy and provide instructions for running tests locally and in CI.

### V. Documentation & Examples
Complete provider documentation including resource schemas, configuration examples, and usage patterns must be maintained. Each resource requires working examples showing typical usage scenarios. API compatibility and version requirements must be clearly documented. Migration guides required for breaking changes.

## Technical Standards

All code must follow Go best practices and HashiCorp's provider development guidelines. The provider MUST implement [Terraform Plugin Protocol Version 6](https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol#protocol-version-6) and declare protocol version 6 compatibility in all plugin manifests and binaries. OpenAPI specifications must accurately reflect the target API endpoints and schemas. Generated code must be validated and enhanced with proper business logic. Error handling must be comprehensive with appropriate HTTP status code mapping. Generated code is using [Hashicorp's Terraform OpenAPI Spec Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator) which produces a Provider Code Spec which is used to create the Terraform Provider code. This is documented [here](https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator).

## Quality Assurance

Code reviews must verify compliance with all constitutional principles. Provider releases require successful test suite execution including unit, integration, and acceptance tests. Breaking changes must be documented with migration paths. Performance impact must be evaluated for large-scale resource management scenarios.

## Governance

This constitution supersedes all other development practices for the HiiRetail Terraform Providers project. Amendments require documentation of rationale, approval from project maintainers, and update of dependent templates and documentation. All development decisions must align with these principles - complexity that violates these standards must be justified or refactored.

<<<<<<< HEAD
**Version**: 1.5.0 | **Ratified**: 2025-09-28 | **Last Amended**: 2025-10-11
=======
**Version**: 1.1.0 | **Ratified**: 2025-09-28 | **Last Amended**: 2025-10-11
>>>>>>> docs: amend constitution to v1.1.0 (require Terraform Plugin Framework & Protocol v6)
