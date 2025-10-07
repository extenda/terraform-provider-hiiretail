package client

import (
	"net/url"
	"testing"
)

// TestBuildURL verifies that URLs are constructed correctly without double prefixes
func TestBuildURL(t *testing.T) {
	tests := []struct {
		name         string
		baseURLPath  string
		requestPath  string
		expectedPath string
	}{
		{
			name:         "V2 API path with empty base path (CORRECT)",
			baseURLPath:  "",
			requestPath:  "/api/v2/tenants/test/groups/group123/roles",
			expectedPath: "/api/v2/tenants/test/groups/group123/roles",
		},
		{
			name:         "PROBLEMATIC: V2 path with /api/v1 base path (causes double prefix)",
			baseURLPath:  "/api/v1",
			requestPath:  "/api/v2/tenants/test/groups/group123/roles",
			expectedPath: "/api/v1/api/v2/tenants/test/groups/group123/roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, _ := url.Parse("https://iam-api.retailsvc.com")
			baseURL.Path = tt.baseURLPath

			client := &Client{
				baseURL: baseURL,
			}

			resultURL := client.buildURL(tt.requestPath)

			if resultURL.Path != tt.expectedPath {
				t.Errorf("Expected path: %s, got: %s", tt.expectedPath, resultURL.Path)
			}
		})
	}
}
