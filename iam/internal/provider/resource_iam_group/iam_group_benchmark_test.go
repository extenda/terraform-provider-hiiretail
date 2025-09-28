package resource_iam_group

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// T032: Add performance benchmarks for Group operations

func BenchmarkGroupResource_Create(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	// Setup basic configuration for benchmarking
	r.baseURL = "https://test-api.example.com"
	r.tenantID = "benchmark-tenant"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data := IamGroupModel{
			Name:        types.StringValue("benchmark-group"),
			Description: types.StringValue("Benchmark test group"),
		}

		// Simulate the validation and processing overhead
		_ = r.validateGroupData(ctx, &data)

		// Note: Actual API calls would be mocked/stubbed in real benchmarks
		// to avoid network overhead and focus on local processing performance
	}
}

func BenchmarkGroupResource_Read(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)

	r.baseURL = "https://test-api.example.com"
	r.tenantID = "benchmark-tenant"

	data := IamGroupModel{
		Id:          types.StringValue("benchmark-group-id"),
		Name:        types.StringValue("benchmark-group"),
		Description: types.StringValue("Benchmark test group"),
		Status:      types.StringValue("active"),
		TenantId:    types.StringValue("benchmark-tenant"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate read operation processing
		_ = data.Id.ValueString()
		_ = data.Name.ValueString()
		_ = data.Description.ValueString()
	}
}

func BenchmarkGroupResource_Update(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	r.baseURL = "https://test-api.example.com"
	r.tenantID = "benchmark-tenant"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data := IamGroupModel{
			Id:          types.StringValue("benchmark-group-id"),
			Name:        types.StringValue("updated-benchmark-group"),
			Description: types.StringValue("Updated benchmark test group"),
		}

		// Simulate validation overhead
		_ = r.validateGroupData(ctx, &data)

		// Note: Actual API calls would be mocked in real benchmarks
	}
}

func BenchmarkGroupResource_Delete(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)

	r.baseURL = "https://test-api.example.com"
	r.tenantID = "benchmark-tenant"

	data := IamGroupModel{
		Id:          types.StringValue("benchmark-group-id"),
		Name:        types.StringValue("benchmark-group"),
		Description: types.StringValue("Benchmark test group"),
		Status:      types.StringValue("active"),
		TenantId:    types.StringValue("benchmark-tenant"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate delete operation processing
		_ = data.Id.ValueString()

		// Note: Actual API calls would be mocked in real benchmarks
	}
}

func BenchmarkGroupResource_Validation(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)
	ctx := context.Background()

	data := IamGroupModel{
		Name:        types.StringValue("test-group-validation"),
		Description: types.StringValue("Test group for validation benchmarking"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.validateGroupData(ctx, &data)
	}
}

func BenchmarkGroupResource_Schema(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = IamGroupResourceSchema(ctx)
	}
}

func BenchmarkGroupResource_ModelBinding(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data := IamGroupModel{
			Id:          types.StringValue("benchmark-id"),
			Name:        types.StringValue("benchmark-name"),
			Description: types.StringValue("benchmark-description"),
			Status:      types.StringValue("active"),
			TenantId:    types.StringValue("benchmark-tenant"),
		}

		// Simulate model operations
		_ = data.Id.ValueString()
		_ = data.Name.ValueString()
		_ = data.Description.ValueString()
		_ = data.Status.ValueString()
		_ = data.TenantId.ValueString()
	}
}

// Benchmark helper function performance
func BenchmarkHelperFunctions(b *testing.B) {
	r := NewIamGroupResource().(*IamGroupResource)

	b.Run("isRetryableError", func(b *testing.B) {
		err := fmt.Errorf("connection timeout")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = r.isRetryableError(err)
		}
	})

	b.Run("mapHTTPError", func(b *testing.B) {
		err := fmt.Errorf("test error")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = r.mapHTTPError(404, err)
		}
	})
}

// Example usage:
// go test -bench=. -benchmem ./internal/provider/resource_iam_group/
// go test -bench=BenchmarkGroupResource_Create -benchtime=10s ./internal/provider/resource_iam_group/
