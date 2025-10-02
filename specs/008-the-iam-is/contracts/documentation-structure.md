# Documentation Structure Contract

## Documentation Organization Contract

**Root Structure**:
```
docs/
├── index.md                     # Provider overview and quick start
├── guides/                      # User guides and tutorials
│   ├── getting-started.md       # First-time setup guide
│   ├── authentication.md        # OAuth2 configuration guide
│   ├── provider-comparison.md   # vs AWS/GCP/Azure patterns
│   └── migration-guides/
│       └── from-hiiretail-iam.md # Migration from old provider
├── resources/                   # Resource documentation
│   ├── iam/
│   │   ├── overview.md          # IAM service introduction
│   │   ├── group.md             # hiiretail_iam_group
│   │   ├── custom_role.md       # hiiretail_iam_custom_role
│   │   └── role_binding.md      # hiiretail_iam_role_binding
│   └── ccc/ (future)
│       └── overview.md
├── data-sources/                # Data source documentation
│   └── iam/
│       └── groups.md            # hiiretail_iam_groups
└── examples/                    # Complete working examples
    ├── basic-iam-setup/
    ├── multi-service-deployment/
    └── enterprise-patterns/
```

## Documentation Content Contract

**Provider Overview** (`docs/index.md`):
- Clear explanation of multi-API capability
- Service overview with links to detailed docs
- Quick configuration example
- Link to getting started guide

**Service Overview Pages** (`docs/resources/{service}/overview.md`):
- Service purpose and key concepts
- Available resources and data sources
- Service-specific configuration options
- Common usage patterns and examples

**Resource Documentation** (`docs/resources/{service}/{resource}.md`):
- Resource purpose and functionality
- Complete schema documentation (required/optional/computed attributes)
- At least 3 working examples (basic, intermediate, advanced)
- Related resources and common patterns
- Troubleshooting section for common issues

**Getting Started Guide** (`docs/guides/getting-started.md`):
- Prerequisites and setup requirements
- Step-by-step first resource creation
- Validation and troubleshooting steps
- Next steps and additional resources

## Example Quality Contract

**All Examples Must**:
- Be complete, runnable Terraform configurations
- Include necessary provider configuration
- Show realistic use cases, not toy examples
- Include comments explaining key concepts
- Be tested and validated before publication

**Example Categories**:
- **Basic**: Single resource with minimal configuration
- **Intermediate**: Multiple resources with relationships
- **Advanced**: Complex scenarios with dependencies and data sources
- **Migration**: Before/after showing old vs new patterns

## Auto-Generation Contract

**Generated Content**:
- Resource schema documentation from Go code
- Example validation from actual test files
- Version compatibility matrices from CI/CD
- Performance benchmarks from automated tests

**Manual Content**:
- Service overviews and conceptual explanations
- Getting started guides and tutorials
- Migration guides and best practices
- Troubleshooting and FAQ sections

## User Experience Contract

**Navigation**:
- Service-based organization for easy discovery
- Cross-references between related resources
- Search functionality for large documentation sets
- Mobile-responsive design for various devices

**Content Quality**:
- Clear, concise writing appropriate for technical audience
- Consistent terminology and formatting
- Regular updates aligned with provider releases  
- User feedback integration and continuous improvement