# Project Directory Tree: terraform-provider-hiiretail

> Every folder and file is annotated with its purpose.

```
.                                   # Project root: Terraform provider for HiiRetail
├── .github                         # GitHub configuration (CI/CD, prompts, workflows)
│   ├── copilot-instructions.md     # Instructions for GitHub Copilot agent
│   ├── prompts/                    # Prompt templates for agent-driven workflows
│   └── workflows/                  # GitHub Actions workflows (CI, release)
├── .gitignore                      # Git ignore rules
├── .golangci.yml                   # GolangCI-Lint configuration
├── .goreleaser.yml                 # GoReleaser config for building/distributing provider
├── .specify                        # Specify agent config, scripts, and templates
│   ├── memory/constitution.md      # Project constitution (governance, principles)
│   ├── scripts/bash/               # Bash scripts for feature/spec/plan automation
│   └── templates/                  # Templates for specs, plans, tasks, agent files
├── .vscode/settings.json           # VS Code workspace settings
├── AUTHENTICATION_ENHANCEMENT.md   # Documentation for authentication improvements
├── CHANGELOG.md                    # Project changelog
├── cleanup.sh                      # Shell script for cleaning up artifacts
├── dist/                           # Distribution artifacts (built binaries, metadata)
│   └── ...                         # Platform-specific provider binaries and metadata
├── docs/                           # Project documentation
│   ├── data-sources/               # Docs for Terraform data sources
│   ├── end-user-testing.md         # End-user testing guide
│   ├── examples/                   # Example Terraform configs
│   ├── guides/                     # Guides (authentication, getting started)
│   ├── index.md                    # Main documentation index
│   ├── maintainer-release-process.md # Release process for maintainers
│   ├── registry-monitoring.md      # Registry monitoring guide
│   └── resources/                  # Docs for Terraform resources
├── examples/                       # Example Terraform configurations
│   ├── provider/                   # Example provider block
│   └── resources/                  # Resource-specific examples
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── iam/                            # IAM provider development, scripts, tests, docs
│   ├── .githooks/                  # Git hooks (pre-commit)
│   ├── .github/prompts/            # Prompt templates (copied from root)
│   ├── .specify/                   # Specify agent config (copied from root)
│   ├── acceptance_tests/           # Acceptance test files for IAM resources
│   ├── demo/                       # Demo code and OAuth2 example
│   ├── dev/                        # Dev environment files
│   ├── scripts/                    # Shell scripts for testing, validation
│   ├── specs/                      # Feature specs and tasks for IAM
│   ├── templates/                  # Documentation templates
│   ├── tests/                      # Test suites (acceptance, unit, contract, integration, performance, validation)
│   └── ...                         # Misc. docs, configs, and test scripts
├── internal/                       # Provider implementation code
│   ├── provider/                   # Main provider code
│   │   ├── iam/                    # IAM API client, resource/data source logic
│   │   │   ├── datasources/        # Data source implementations (groups, roles)
│   │   │   ├── resources/          # Resource implementations (custom_role, group, resource)
│   │   │   └── service.go          # IAM service client (CRUD, models, API logic)
│   │   ├── provider_hiiretail_iam/ # Auto-generated provider schema/model (do not edit)
│   │   │   └── hiiretail_iam_provider_gen.go
│   │   ├── provider_integration_test.go # Provider integration tests
│   │   ├── provider_oauth2_test.go      # OAuth2 authentication tests
│   │   ├── provider_test.go             # Provider unit tests
│   │   ├── provider.go                  # Provider entrypoint and registration
│   │   ├── resource_iam_custom_role/    # Custom role resource code and tests
│   │   ├── resource_iam_group/          # Group resource code and tests
│   │   ├── resource_iam_resource/       # Resource code and tests
│   │   ├── resource_iam_role/           # Role resource code and tests
│   │   ├── resource_iam_role_binding/   # Role binding resource code and tests
│   │   ├── shared/                      # Shared code (auth, client, validators)
│   │   └── testutils/                   # Test utilities (mock server)
│   └── validation/                      # Validation logic for provider schemas
├── LICENSE                             # Project license
├── main.go                             # Main entrypoint for provider binary
├── Makefile                            # Build and test automation
├── PROJECT_CONSTITUTION.md              # Project governance and principles
├── README.md                            # Project overview and documentation
├── REPOSITORY_CLEANUP_REPORT.md         # Report on repository cleanup
├── sonar-project.properties             # SonarCloud static analysis config
├── specs/                               # Feature specifications and plans
│   └── [feature folders]                # Each feature: contracts, data-model, plan, quickstart, research, spec, tasks
├── terraform-provider-hiiretail         # Provider binary (legacy or build artifact)
├── terraform-registry-manifest.json     # Terraform Registry manifest
├── test_build_artifacts.sh              # Shell script for build artifact validation
├── test_github_actions.sh               # Shell script for GitHub Actions validation
├── test_goreleaser_config.sh            # Shell script for GoReleaser config validation
└── test_release_process.sh              # Shell script for release process validation
```

---

## Folder/File Purpose Summary

- **.github/**: CI/CD, workflow, and agent prompt configuration.
- **.specify/**: Agent automation, governance, and template files.
- **dist/**: Built provider binaries and metadata for distribution.
- **docs/**: All documentation for users, maintainers, and registry.
- **examples/**: Example Terraform configurations for users.
- **iam/**: IAM provider development, tests, scripts, and documentation.
- **internal/provider/**: Main provider implementation, including IAM API client, resources, data sources, shared code, and tests.
- **LICENSE, README.md, PROJECT_CONSTITUTION.md**: Legal, documentation, and governance.
- **specs/**: Feature specifications, plans, and supporting docs.
- **terraform-provider-hiiretail, terraform-registry-manifest.json**: Provider binary and registry manifest.
- **test_*.sh**: Shell scripts for validating build, release, and CI/CD processes.

If you need a comment on any specific file or folder, let me know!
