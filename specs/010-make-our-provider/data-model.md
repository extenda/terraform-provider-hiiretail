# Data Model: Provider Distribution and Availability

## Overview
Data models and entities involved in the Terraform provider distribution system.

## Core Entities

### Release Artifact
**Description**: Binary distribution package for a specific platform and version
**Attributes**:
- `filename`: Binary filename (e.g., `terraform-provider-hiiretail_v1.0.0_linux_amd64.zip`)
- `platform`: Target platform (linux/amd64, darwin/arm64, windows/amd64, etc.)
- `version`: Semantic version string (v1.0.0)
- `checksum`: SHA256 hash of the binary
- `signature`: GPG signature for verification
- `size`: File size in bytes

**Relationships**: 
- Belongs to one Release
- Has one Checksum entry
- Has one GPG Signature

### Release
**Description**: Complete provider release containing multiple platform artifacts
**Attributes**:
- `tag`: Git tag (v1.0.0)
- `version`: Semantic version without 'v' prefix (1.0.0)
- `name`: Release name/title
- `body`: Release notes/changelog
- `draft`: Boolean release status
- `prerelease`: Boolean pre-release indicator
- `created_at`: Release creation timestamp
- `published_at`: Release publication timestamp

**Relationships**:
- Has many Release Artifacts
- Has one SHASUMS file
- Has one SHASUMS signature

### Build Configuration
**Description**: GoReleaser build configuration for cross-compilation
**Attributes**:
- `project_name`: terraform-provider-hiiretail
- `binary_template`: terraform-provider-{{ .ProjectName }}_v{{ .Version }}
- `platforms`: Array of target OS/arch combinations
- `archive_format`: Archive type per platform (zip/tar.gz)
- `checksum_template`: SHASUMS filename pattern

**Relationships**:
- Generates multiple Release Artifacts
- Defines build matrix for CI/CD

### Documentation Bundle
**Description**: Provider documentation and examples for distribution
**Attributes**:
- `readme`: Installation and usage guide
- `provider_docs`: Generated provider schema documentation
- `resource_docs`: Individual resource documentation files
- `examples`: Working Terraform configuration examples
- `changelog`: Version history and breaking changes

**Relationships**:
- Packaged with each Release
- Synchronized with Terraform Registry

## State Transitions

### Release Lifecycle
1. **Draft**: Release created but not published
2. **Published**: Release available for download
3. **Deprecated**: Older release, users encouraged to upgrade
4. **Yanked**: Release removed due to critical issues

### Build Process Flow
1. **Triggered**: Tag push triggers GitHub Actions
2. **Building**: Cross-compilation in progress
3. **Signing**: GPG signature generation
4. **Packaging**: Archive creation and checksum generation
5. **Publishing**: GitHub release creation
6. **Registry Sync**: Terraform Registry detects and indexes release

## Validation Rules

### Version Constraints
- Must follow semantic versioning (MAJOR.MINOR.PATCH)
- Git tags must include 'v' prefix (v1.0.0)
- No duplicate versions allowed
- Pre-release versions use suffix (-alpha, -beta, -rc)

### Artifact Requirements
- All supported platforms must have artifacts
- Checksums required for all binaries
- GPG signatures mandatory for security
- Archive format consistent per platform

### Documentation Standards
- README must include installation instructions
- Examples must be valid Terraform configurations
- Resource documentation auto-generated from schemas
- Changelog follows conventional format

## Data Flow

```
Git Tag (v1.0.0) 
  → GitHub Actions Trigger
  → GoReleaser Build Matrix
  → Platform Binaries Generation
  → Checksum Calculation
  → GPG Signing
  → Archive Creation
  → GitHub Release Publication
  → Terraform Registry Sync
```

## Storage and Persistence

**GitHub Releases**: Primary storage for release artifacts and metadata
**Terraform Registry**: Indexed metadata and documentation mirror
**Build Artifacts**: Temporary storage during CI/CD process (auto-cleanup)
**Documentation**: Version-controlled in repository, deployed to registry