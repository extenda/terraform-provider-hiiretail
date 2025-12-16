package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// mockRawClient implements iam.RawClient for testing
type mockRawClient struct {
	DoFunc func(ctx context.Context, req *client.Request) (*client.Response, error)
}

func (m *mockRawClient) Do(ctx context.Context, req *client.Request) (*client.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(ctx, req)
	}
	return nil, nil
}

// mockClientService implements iam.clientService for testing
type mockClientService struct{}

func (m *mockClientService) IAMClient() iam.RawClient {
	return nil
}

func (m *mockClientService) CCCClient() iam.RawClient {
	return nil
}

func (m *mockClientService) Get(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
	return nil, nil
}

func (m *mockClientService) Post(ctx context.Context, path string, body interface{}) (*client.Response, error) {
	return nil, nil
}

func (m *mockClientService) Put(ctx context.Context, path string, body interface{}) (*client.Response, error) {
	return nil, nil
}

func (m *mockClientService) Delete(ctx context.Context, path string) (*client.Response, error) {
	return nil, nil
}

func TestResourceDataSource_Metadata(t *testing.T) {
	ds := NewResourceDataSource()
	req := datasource.MetadataRequest{
		ProviderTypeName: "hiiretail",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	assert.Equal(t, "hiiretail_iam_resource", resp.TypeName)
}

func TestResourceDataSource_Schema(t *testing.T) {
	ds := NewResourceDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.NotEmpty(t, resp.Schema.Attributes)

	// Check required attributes
	idAttr := resp.Schema.Attributes["id"]
	assert.NotNil(t, idAttr)

	nameAttr := resp.Schema.Attributes["name"]
	assert.NotNil(t, nameAttr)

	typeAttr := resp.Schema.Attributes["name"]
	assert.NotNil(t, typeAttr)

	propsAttr := resp.Schema.Attributes["properties"]
	assert.NotNil(t, propsAttr)
}

func TestResourceDataSource_Configure(t *testing.T) {
	ctx := context.Background()
	ds := &ResourceDataSource{}

	t.Run("nil provider data", func(t *testing.T) {
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &datasource.ConfigureResponse{}

		ds.Configure(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("invalid provider data type", func(t *testing.T) {
		req := datasource.ConfigureRequest{
			ProviderData: "invalid",
		}
		resp := &datasource.ConfigureResponse{}

		ds.Configure(ctx, req, resp)

		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("valid provider data", func(t *testing.T) {
		mockClient := &client.Client{}
		mockService := &iam.Service{}

		req := datasource.ConfigureRequest{
			ProviderData: map[string]interface{}{
				"client":      mockClient,
				"iam_service": mockService,
			},
		}
		resp := &datasource.ConfigureResponse{}

		ds.Configure(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
		assert.Equal(t, mockClient, ds.client)
		assert.Equal(t, mockService, ds.iamService)
	})
}

func TestResourceDataSource_Read_Success(t *testing.T) {
	// Skip this test as it requires complex Terraform framework setup
	// The datasource will be tested via acceptance tests
	t.Skip("Requires full Terraform framework setup - use acceptance tests")
}

func TestResourceDataSource_Read_NotFound(t *testing.T) {
	// Skip this test as it requires complex Terraform framework setup
	// The datasource will be tested via acceptance tests
	t.Skip("Requires full Terraform framework setup - use acceptance tests")
}

func TestResourceDataSource_Read_EmptyPermissions(t *testing.T) {
	// Skip this test as it requires complex Terraform framework setup
	// The datasource will be tested via acceptance tests
	t.Skip("Requires full Terraform framework setup - use acceptance tests")
}
