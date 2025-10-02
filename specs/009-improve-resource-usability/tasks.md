# Tasks: Improve Resource Usability

**Input**: Design documents from `/specs/009-improve-resource-usability/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory ✅
   → Tech stack found: Go 1.19+ (Terraform Plugin Framework)
   → Extract: HashiCorp Terraform Plugin Framework, OAuth2 client libraries, HiiRetail IAM API clients
2. Load optional design documents: ✅
   → data-model.md: ValidationMessage, ResourceReference, ConfigurationExample, PermissionPattern
   → contracts/: validation-api.md, resource-schemas.md
   → research.md: Terraform provider best practices, validation patterns
3. Generate tasks by category: ✅
   → Setup: Enhanced error system, base validators
   → Tests: Contract tests for validation scenarios
   → Core: Resource schema enhancements, custom validators
   → Integration: Cross-resource validation, reference resolution
   → Polish: Documentation, examples, comprehensive testing
4. Apply task rules: ✅
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...) ✅
6. Generate dependency graph ✅
7. Create parallel execution examples ✅
8. Validate task completeness: ✅
   → All contracts have tests ✅
   → All entities have validators ✅
   → All resources enhanced ✅
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Based on plan.md structure (single Terraform provider project):
- **Provider resources**: `iam/internal/provider/resource_{type}/`
- **Validation code**: `iam/internal/validation/`
- **Tests**: `iam/tests/{acceptance,unit,validation}/`
- **Examples**: `iam/examples/{basic,enterprise,troubleshooting}/`
- **Documentation**: `iam/docs/{resources,guides}/`

## Phase 3.1: Foundation Setup
- [x] T001 Create enhanced validation package structure in `iam/internal/validation/`
- [x] T002 Create enhanced error system in `iam/internal/validation/errors.go`
- [x] T003 [P] Create base validator interfaces in `iam/internal/validation/interfaces.go`
- [x] T004 [P] Create validation test framework in `iam/tests/validation/framework.go`

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Validation Contract Tests
- [ ] T005 [P] Field validation contract tests in `iam/tests/validation/field_validation_test.go`
- [ ] T006 [P] Reference resolution contract tests in `iam/tests/validation/reference_validation_test.go`
- [ ] T007 [P] Permission validation contract tests in `iam/tests/validation/permission_validation_test.go`
- [ ] T008 [P] Enhanced error message tests in `iam/tests/validation/error_message_test.go`

### Resource-Specific Validation Tests
- [ ] T009 [P] Group resource validation tests in `iam/tests/validation/group_validation_test.go`
- [ ] T010 [P] Custom role resource validation tests in `iam/tests/validation/custom_role_validation_test.go`
- [ ] T011 [P] Role binding resource validation tests in `iam/tests/validation/role_binding_validation_test.go`

### Integration Scenario Tests
- [ ] T012 [P] New user experience integration test in `iam/tests/acceptance/new_user_test.go`
- [ ] T013 [P] Error recovery integration test in `iam/tests/acceptance/error_recovery_test.go`
- [ ] T014 [P] Complex configuration integration test in `iam/tests/acceptance/complex_config_test.go`

## Phase 3.3: Core Validation Implementation (ONLY after tests are failing)

### Base Validation Components
- [ ] T015 Implement ValidationMessage entity in `iam/internal/validation/message.go`
- [ ] T016 Implement ResourceReference entity in `iam/internal/validation/reference.go`
- [ ] T017 Implement PermissionPattern entity in `iam/internal/validation/permission.go`
- [ ] T018 Implement enhanced diagnostic converter in `iam/internal/validation/diagnostics.go`

### Custom Validators (Parallel Implementation)
- [ ] T019 [P] Group name uniqueness validator in `iam/internal/validation/group_validators.go`
- [ ] T020 [P] Permission format validator in `iam/internal/validation/permission_validators.go`
- [ ] T021 [P] Resource reference validator in `iam/internal/validation/reference_validators.go`
- [ ] T022 [P] Member format validator in `iam/internal/validation/member_validators.go`
- [ ] T023 [P] Condition expression validator in `iam/internal/validation/condition_validators.go`

### Resource Schema Enhancements (Parallel per Resource)
- [ ] T024 [P] Enhance Group resource schema in `iam/internal/provider/resource_iam_group/iam_group_resource_gen.go`
- [ ] T025 [P] Enhance Custom Role resource schema in `iam/internal/provider/resource_iam_custom_role/iam_custom_role_resource_gen.go`
- [ ] T026 [P] Enhance Role Binding resource schema in `iam/internal/provider/resource_iam_role_binding/iam_role_binding_resource_gen.go`

## Phase 3.4: Cross-Resource Integration
- [ ] T027 Implement reference resolution service in `iam/internal/validation/resolver.go`
- [ ] T028 Add plan-time validation for Group resource in `iam/internal/provider/resource_iam_group/`
- [ ] T029 Add plan-time validation for Custom Role resource in `iam/internal/provider/resource_iam_custom_role/`
- [ ] T030 Add plan-time validation for Role Binding resource in `iam/internal/provider/resource_iam_role_binding/`
- [ ] T031 Implement circular dependency detection in `iam/internal/validation/dependency_checker.go`
- [ ] T032 Add API connectivity validation in `iam/internal/validation/api_validator.go`

## Phase 3.5: Documentation and Examples (Parallel Creation)
- [ ] T033 [P] Create basic usage examples in `iam/examples/basic/`
- [ ] T034 [P] Create enterprise setup examples in `iam/examples/enterprise/`
- [ ] T035 [P] Create troubleshooting examples in `iam/examples/troubleshooting/`
- [ ] T036 [P] Create resource documentation in `iam/docs/resources/`
- [ ] T037 [P] Create validation guide in `iam/docs/guides/validation.md`
- [ ] T038 [P] Create troubleshooting guide in `iam/docs/guides/troubleshooting.md`

## Phase 3.6: Polish and Comprehensive Testing
- [ ] T039 [P] Unit tests for validation utilities in `iam/tests/unit/validation_utils_test.go`
- [ ] T040 [P] Performance tests for validation (<2s target) in `iam/tests/unit/validation_performance_test.go`
- [ ] T041 [P] Reference resolution performance tests in `iam/tests/unit/reference_performance_test.go`
- [ ] T042 Update provider documentation with new validation features
- [ ] T043 Create migration guide for enhanced validation
- [ ] T044 Run comprehensive quickstart validation scenarios

## Dependencies

### Sequential Dependencies
- **Foundation**: T001 → T002 → T003, T004
- **Tests First**: T005-T014 MUST complete before T015-T032
- **Base Components**: T015-T018 before custom validators T019-T023
- **Validators**: T019-T023 before schema enhancements T024-T026
- **Schema Enhancement**: T024-T026 before integration T027-T032
- **Integration**: T027-T032 before polish T039-T044

### Parallel Groups
```
Group A (Contract Tests): T005, T006, T007, T008
Group B (Resource Tests): T009, T010, T011  
Group C (Integration Tests): T012, T013, T014
Group D (Custom Validators): T019, T020, T021, T022, T023
Group E (Schema Enhancements): T024, T025, T026
Group F (Documentation): T033, T034, T035, T036, T037, T038
Group G (Polish Tests): T039, T040, T041
```

## Parallel Execution Examples

### Contract Tests (After T004)
```bash
# Launch all contract tests together:
Task: "Field validation contract tests in iam/tests/validation/field_validation_test.go"
Task: "Reference resolution contract tests in iam/tests/validation/reference_validation_test.go"
Task: "Permission validation contract tests in iam/tests/validation/permission_validation_test.go"
Task: "Enhanced error message tests in iam/tests/validation/error_message_test.go"
```

### Custom Validators (After T018)
```bash
# Launch all validator implementations together:
Task: "Group name uniqueness validator in iam/internal/validation/group_validators.go"
Task: "Permission format validator in iam/internal/validation/permission_validators.go"
Task: "Resource reference validator in iam/internal/validation/reference_validators.go"
Task: "Member format validator in iam/internal/validation/member_validators.go"
Task: "Condition expression validator in iam/internal/validation/condition_validators.go"
```

### Resource Schema Enhancements (After T023)
```bash
# Launch all resource schema updates together:
Task: "Enhance Group resource schema in iam/internal/provider/resource_iam_group/iam_group_resource_gen.go"
Task: "Enhance Custom Role resource schema in iam/internal/provider/resource_iam_custom_role/iam_custom_role_resource_gen.go"  
Task: "Enhance Role Binding resource schema in iam/internal/provider/resource_iam_role_binding/iam_role_binding_resource_gen.go"
```

### Documentation Creation (After T032)
```bash
# Launch all documentation tasks together:
Task: "Create basic usage examples in iam/examples/basic/"
Task: "Create enterprise setup examples in iam/examples/enterprise/"
Task: "Create troubleshooting examples in iam/examples/troubleshooting/"
Task: "Create resource documentation in iam/docs/resources/"
Task: "Create validation guide in iam/docs/guides/validation.md"
Task: "Create troubleshooting guide in iam/docs/guides/troubleshooting.md"
```

## Success Criteria Validation

### Error Message Quality (T005-T008, T015-T018)
- ✅ Error messages include specific field paths
- ✅ Current values displayed in error messages
- ✅ Expected formats clearly described  
- ✅ Working examples provided in error messages
- ✅ Actionable guidance included for resolution

### Validation Coverage (T009-T011, T019-T026)
- ✅ All required fields have appropriate validation
- ✅ Format validation for names, emails, permissions
- ✅ Cross-resource references validated
- ✅ Permission strings follow expected patterns
- ✅ Conditional expressions validated

### User Experience (T012-T014, T027-T032)
- ✅ Users understand errors without consulting documentation
- ✅ Suggestions provided for common typos
- ✅ Related valid options shown when validation fails
- ✅ Plan-time validation catches issues before apply

### Documentation Quality (T033-T038, T042-T043)
- ✅ Resource schema descriptions comprehensive
- ✅ Examples demonstrate real-world usage patterns
- ✅ Field descriptions explain business purpose
- ✅ Troubleshooting guidance accessible

## Performance Targets
- **Validation Response**: <2 seconds for complex configurations (T040)
- **Reference Resolution**: <5 seconds including API calls (T041)
- **Provider Load Time**: No significant impact on startup

## Risk Mitigation
- **Backward Compatibility**: All enhancements maintain existing API contracts
- **Feature Flags**: Enhanced validation can be disabled if issues arise
- **Rollback Plan**: T043 includes rollback documentation
- **Testing Coverage**: Comprehensive test suite ensures reliability

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing (Critical for T005-T014)
- Commit after each task completion
- Integration tests (T012-T014) validate end-to-end user scenarios
- Performance tests (T040-T041) ensure scalability requirements

## Task Generation Rules Applied
*Applied during main() execution*

1. **From Contracts**: ✅
   - validation-api.md → contract tests T005-T008
   - resource-schemas.md → schema enhancement tasks T024-T026
   
2. **From Data Model**: ✅
   - ValidationMessage → T015, enhanced error system T002
   - ResourceReference → T016, reference validator T021
   - PermissionPattern → T017, permission validator T020
   - ConfigurationExample → documentation tasks T033-T038
   
3. **From User Stories (Quickstart)**: ✅
   - New user experience → T012 integration test
   - Error recovery → T013 integration test  
   - Complex configuration → T014 integration test
   - Validation scenarios → T039-T044 comprehensive testing

4. **Ordering Applied**: ✅
   - Setup (T001-T004) → Tests (T005-T014) → Implementation (T015-T032) → Polish (T033-T044)
   - Dependencies prevent premature parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T005-T008)
- [x] All entities have implementation tasks (T015-T017)
- [x] All tests come before implementation (T005-T014 before T015-T032)
- [x] Parallel tasks truly independent ([P] tasks use different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Resource schemas from simple_test.tf fully covered (T024-T026)
- [x] Cross-resource validation addresses role binding complexity (T027-T032)

---

**Task Status**: READY FOR EXECUTION ✅  
**Total Tasks**: 44  
**Parallel Groups**: 7  
**Dependencies**: Properly ordered ✅  
**Coverage**: Complete feature implementation ✅