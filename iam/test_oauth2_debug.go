package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
)

func main() {
	// Set up test environment like in acceptance tests
	env := &testutils.TestEnvironment{}
	env = testutils.SetupTestEnvironment(nil)

	// Simulate the mock server setup but with debug output
	fmt.Println("Setting up mock server...")
	// Note: We can't use t.Cleanup here since this isn't a test, so we'll manually close
	env.SetupMockServer(nil)
	defer env.MockServer.Close()

	fmt.Printf("Mock server running at: %s\n", env.BaseURL)

	// Try to authenticate like the provider does
	config := &clientcredentials.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     fmt.Sprintf("%s/oauth2/token", env.BaseURL),
	}

	fmt.Println("Creating OAuth2 client...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Getting OAuth2 token...")
	httpClient := config.Client(ctx)

	// Try to make a test request to trigger token acquisition
	fmt.Println("Making test request to trigger token...")
	resp, err := httpClient.Get(env.BaseURL + "/test")
	if err != nil {
		log.Printf("Error making request: %v", err)
	} else {
		resp.Body.Close()
		fmt.Printf("Got response: %d\n", resp.StatusCode)
	}

	fmt.Println("OAuth2 test completed")
}
