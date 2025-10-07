# Quickstart: Provider Distribution Validation

## Overview
This quickstart guide validates the complete provider distribution pipeline from code changes to public availability via Terraform Registry.

## Prerequisites
- [x] GitHub repository with admin access
- [x] GPG key pair generated and configured
- [x] GitHub repository secrets configured:
  - `GPG_PRIVATE_KEY`: ASCII-armored private key
  - `PASSPHRASE`: GPG key passphrase
- [x] Clean repository state (no uncommitted changes)

## Phase 1: Configuration Setup

### Step 1: Configure GoReleaser
```bash
# Create .goreleaser.yml in repository root
cat > .goreleaser.yml << 'EOF'
project_name: terraform-provider-hiiretail

before:
  hooks:
    - go mod tidy

builds:
  - binary: terraform-provider-{{ .ProjectName }}_v{{ .Version }}
    main: ./
    goos:
      - linux
      - darwin  
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{ .Version }}
      - -X main.commit={{ .ShortCommit }}

archives:
  - format_overrides:
      - goos: windows
        format: zip
    name_template: "{{ .ProjectName }}_v{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "{{ .ProjectName }}_v{{ .Version }}_SHA256SUMS"

signs:
  - artifacts: checksum
    args: ["--batch", "--local-user", "{{ .Env.GPG_FINGERPRINT }}", "--output", "${signature}", "--detach-sign", "${artifact}"]

release:
  draft: false
  prerelease: auto
EOF
```

### Step 2: Create GitHub Actions Workflow
```bash
# Create .github/workflows/release.yml
mkdir -p .github/workflows
cat > .github/workflows/release.yml << 'EOF'
name: release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  id-token: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - name: Import GPG key
        id: import_gpg
        uses: crazy-max/ghaction-import-gpg@v6
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.PASSPHRASE }}
          
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
EOF
```

## Phase 2: Local Validation

### Step 3: Test GoReleaser Configuration
```bash
# Install GoReleaser locally for testing
brew install goreleaser

# Validate configuration without releasing
goreleaser check

# Test build process (creates dist/ directory)
goreleaser build --snapshot --clean

# Verify build artifacts
ls -la dist/
```

**Expected Output**:
```
dist/
├── terraform-provider-hiiretail_v0.0.0-next_darwin_amd64/
├── terraform-provider-hiiretail_v0.0.0-next_darwin_arm64/
├── terraform-provider-hiiretail_v0.0.0-next_linux_amd64/
├── terraform-provider-hiiretail_v0.0.0-next_linux_arm64/
└── terraform-provider-hiiretail_v0.0.0-next_windows_amd64/
```

### Step 4: Validate Binary Functionality
```bash
# Test a built binary
./dist/terraform-provider-hiiretail_v0.0.0-next_linux_amd64/terraform-provider-hiiretail_v0.0.0-next

# Should output provider version and help information
```

## Phase 3: Release Process Testing

### Step 5: Create Test Release
```bash
# Ensure clean repository state
git status

# Create and push a test tag
git tag v0.1.0-beta.1
git push origin v0.1.0-beta.1
```

### Step 6: Monitor GitHub Actions
1. Navigate to GitHub repository → Actions tab
2. Verify "release" workflow is triggered
3. Monitor workflow execution logs
4. Check for successful completion

**Expected Workflow Steps**:
- [x] Checkout repository
- [x] Set up Go 1.21
- [x] Import GPG key
- [x] Run GoReleaser
- [x] Create GitHub release
- [x] Upload artifacts

### Step 7: Validate Release Artifacts
Navigate to GitHub repository → Releases and verify:

**Release Assets**:
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_linux_amd64.tar.gz`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_linux_arm64.tar.gz`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_darwin_amd64.tar.gz`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_darwin_arm64.tar.gz`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_windows_amd64.zip`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_SHA256SUMS`
- [x] `terraform-provider-hiiretail_v0.1.0-beta.1_SHA256SUMS.sig`

## Phase 4: Terraform Registry Validation

### Step 8: Manual Installation Test
```bash
# Create test Terraform configuration
mkdir test-installation
cd test-installation

cat > main.tf << 'EOF'
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 0.1.0-beta.1"
    }
  }
}

provider "hiiretail" {
  # Configuration will be added after registry publication
}
EOF

# Initialize and verify provider download
terraform init
```

### Step 9: Registry Publication Verification
1. Wait for Terraform Registry to detect the release (may take several hours)
2. Search for provider at: `https://registry.terraform.io/providers/extenda/hiiretail`
3. Verify provider page shows correct version and documentation
4. Test installation via `terraform init` with registry source

## Phase 5: Production Release

### Step 10: Create Production Release
```bash
# After successful testing, create production release
git tag v1.0.0
git push origin v1.0.0
```

### Step 11: Final Validation
1. Monitor GitHub Actions workflow for v1.0.0
2. Verify GitHub release is created successfully  
3. Test provider installation from Terraform Registry
4. Validate provider functionality with real configuration

## Success Criteria
- [x] GoReleaser configuration valid and tested
- [x] GitHub Actions workflow executes successfully
- [x] Multi-platform binaries generated and signed
- [x] GitHub releases created with all required artifacts
- [x] Terraform Registry detects and publishes provider
- [x] Users can discover and install provider via `terraform init`
- [x] Provider functions correctly after installation

## Troubleshooting

### Common Issues and Solutions

**GPG Signing Fails**:
```bash
# Verify GPG key setup
gpg --list-secret-keys
# Ensure secrets are correctly configured in GitHub
```

**GoReleaser Build Fails**:
```bash
# Check Go module validity
go mod verify
go test ./...
```

**Registry Detection Delayed**:
- Registry scans for new releases periodically
- Publishing may take 1-24 hours
- Contact HashiCorp support if detection fails after 48 hours

## Next Steps
After successful validation:
1. Document release process for team members
2. Set up automated changelog generation
3. Configure release notifications
4. Plan regular release cadence