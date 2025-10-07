# Changelog

All notable changes to the HiiRetail Terraform Provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release preparation and Terraform Registry publishing setup

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [1.0.0] - TBD

### Added
- **Provider Distribution**: Automated publishing to Terraform Registry via GitHub Actions
- **Installation Methods**: Support for installation from Terraform Registry, manual installation, and local development
- **Documentation**: Comprehensive documentation for registry installation and usage
- **Examples**: Complete examples demonstrating provider installation from Terraform Registry
- **Release Automation**: GoReleaser configuration for multi-platform releases with GPG signing
- **Version Constraints**: Support for semantic versioning and version constraints

### Changed
- **Provider Source**: Updated from development source to official registry source `extenda/hiiretail`
- **Installation**: Simplified installation process via Terraform Registry
- **Documentation**: Enhanced provider documentation with registry-specific instructions

### Security
- **GPG Signing**: All release artifacts are GPG signed for integrity verification
- **Automated Releases**: Secure release pipeline with GitHub Actions and repository secrets

## [0.1.0] - Initial Development Release

### Added
- **IAM Group Management**: Create, read, update, and delete IAM groups
- **Custom Role Management**: Full CRUD operations for custom roles with permissions
- **Resource Management**: Manage IAM resources with proper typing and validation
- **Role Binding Management**: Complete role binding functionality linking groups, roles, and resources
- **OAuth2 Authentication**: Secure authentication using client credentials flow
- **Multi-tenant Support**: Tenant-scoped operations with flexible tenant configuration
- **Environment Variables**: Support for multiple environment variable formats (HIIRETAIL_*, TF_VAR_*)
- **Error Handling**: Comprehensive error handling with descriptive messages
- **Validation**: Input validation for all resource attributes
- **Documentation**: Auto-generated documentation using terraform-plugin-docs
- **Examples**: Working examples for all resources and configurations
- **Testing**: Unit tests and acceptance tests for all resources

### Technical Details
- Built with Terraform Plugin Framework v1.4.2
- Go 1.21+ compatibility
- OAuth2 integration with automatic token refresh
- Comprehensive logging and debugging support
- Cross-platform compatibility (Linux, macOS, Windows)

---

## Release Notes Format

Each release follows this structure:

### Added
- New features and capabilities

### Changed  
- Changes to existing functionality

### Deprecated
- Features marked for removal in future versions

### Removed
- Features removed in this version

### Fixed
- Bug fixes and corrections

### Security
- Security-related changes and improvements

## Version Numbering

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions  
- **PATCH** version for backwards-compatible bug fixes

## Release Process

1. Update CHANGELOG.md with new version and changes
2. Update version constraints in examples and documentation
3. Create and push version tag: `git tag v1.0.0 && git push origin v1.0.0`
4. GitHub Actions automatically creates release with artifacts
5. Terraform Registry automatically detects and publishes new version

## Links

- [GitHub Releases](https://github.com/extenda/terraform-provider-hiiretail/releases)
- [Terraform Registry](https://registry.terraform.io/providers/extenda/hiiretail)
- [Documentation](https://registry.terraform.io/providers/extenda/hiiretail/latest/docs)