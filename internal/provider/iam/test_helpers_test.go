package iam

import "context"

// NewServiceForTest creates a Service using provided raw client and clientService implementations.
// This helper is only built during `go test` (file ends with _test.go) and is intended for use by
// package-external tests to avoid constructing production *client.Client instances.
// Basic context helper to satisfy some tests if needed.
func Ctx() context.Context { return context.Background() }
