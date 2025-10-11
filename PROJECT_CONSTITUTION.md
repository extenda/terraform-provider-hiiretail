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

**Version**: 1.2.0 | **Ratified**: 2025-09-28 | **Last Amended**: 2025-10-11
