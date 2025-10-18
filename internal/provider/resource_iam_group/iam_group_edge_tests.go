package resource_iam_group

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_retryOperation_exhausts_and_returns_error(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// operation that always returns a server error (retryable)
	calls := 0
	err := r.retryOperation(context.Background(), func() error {
		calls++
		return errors.New("server error: status 502")
	})
	require.Error(t, err)
	require.Equal(t, 3, calls)
}

func Test_isRetryableError_false_cases(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	require.False(t, r.isRetryableError(nil))
	require.False(t, r.isRetryableError(errors.New("not retryable")))
}

func Test_contains_empty_substr_behavior(t *testing.T) {
	// empty substring should be false per implementation because len(s) >= len(substr) check
	if contains("", "") {
		t.Fatalf("expected contains(empty, empty) to be false under this implementation")
	}
}
