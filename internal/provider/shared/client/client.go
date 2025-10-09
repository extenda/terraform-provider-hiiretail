package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
)

// Config holds the configuration for the API client
type Config struct {
	BaseURL      string
	IAMEndpoint  string
	CCCEndpoint  string
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

// DefaultConfig returns a default client configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL:      "https://iam-api.retailsvc.com",
		IAMEndpoint:  "/api/v1", // V1 API prefix for most IAM resources (groups, roles, custom roles, resources)
		CCCEndpoint:  "/ccc/v1",
		UserAgent:    "terraform-provider-hiiretail/1.0.0",
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}
}

// Client is the unified API client for HiiRetail services
type Client struct {
	config     *Config
	httpClient *http.Client
	auth       *auth.Config
	baseURL    *url.URL
	tenantID   string
}

// New creates a new HiiRetail API client
func New(authConfig *auth.Config, clientConfig *Config) (*Client, error) {
	if clientConfig == nil {
		clientConfig = DefaultConfig()
	}

	baseURL, err := url.Parse(clientConfig.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	var httpClient *http.Client
	if authConfig != nil && authConfig.TestToken != "" {
		// Use basic http.Client for contract tests with dummy token
		httpClient = &http.Client{Timeout: clientConfig.Timeout}
	} else {
		// Create OAuth2 HTTP client
		var err error
		httpClient, err = auth.NewHTTPClient(context.Background(), authConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create OAuth2 client: %w", err)
		}
		// Set timeout
		httpClient.Timeout = clientConfig.Timeout
	}

	return &Client{
		config:     clientConfig,
		httpClient: httpClient,
		auth:       authConfig,
		baseURL:    baseURL,
		tenantID:   authConfig.TenantID,
	}, nil
}

// Request represents an API request
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   map[string]string
}

// Response represents an API response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Do executes an API request
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	// Build URL
	reqURL := c.buildURL(req.Path)
	if len(req.Query) > 0 {
		q := reqURL.Query()
		for key, value := range req.Query {
			q.Set(key, value)
		}
		reqURL.RawQuery = q.Encode()
	}

	// Prepare body
	var body io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		body = strings.NewReader(string(bodyBytes))
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, reqURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.config.UserAgent)
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	// If TestToken is set, use it for Authorization and skip real OAuth2
	if c.auth != nil && c.auth.TestToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.auth.TestToken)
	}

	// Execute request with retries
	resp, err := c.doWithRetry(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// buildURL constructs the full URL for a request path
func (c *Client) buildURL(path string) *url.URL {
	u := *c.baseURL // Copy
	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")
	fmt.Fprintf(os.Stderr, "[DEBUG buildURL] Input path: '%s', BaseURL: '%s', Final URL: '%s'\n", path, c.baseURL.String(), u.String())
	return &u
}

// doWithRetry executes HTTP request with retry logic
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay
			delay := c.calculateBackoff(attempt)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Check if we should retry based on status code
		if c.shouldRetry(resp.StatusCode) {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// shouldRetry determines if a request should be retried based on status code
func (c *Client) shouldRetry(statusCode int) bool {
	// Retry on server errors and rate limiting
	return statusCode >= 500 || statusCode == 429
}

// calculateBackoff calculates exponential backoff with jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	min := c.config.RetryWaitMin
	max := c.config.RetryWaitMax

	// Exponential backoff: min * 2^(attempt-1)
	backoff := min * time.Duration(1<<uint(attempt-1))

	if backoff > max {
		backoff = max
	}

	// Add jitter (±25%)
	jitter := time.Duration(float64(backoff) * 0.25)
	jitterMultiplier := (2*float64(time.Now().UnixNano()%1000)/1000.0 - 1)
	backoff += time.Duration(float64(jitter) * jitterMultiplier)

	if backoff < min {
		backoff = min
	}

	return backoff
}

// IAMClient returns a client configured for IAM service endpoints
func (c *Client) IAMClient() *ServiceClient {
	return &ServiceClient{
		client:   c,
		endpoint: c.config.IAMEndpoint,
		service:  "iam",
	}
}

// CCCClient returns a client configured for CCC service endpoints
func (c *Client) CCCClient() *ServiceClient {
	return &ServiceClient{
		client:   c,
		endpoint: c.config.CCCEndpoint,
		service:  "ccc",
	}
}

// ServiceClient wraps the main client for service-specific operations
type ServiceClient struct {
	client   *Client
	endpoint string
	service  string
}

// Do executes a request with the service endpoint prefix
func (sc *ServiceClient) Do(ctx context.Context, req *Request) (*Response, error) {
	originalPath := req.Path
	// Prefix path with service endpoint
	req.Path = strings.TrimSuffix(sc.endpoint, "/") + "/" + strings.TrimPrefix(req.Path, "/")

	fmt.Fprintf(os.Stderr, "[DEBUG ServiceClient.Do] Service: '%s', Endpoint: '%s', Original path: '%s', Final path: '%s'\n",
		sc.service, sc.endpoint, originalPath, req.Path)

	return sc.client.Do(ctx, req)
}

// Get performs a GET request
func (sc *ServiceClient) Get(ctx context.Context, path string, query map[string]string) (*Response, error) {
	return sc.Do(ctx, &Request{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	})
}

// Post performs a POST request
func (sc *ServiceClient) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return sc.Do(ctx, &Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request
func (sc *ServiceClient) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return sc.Do(ctx, &Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

// Patch performs a PATCH request
func (sc *ServiceClient) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	return sc.Do(ctx, &Request{
		Method: http.MethodPatch,
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request
func (sc *ServiceClient) Delete(ctx context.Context, path string) (*Response, error) {
	return sc.Do(ctx, &Request{
		Method: http.MethodDelete,
		Path:   path,
	})
}

// TenantID returns the tenant ID configured for this client
func (c *Client) TenantID() string {
	return c.tenantID
}

// HTTPClient returns the underlying HTTP client
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// BaseURL returns the base URL configured for this client
func (c *Client) BaseURL() string {
	return c.baseURL.String()
}
