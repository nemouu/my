package ai

// This file will contain the Ollama client wrapper for communicating with
// the local LLM API. It handles connection management, request/response
// handling, and error management.
//
// TODO: Implement Ollama client as described in TODO.md Milestone 3 "Ollama Client"
// - Create Client struct with baseURL and http.Client
// - Implement Generate(model, prompt string) (string, error)
// - Add connection testing
// - Handle timeouts and Ollama not running errors
// - Add streaming support for future use

// Client represents an Ollama API client
type Client struct {
	// TODO: Add fields for baseURL, http.Client, etc.
}

// NewClient creates a new Ollama client
func NewClient(baseURL string) *Client {
	// TODO: Initialize and return client
	return nil
}

// Generate sends a prompt to the specified model and returns the response
func (c *Client) Generate(model, prompt string) (string, error) {
	// TODO: Implement LLM generation
	return "", nil
}
