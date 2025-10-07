# Maintainer Release Process Documentation

This document provides comprehensive guidance for maintainers on managing releases of the HiiRetail Terraform provider.

## 🚀 Release Workflow Overview

The provider uses an automated release process through GitHub Actions, triggered by semantic version tags.

### Release Types

1. **Major Release** (e.g., v1.0.0, v2.0.0)
   - Breaking changes to provider API
   - Requires migration guide
   - Full regression testing

2. **Minor Release** (e.g., v1.1.0, v1.2.0)
   - New features and resources
   - Backward compatible
   - Extended testing

3. **Patch Release** (e.g., v1.0.1, v1.0.2)
   - Bug fixes only
   - Backward compatible
   - Standard testing

4. **Pre-release** (e.g., v1.0.0-beta, v1.0.0-rc1)
   - Testing and validation
   - Not for production use
   - Community feedback

## 📋 Pre-Release Checklist

### Code Preparation
- [ ] All intended changes merged to main branch
- [ ] Code review completed and approved
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] Documentation updated
- [ ] Examples verified

### Version Planning
- [ ] Determine semantic version number
- [ ] Check for breaking changes
- [ ] Review dependency updates
- [ ] Plan communication strategy

### Environment Verification
- [ ] GitHub repository secrets configured
  - [ ] GPG_PRIVATE_KEY
  - [ ] PASSPHRASE
- [ ] GoReleaser configuration validated
- [ ] GitHub Actions workflow tested

## 🔧 Release Process Steps

### 1. Prepare Release Branch
```bash
# Start from main branch
git checkout main
git pull origin main

# Create release branch (for major/minor releases)
git checkout -b release/v1.1.0

# Or work directly on main for patches
```

### 2. Update Version References
```bash
# Update version in go.mod if needed
go mod edit -version=v1.1.0

# Update documentation references
# docs/index.md, README.md, examples/
```

### 3. Update CHANGELOG.md
```markdown
## [1.1.0] - 2025-10-07

### Added
- New IAM resource for custom permissions
- Enhanced error handling for OAuth2 failures

### Changed  
- Improved performance for bulk operations
- Updated dependencies to latest versions

### Fixed
- Fixed issue with role binding validation
- Resolved timeout issues in large deployments

### Breaking Changes
- None

### Migration Guide
- No migration required for this release
```

### 4. Final Testing
```bash
# Run comprehensive tests
make test-all

# Validate release configuration
make validate-release-config

# Test local build
make test-build-local

# Run integration tests
make test-integration
```

### 5. Create and Push Release Tag
```bash
# Commit final changes
git add .
git commit -m "chore: prepare release v1.1.0

- Update CHANGELOG.md with release notes
- Update version references in documentation
- Final pre-release validation complete"

# Create annotated tag
git tag -a v1.1.0 -m "Release v1.1.0

New Features:
- Enhanced IAM resource management
- Improved error handling and user experience
- Performance optimizations for large deployments

Bug Fixes:
- Fixed role binding validation issues
- Resolved OAuth2 timeout problems

This release is fully backward compatible with v1.0.x.
See CHANGELOG.md for detailed changes."

# Push tag to trigger release
git push origin v1.1.0

# Push branch changes (if using release branch)
git push origin release/v1.1.0
```

### 6. Monitor Release Process

#### GitHub Actions Monitoring
- Navigate to: https://github.com/extenda/terraform-provider-hiiretail/actions
- Monitor the release workflow progress
- Verify all steps complete successfully

**Expected Timeline:**
- Workflow trigger: Immediate
- Build completion: 5-10 minutes
- GitHub release creation: 5-10 minutes
- Registry detection: 15-30 minutes
- Registry publication: 30-60 minutes

#### Artifact Verification
Check GitHub release includes:
- [ ] Multi-platform binary archives (.zip files)
- [ ] SHA256SUMS checksum file
- [ ] GPG signature file (.sig)
- [ ] Release notes from CHANGELOG.md
- [ ] Proper semantic version tag

### 7. Post-Release Verification

#### Registry Verification
- Visit: https://registry.terraform.io/providers/extenda/hiiretail
- Verify new version appears
- Check documentation updated
- Validate installation instructions

#### Installation Testing
```bash
# Test fresh installation
mkdir test-release && cd test-release
cat > main.tf << 'EOF'
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 1.1.0"
    }
  }
}
EOF

terraform init
terraform providers
```

#### Smoke Testing
```bash
# Run basic provider functionality
terraform plan
terraform apply
terraform destroy
```

## 🚨 Troubleshooting Guide

### GitHub Actions Failures

#### GPG Import Issues
**Symptoms:**
```
Error: gpg: can't connect to the agent: IPC connect call failed
```

**Solutions:**
1. Verify GPG_PRIVATE_KEY secret format:
   ```bash
   gpg --armor --export-secret-keys KEY_ID
   ```
2. Check PASSPHRASE secret is correct
3. Ensure key hasn't expired

#### Build Failures
**Symptoms:**
```
Error: failed to build for linux/amd64: exit status 1
```

**Solutions:**
1. Check Go module dependencies:
   ```bash
   go mod tidy
   go mod verify
   ```
2. Verify cross-compilation works locally:
   ```bash
   GOOS=linux GOARCH=amd64 go build
   ```
3. Review recent dependency changes

#### GoReleaser Configuration Errors
**Symptoms:**
```
Error: yaml: unmarshal errors
```

**Solutions:**
1. Validate configuration locally:
   ```bash
   goreleaser check
   ```
2. Test with dry run:
   ```bash
   goreleaser release --snapshot --skip-publish
   ```

### Registry Publication Issues

#### Provider Not Detected
**Cause:** Registry hasn't detected the GitHub release

**Solutions:**
1. Wait 30-60 minutes for automatic detection
2. Verify release has all required artifacts
3. Check release follows semantic versioning
4. Ensure repository is public

#### Documentation Missing
**Cause:** Documentation not generated properly

**Solutions:**
1. Verify docs/ directory structure:
   ```
   docs/
   ├── index.md
   ├── resources/
   └── data-sources/
   ```
2. Check resource documentation format
3. Validate examples/ directory exists

#### Installation Failures
**Cause:** Binary compatibility or signing issues

**Solutions:**
1. Verify GPG signatures:
   ```bash
   gpg --verify *.sig
   ```
2. Test multi-platform binaries
3. Check binary naming convention

### Emergency Procedures

#### Bad Release Published
1. **Immediate Actions:**
   - Add warning to README.md
   - Create GitHub issue describing problems
   - Communicate to users via appropriate channels

2. **Fix and Re-release:**
   - Fix issues on main branch
   - Create patch release (e.g., v1.1.1)
   - Follow standard release process

3. **Registry Management:**
   - Contact HashiCorp if release needs removal
   - Document known issues clearly

#### Critical Security Issue
1. **Immediate Response:**
   - Assess severity and impact
   - Coordinate with security team
   - Prepare security advisory

2. **Release Process:**
   - Create security patch
   - Fast-track release process
   - Notify users of security update

## 📊 Release Metrics and Monitoring

### Success Metrics
Track the following for each release:
- [ ] Release automation success rate
- [ ] Time from tag to registry publication
- [ ] Download statistics
- [ ] User feedback and issues

### Monitoring Dashboards
- **GitHub Actions**: Monitor workflow success rates
- **Registry Analytics**: Track downloads and adoption
- **Issue Tracker**: Monitor bug reports and feature requests

## 🔄 Continuous Improvement

### Post-Release Review
After each release, conduct a retrospective:
1. What went well?
2. What could be improved?
3. Any process updates needed?
4. Documentation gaps identified?

### Process Updates
- Update this documentation based on learnings
- Improve automation where possible
- Enhance testing coverage
- Streamline communication

## 📚 Additional Resources

### Documentation Links
- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Terraform Registry Requirements](https://www.terraform.io/docs/registry/providers/publishing.html)
- [Semantic Versioning](https://semver.org/)

### Internal Resources
- [GPG Setup Guide](./gpg-setup.md)
- [Registry Monitoring Guide](./registry-monitoring.md)
- [End-User Testing Guide](./end-user-testing.md)

### Emergency Contacts
- **Release Manager**: [Contact Information]
- **Security Team**: [Contact Information]  
- **Platform Team**: [Contact Information]

## 📋 Quick Reference Commands

### Local Development
```bash
# Validate configuration
make validate-release-config

# Test build locally
make test-build-local

# Run all tests
make test-all
```

### Release Commands
```bash
# Create release tag
git tag -a v1.1.0 -m "Release v1.1.0"

# Push tag (triggers release)
git push origin v1.1.0

# Check release status
gh release view v1.1.0
```

### Troubleshooting
```bash
# Check GoReleaser config
goreleaser check

# Test GPG signing
gpg --list-secret-keys

# Verify binaries
file dist/*/terraform-provider-*
```

## ✅ Release Checklist Summary

### Pre-Release (T-1 week)
- [ ] Plan release scope and version
- [ ] Review and merge all intended changes
- [ ] Update documentation and examples
- [ ] Coordinate with stakeholders

### Release Day (T-0)
- [ ] Final testing and validation
- [ ] Update CHANGELOG.md
- [ ] Create and push release tag
- [ ] Monitor GitHub Actions workflow
- [ ] Verify GitHub release creation

### Post-Release (T+1 hour)
- [ ] Verify Terraform Registry publication
- [ ] Test provider installation
- [ ] Monitor for issues
- [ ] Communicate release to users

### Follow-up (T+1 week)
- [ ] Review release metrics
- [ ] Address any issues found
- [ ] Update documentation based on feedback
- [ ] Plan next release if needed

---

**Remember**: The goal is reliable, automated releases that provide value to users while maintaining high quality and security standards.