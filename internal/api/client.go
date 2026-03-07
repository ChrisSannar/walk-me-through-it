package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Client struct {
	provider   *Provider
	apiKey     string
	httpClient *http.Client
}

func NewClient(provider *Provider, apiKey string) *Client {
	return &Client{
		provider:   provider,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	switch c.provider.ID {
	case "google":
		return c.chatGoogle(ctx, messages)
	case "openai", "groq":
		return c.chatOpenAI(ctx, messages)
	case "anthropic":
		return c.chatAnthropic(ctx, messages)
	case "ollama":
		return c.chatOllama(ctx, messages)
	default:
		return "", fmt.Errorf("unsupported provider: %s", c.provider.ID)
	}
}

func (c *Client) chatOpenAI(ctx context.Context, messages []Message) (string, error) {
	url := c.provider.BaseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model:       c.provider.DefaultModels[0],
		Messages:    messages,
		MaxTokens:   4096,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}

type GoogleRequest struct {
	Contents         []GoogleContent        `json:"contents"`
	GenerationConfig GoogleGenerationConfig `json:"generationConfig"`
}

type GoogleContent struct {
	Parts []GooglePart `json:"parts"`
}

type GooglePart struct {
	Text string `json:"text"`
}

type GoogleGenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GoogleResponse struct {
	Candidates []GoogleCandidate `json:"candidates"`
}

type GoogleCandidate struct {
	Content GoogleContent `json:"content"`
}

func (c *Client) chatGoogle(ctx context.Context, messages []Message) (string, error) {
	model := c.provider.DefaultModels[0]
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.provider.BaseURL, model, c.apiKey)

	systemPrompt := ""
	userMessages := []string{}
	for _, m := range messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else if m.Role == "user" {
			userMessages = append(userMessages, m.Content)
		}
	}

	contents := []GoogleContent{}
	if systemPrompt != "" {
		contents = append(contents, GoogleContent{
			Parts: []GooglePart{{Text: systemPrompt}},
		})
	}
	for _, text := range userMessages {
		contents = append(contents, GoogleContent{
			Parts: []GooglePart{{Text: text}},
		})
	}

	reqBody := GoogleRequest{
		Contents: contents,
		GenerationConfig: GoogleGenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 4096,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var googleResp GoogleResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(googleResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	parts := googleResp.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return "", fmt.Errorf("no content parts returned")
	}

	return parts[0].Text, nil
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []AnthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
}

type AnthropicContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func (c *Client) chatAnthropic(ctx context.Context, messages []Message) (string, error) {
	url := c.provider.BaseURL + "/messages"

	var systemPrompt string
	anthropicMessages := []AnthropicMessage{}

	for _, m := range messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else {
			anthropicMessages = append(anthropicMessages, AnthropicMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}

	reqBody := AnthropicRequest{
		Model:     c.provider.DefaultModels[0],
		Messages:  anthropicMessages,
		MaxTokens: 4096,
		System:    systemPrompt,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no content returned")
	}

	return anthropicResp.Content[0].Text, nil
}

type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type OllamaResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func (c *Client) chatOllama(ctx context.Context, messages []Message) (string, error) {
	url := c.provider.BaseURL + "/api/chat"

	reqBody := OllamaRequest{
		Model:    c.provider.DefaultModels[0],
		Messages: messages,
		Stream:   false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

func (c *Client) TestConnection(ctx context.Context) error {
	testMessages := []Message{
		{Role: "user", Content: "Hi"},
	}
	_, err := c.Chat(ctx, testMessages)
	return err
}

func ParseMarkdownCodeBlocks(text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if !inCodeBlock {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}
