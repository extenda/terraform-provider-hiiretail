package resource_iam_group

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// T033: Add concurrent operation safety tests

func TestGroupResource_ConcurrentValidation(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	errors := make([]error, numGoroutines*numOperations)
	errorIndex := 0
	var errorMutex sync.Mutex

	// Test concurrent validation calls
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				data := IamGroupModel{
					Name:        types.StringValue(fmt.Sprintf("concurrent-group-%d-%d", goroutineID, j)),
					Description: types.StringValue("Concurrent validation test"),
				}

				err := r.validateGroupData(ctx, &data)

				errorMutex.Lock()
				errors[errorIndex] = err
				errorIndex++
				errorMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Check that all validations succeeded
	for i, err := range errors {
		if i >= errorIndex {
			break
		}
		assert.NoError(t, err, "Validation should succeed for concurrent operation %d", i)
	}
}

func TestGroupResource_ConcurrentErrorMapping(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	const numGoroutines = 5
	const numOperations = 50

	var wg sync.WaitGroup
	results := make([]error, numGoroutines*numOperations)
	resultIndex := 0
	var resultMutex sync.Mutex

	// Test concurrent HTTP error mapping
	statusCodes := []int{400, 401, 403, 404, 409, 500, 503}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				statusCode := statusCodes[j%len(statusCodes)]
				originalErr := fmt.Errorf("test error %d-%d", goroutineID, j)

				mappedErr := r.mapHTTPError(statusCode, originalErr)

				resultMutex.Lock()
				results[resultIndex] = mappedErr
				resultIndex++
				resultMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify all error mappings were successful
	for i := 0; i < resultIndex; i++ {
		assert.NotNil(t, results[i], "Error mapping should produce a result for operation %d", i)
		assert.Error(t, results[i], "Mapped result should be an error for operation %d", i)
	}
}

func TestGroupResource_ConcurrentRetryLogic(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	const numGoroutines = 3
	const numOperations = 20

	var wg sync.WaitGroup
	successCount := int64(0)
	var successMutex sync.Mutex

	// Test concurrent retry operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				operationCalled := false

				err := r.retryOperation(ctx, func() error {
					operationCalled = true
					// Simulate success after some operations
					if (goroutineID+j)%3 == 0 {
						return nil
					}
					return fmt.Errorf("simulated transient error")
				})

				if err == nil && operationCalled {
					successMutex.Lock()
					successCount++
					successMutex.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify that some operations succeeded
	assert.Greater(t, successCount, int64(0), "Some retry operations should have succeeded")
}

func TestGroupResource_ConcurrentHelperFunctions(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// Test concurrent isRetryableError calls
	t.Run("isRetryableError", func(t *testing.T) {
		results := make([]bool, numGoroutines*numOperations)
		resultIndex := 0
		var resultMutex sync.Mutex

		errors := []error{
			fmt.Errorf("timeout occurred"),
			fmt.Errorf("connection refused"),
			fmt.Errorf("service temporarily unavailable"),
			fmt.Errorf("server error"),
			fmt.Errorf("not a retryable error"),
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					err := errors[j%len(errors)]
					result := r.isRetryableError(err)

					resultMutex.Lock()
					results[resultIndex] = result
					resultIndex++
					resultMutex.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Verify all calls completed successfully
		assert.Equal(t, numGoroutines*numOperations, resultIndex, "All isRetryableError calls should complete")
	})
}

func TestGroupResource_ConcurrentSchemaAccess(t *testing.T) {
	ctx := context.Background()

	const numGoroutines = 5
	const numOperations = 50

	var wg sync.WaitGroup

	// Test concurrent schema access
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				schema := IamGroupResourceSchema(ctx)

				// Verify schema has expected attributes
				assert.NotNil(t, schema.Attributes, "Schema should have attributes")
				assert.Contains(t, schema.Attributes, "name", "Schema should contain name attribute")
				assert.Contains(t, schema.Attributes, "id", "Schema should contain id attribute")
				assert.Contains(t, schema.Attributes, "description", "Schema should contain description attribute")
				assert.Contains(t, schema.Attributes, "status", "Schema should contain status attribute")
				assert.Contains(t, schema.Attributes, "tenant_id", "Schema should contain tenant_id attribute")
			}
		}()
	}

	wg.Wait()
}

func TestGroupResource_ConcurrentModelOperations(t *testing.T) {
	const numGoroutines = 8
	const numOperations = 100

	var wg sync.WaitGroup

	// Test concurrent model field access
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				data := IamGroupModel{
					Id:          types.StringValue(fmt.Sprintf("concurrent-id-%d-%d", goroutineID, j)),
					Name:        types.StringValue(fmt.Sprintf("concurrent-group-%d-%d", goroutineID, j)),
					Description: types.StringValue("Concurrent model test"),
					Status:      types.StringValue("active"),
					TenantId:    types.StringValue("concurrent-tenant"),
				}

				// Access all fields concurrently
				_ = data.Id.ValueString()
				_ = data.Name.ValueString()
				_ = data.Description.ValueString()
				_ = data.Status.ValueString()
				_ = data.TenantId.ValueString()
			}
		}(i)
	}

	wg.Wait()
}

// TestGroupResource_RaceConditionDetection tests for potential race conditions
func TestGroupResource_RaceConditionDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	// Configure the resource
	r.baseURL = "https://test-api.example.com"
	r.tenantID = "race-test-tenant"

	const numGoroutines = 20
	const duration = 2 * time.Second

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Start multiple goroutines performing different operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-stop:
					return
				default:
					// Perform various operations that might race
					data := IamGroupModel{
						Name:        types.StringValue(fmt.Sprintf("race-group-%d", id)),
						Description: types.StringValue("Race condition test"),
					}

					_ = r.validateGroupData(ctx, &data)
					_ = r.isRetryableError(fmt.Errorf("test error"))
					_ = r.mapHTTPError(404, fmt.Errorf("not found"))
				}
			}
		}(i)
	}

	// Let it run for a while
	time.Sleep(duration)
	close(stop)
	wg.Wait()

	// If we reach here without data races, the test passes
	t.Log("Race condition test completed without detected races")
}

// Example usage:
// go test -race -v ./internal/provider/resource_iam_group/ -run TestGroupResource_Concurrent
// go test -race -v ./internal/provider/resource_iam_group/ -run TestGroupResource_RaceConditionDetection
