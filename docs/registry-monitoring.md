# Terraform Registry Monitoring Guide

This guide explains how to monitor the provider publication process to the Terraform Registry.

## Release Process Timeline

### 1. GitHub Release (Automated)
When a version tag is pushed (e.g., `v0.1.0-beta`):
- GitHub Actions workflow triggers automatically
- GoReleaser builds multi-platform binaries
- Artifacts are signed with GPG
- GitHub release is created with artifacts

**Timeline**: 5-10 minutes after tag push

### 2. Terraform Registry Detection (Automatic)
The Terraform Registry monitors GitHub releases:
- Detects new releases via GitHub API
- Validates provider structure and artifacts
- Imports provider metadata and documentation

**Timeline**: 15-30 minutes after GitHub release

### 3. Provider Publication (Automatic)
Once validated, the provider becomes available:
- Listed on [registry.terraform.io](https://registry.terraform.io/providers/extenda/hiiretail)
- Available for `terraform init` installation
- Documentation published automatically

**Timeline**: 30-60 minutes after GitHub release

## Monitoring Checklist

### ✅ Phase 1: GitHub Actions Workflow
Monitor at: https://github.com/extenda/terraform-provider-hiiretail/actions

**What to check:**
- [ ] Workflow triggered by tag push
- [ ] Go setup completed successfully  
- [ ] GPG key import successful
- [ ] GoReleaser execution completed
- [ ] Multi-platform binaries generated
- [ ] GPG signatures created
- [ ] GitHub release published

**Expected artifacts:**
```
terraform-provider-hiiretail_v0.1.0-beta_linux_amd64.zip
terraform-provider-hiiretail_v0.1.0-beta_linux_arm64.zip
terraform-provider-hiiretail_v0.1.0-beta_darwin_amd64.zip
terraform-provider-hiiretail_v0.1.0-beta_darwin_arm64.zip
terraform-provider-hiiretail_v0.1.0-beta_windows_amd64.zip
terraform-provider-hiiretail_v0.1.0-beta_windows_arm64.zip
terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS
terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS.sig
```

### ✅ Phase 2: GitHub Release
Monitor at: https://github.com/extenda/terraform-provider-hiiretail/releases

**What to check:**
- [ ] Release created with correct tag (v0.1.0-beta)
- [ ] All platform archives present
- [ ] SHA256SUMS file present
- [ ] GPG signature (.sig) file present
- [ ] Release notes populated from CHANGELOG.md
- [ ] Pre-release flag set correctly for beta versions

### ✅ Phase 3: Terraform Registry
Monitor at: https://registry.terraform.io/providers/extenda/hiiretail

**What to check:**
- [ ] Provider listing appears
- [ ] Version v0.1.0-beta is available
- [ ] Documentation generated correctly
- [ ] Installation instructions accurate
- [ ] Provider metadata correct (description, homepage, source)

## Troubleshooting Common Issues

### GitHub Actions Failures

**GPG Import Failure**
```
Error: gpg: can't connect to the agent: IPC connect call failed
```
- Check GPG_PRIVATE_KEY secret format
- Verify PASSPHRASE secret is correct
- Ensure secrets are not expired

**Build Failure**
```
Error: failed to build: exit status 1
```
- Check Go module configuration
- Verify dependencies are available
- Review go.mod and go.sum files

**GoReleaser Configuration Error**
```
Error: yaml: unmarshal errors
```
- Validate .goreleaser.yml syntax
- Run `goreleaser check` locally
- Verify all required fields present

### Terraform Registry Issues

**Provider Not Detected**
- Verify GitHub release has all required artifacts
- Check release follows semantic versioning (v1.0.0)
- Ensure repository is public
- Wait 30-60 minutes for detection

**Documentation Missing**
- Verify docs/ directory structure
- Check resource and data source documentation
- Ensure examples/ directory exists
- Review provider metadata in main.go

**Installation Failures**
- Verify binary naming matches Terraform expectations
- Check GPG signatures are valid
- Ensure multi-platform builds succeeded

## Manual Verification Commands

### Test Provider Installation
```bash
# Create test Terraform configuration
mkdir test-provider && cd test-provider

cat > main.tf << 'EOF'
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 0.1.12-beta"
    }
  }
}

provider "hiiretail" {
  # Configuration
}
EOF

# Test installation
terraform init
terraform providers
```

### Verify GPG Signatures
```bash
# Download release artifacts
curl -LO https://github.com/extenda/terraform-provider-hiiretail/releases/download/v0.1.0-beta/terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS
curl -LO https://github.com/extenda/terraform-provider-hiiretail/releases/download/v0.1.0-beta/terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS.sig

# Import public key and verify
gpg --keyserver keyserver.ubuntu.com --recv-keys E8C29730E7CEBC4DB6294298F9549AA602E3C9DC
gpg --verify terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS.sig terraform-provider-hiiretail_v0.1.0-beta_SHA256SUMS
```

## Success Criteria

### ✅ Complete Success
- [ ] GitHub release published automatically
- [ ] All artifacts present and signed
- [ ] Provider available on Terraform Registry
- [ ] Documentation published correctly
- [ ] Provider installs successfully with `terraform init`
- [ ] All platform binaries work correctly
- [ ] GPG signatures verify successfully

### 📋 Registry Publication Checklist
1. **GitHub Release**: Check https://github.com/extenda/terraform-provider-hiiretail/releases/tag/v0.1.0-beta
2. **Registry Listing**: Visit https://registry.terraform.io/providers/extenda/hiiretail/0.1.0-beta
3. **Installation Test**: Run `terraform init` with provider configuration
4. **Documentation**: Verify docs display correctly on registry
5. **Multi-platform**: Test binaries work on different platforms

## Next Steps After Registry Publication

1. **Create Stable Release**: Tag v1.0.0 for first stable release
2. **Monitor Usage**: Track downloads and adoption
3. **Collect Feedback**: Gather user feedback for improvements
4. **Continuous Updates**: Regular releases with bug fixes and features
5. **Community Engagement**: Support users and respond to issues