# Quickstart: Fix Remaining Tests

## Steps to Validate All Resource and Provider Tests

1. Ensure all resource models and provider schemas include required fields (especially tenant_id)
2. Run `go test ./internal/provider/resource_iam_role` and fix any build or test errors
3. Run `go test ./internal/provider/resource_iam_role_binding` and fix argument and field mismatches
4. Run `go test ./internal/provider/testutils` and ensure all utility tests pass
5. Run `go test ./internal/validation` and ensure all validation logic is covered
6. Run `go test ./internal/provider/shared/validators` and ensure all shared validation logic is covered
7. Run `go test ./internal/provider/resource_iam_group` and fix any contract or quickstart validation errors
8. Run `go test ./internal/provider/resource_iam_resource` and fix panic and contract test errors
9. Validate that all contract and quickstart tests either pass or are skipped (not failed)
10. Confirm all fixes align with the HiiRetail Provider Constitution and do not introduce schema or security violations

## Success Criteria
- All specified test commands succeed without build or runtime errors
- No panics or nil pointer dereferences in any test
- All resource and provider models are schema-aligned and validated
- All quickstart and contract tests are either implemented or skipped

---
