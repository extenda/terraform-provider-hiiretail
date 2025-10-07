# GoReleaser Configuration Contract

## Build Configuration Contract

```yaml
# Contract: Project metadata and build settings
project_name: terraform-provider-hiiretail

before:
  hooks:
    - go mod tidy

builds:
  - binary: terraform-provider-{{ .ProjectName }}_v{{ .Version }}
    main: ./
    
    # Contract: Support all major platforms for Terraform providers
    goos:
      - linux
      - darwin
      - windows
    
    goarch:
      - amd64
      - arm64
    
    # Contract: Build flags for production releases
    ldflags:
      - -s -w
      - -X main.version={{ .Version }}
      - -X main.commit={{ .ShortCommit }}
    
    # Contract: Exclude unsupported platform combinations
    ignore:
      - goos: windows
        goarch: arm64
```

## Archive Configuration Contract

```yaml
# Contract: Archive format per platform
archives:
  - format: zip
    name_template: "{{ .ProjectName }}_v{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    
    # Contract: Platform-specific archive formats
    format_overrides:
      - goos: linux
        format: tar.gz
      - goos: darwin
        format: tar.gz
      - goos: windows
        format: zip
```

## Checksum and Signing Contract

```yaml
# Contract: Checksum generation for all artifacts
checksum:
  name_template: "{{ .ProjectName }}_v{{ .Version }}_SHA256SUMS"
  algorithm: sha256

# Contract: GPG signing for security verification
signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"
```

## Release Configuration Contract

```yaml
# Contract: GitHub release settings
release:
  github:
    owner: extenda
    name: terraform-provider-hiiretail
  
  # Contract: Release naming and content
  name_template: "v{{ .Version }}"
  draft: false
  prerelease: auto
  
  # Contract: Include changelog in release notes
  footer: |
    ## Installation
    
    See the [documentation](https://registry.terraform.io/providers/extenda/hiiretail/{{ .Version }}/docs) for installation instructions.
```

## Validation Contract

```yaml
# Contract: Pre-release validation requirements
validation:
  # Require clean git state
  - git status --porcelain | wc -l | grep -q "^0$"
  
  # Require semantic version tag
  - echo "{{ .Tag }}" | grep -qE "^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.-]+)?$"
  
  # Require Go module validity
  - go mod verify
  
  # Require test passing
  - go test ./...
```

## Output Guarantees

```yaml
# Contract: Guaranteed output artifacts
artifacts_produced:
  per_platform:
    - binary: terraform-provider-hiiretail_v{version}_{os}_{arch}
    - archive: terraform-provider-hiiretail_v{version}_{os}_{arch}.{ext}
  
  global:
    - checksums: terraform-provider-hiiretail_v{version}_SHA256SUMS
    - signature: terraform-provider-hiiretail_v{version}_SHA256SUMS.sig
    - release: GitHub release with all artifacts attached

# Contract: Artifact naming consistency
naming_pattern: "terraform-provider-hiiretail_v{major}.{minor}.{patch}_{os}_{arch}.{extension}"
```