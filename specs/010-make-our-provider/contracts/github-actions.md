# GitHub Actions Workflow Contract

## Workflow Trigger Contract

```yaml
# Contract: Workflow triggers on semantic version tags
name: release
on:
  push:
    tags:
      - 'v*'

# Expected: Tag format v1.0.0, v2.1.3, v1.0.0-beta.1
# Rejected: latest, main, feature-branch, non-semantic tags
```

## Job Execution Contract

```yaml
# Contract: Release job with required environment and permissions
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    permissions:
      contents: write  # Required for release creation
      id-token: write  # Required for attestation
    
    steps:
      # Contract: Go version must match project requirements
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      # Contract: GPG setup for artifact signing
      - uses: crazy-max/ghaction-import-gpg@v6
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.PASSPHRASE }}
      
      # Contract: GoReleaser execution with required environment
      - uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
```

## Secrets Contract

```yaml
# Required GitHub Repository Secrets
secrets:
  GPG_PRIVATE_KEY:
    description: "ASCII-armored GPG private key for signing"
    required: true
    format: "-----BEGIN PGP PRIVATE KEY BLOCK-----\n...\n-----END PGP PRIVATE KEY BLOCK-----"
  
  PASSPHRASE:
    description: "GPG private key passphrase"
    required: true
    format: "string"
  
  GITHUB_TOKEN:
    description: "Automatic token for GitHub API access"
    required: true
    auto_generated: true
```

## Output Contract

```yaml
# Expected GitHub Release Artifacts
release_artifacts:
  binaries:
    - terraform-provider-hiiretail_v{version}_linux_amd64.zip
    - terraform-provider-hiiretail_v{version}_linux_arm64.zip
    - terraform-provider-hiiretail_v{version}_darwin_amd64.zip
    - terraform-provider-hiiretail_v{version}_darwin_arm64.zip
    - terraform-provider-hiiretail_v{version}_windows_amd64.zip
  
  checksums:
    - terraform-provider-hiiretail_v{version}_SHA256SUMS
    - terraform-provider-hiiretail_v{version}_SHA256SUMS.sig
  
  metadata:
    - release notes from tag annotation or CHANGELOG.md
    - semantic version tag (v1.0.0)
    - publication timestamp
```