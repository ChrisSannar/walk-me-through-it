package api

import (
	"context"
	"fmt"
)

// API handles optional AI API integration for generating walkthroughs
type API struct {
	provider *Provider
	apiKey   string
	model    string
}

// NewAPI creates a new API client
func NewAPI(provider *Provider, apiKey string, model string) *API {
	return &API{
		provider: provider,
		apiKey:   apiKey,
		model:    model,
	}
}

// GenerateWalkthrough generates a walkthrough document from a codebase
func (a *API) GenerateWalkthrough(ctx context.Context, codebasePath string, prompt string) (string, error) {
	if a.provider == nil {
		return "", fmt.Errorf("no provider configured")
	}

	client := NewClient(a.provider, a.apiKey)

	messages := []Message{
		{Role: "system", Content: "You are a code walkthrough generator. Generate a JSON walkthrough file based on the provided code."},
		{Role: "user", Content: prompt},
	}

	response, err := client.Chat(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	return response, nil
}

// TestConnection tests the API connection
func (a *API) TestConnection(ctx context.Context) error {
	if a.provider == nil {
		return fmt.Errorf("no provider configured")
	}

	client := NewClient(a.provider, a.apiKey)
	return client.TestConnection(ctx)
}
