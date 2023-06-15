package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMainExpectBadRequestError(t *testing.T) {
	go func() {
		// Run the main function.
		main()
	}()

	// Wait for a brief moment to allow the server to start.
	time.Sleep(100 * time.Millisecond)

	// Send a test request to the server.
	resp, err := http.Get("http://localhost:3000/ptlist")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status Bad Request; got %v", resp.StatusCode)
	}

	// Shut down the server gracefully.
	if err := resp.Body.Close(); err != nil {
		t.Errorf("failed to close response body: %v", err)
	}
}
