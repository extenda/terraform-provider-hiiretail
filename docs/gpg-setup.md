# GPG Setup for Terraform Provider Releases

This document outlines the GPG key configuration used for signing Terraform provider releases.

## GPG Key Information

- **Key ID**: F9549AA602E3C9DC
- **Full Fingerprint**: E8C29730E7CEBC4DB6294298F9549AA602E3C9DC
- **Owner**: Shayne Clausson <shayne.clausson@extendaretail.com>
- **Algorithm**: RSA 4096-bit
- **Created**: 2025-10-07

## GitHub Secrets Configuration

The following secrets must be configured in the GitHub repository:

### Required Secrets

1. **GPG_PRIVATE_KEY**
   - The private GPG key in ASCII-armored format
   - Export with: `gpg --armor --export-secret-keys E8C29730E7CEBC4DB6294298F9549AA602E3C9DC`

2. **PASSPHRASE**
   - The passphrase for the GPG private key
   - Used by GitHub Actions to unlock the key for signing

### How GitHub Actions Uses GPG

1. The `crazy-max/ghaction-import-gpg@v6` action imports the private key
2. The action automatically detects the key fingerprint
3. GoReleaser uses the fingerprint via `GPG_FINGERPRINT` environment variable
4. All release artifacts are signed with the imported key

## Local Development

For local testing of GPG signatures:

```bash
# Test GPG signing manually
export GPG_TTY=$(tty)
export GPG_FINGERPRINT="E8C29730E7CEBC4DB6294298F9549AA602E3C9DC"

# Test signing a file
gpg --batch --local-user $GPG_FINGERPRINT --output test.sig --detach-sign somefile

# Verify signature
gpg --verify test.sig somefile
```

## Security Notes

- The GPG key is used exclusively for signing Terraform provider releases
- Private key is stored securely in GitHub repository secrets
- Key fingerprint is publicly available and used for verification
- All signatures can be verified using the public key
- The key has a long expiration date (2067-12-21) for release stability

## Verification

Users can verify release signatures using the public key:

```bash
# Import the public key
gpg --keyserver keyserver.ubuntu.com --recv-keys E8C29730E7CEBC4DB6294298F9549AA602E3C9DC

# Verify a release signature
gpg --verify terraform-provider-hiiretail_1.0.0_SHA256SUMS.sig terraform-provider-hiiretail_1.0.0_SHA256SUMS
```