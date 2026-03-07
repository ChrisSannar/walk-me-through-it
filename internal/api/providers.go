package api

type AuthType string

const (
	AuthTypeBearer AuthType = "bearer"
	AuthTypeAPIKey AuthType = "api-key"
	AuthTypeNone   AuthType = "none"
)

type Provider struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name"`
	BaseURL        string   `json:"base_url"`
	AuthType       AuthType `json:"auth_type"`
	AuthHeader     string   `json:"auth_header"`
	DefaultModels  []string `json:"default_models"`
	SupportsStream bool     `json:"supports_stream"`
}

var Providers = []Provider{
	{
		ID:          "google",
		Name:        "google",
		DisplayName: "Google AI (Gemini)",
		BaseURL:     "https://generativelanguage.googleapis.com",
		AuthType:    "api-key",
		AuthHeader:  "X-Goog-Api-Key",
		DefaultModels: []string{
			"gemini-2.0-flash",
			"gemini-1.5-flash",
			"gemini-1.5-pro",
		},
		SupportsStream: false,
	},
	{
		ID:          "openai",
		Name:        "openai",
		DisplayName: "OpenAI",
		BaseURL:     "https://api.openai.com/v1",
		AuthType:    "bearer",
		AuthHeader:  "Authorization",
		DefaultModels: []string{
			"gpt-4o",
			"gpt-4o-mini",
			"gpt-4-turbo",
			"gpt-3.5-turbo",
		},
		SupportsStream: true,
	},
	{
		ID:          "anthropic",
		Name:        "anthropic",
		DisplayName: "Anthropic (Claude)",
		BaseURL:     "https://api.anthropic.com/v1",
		AuthType:    "api-key",
		AuthHeader:  "x-api-key",
		DefaultModels: []string{
			"claude-sonnet-4-20250514",
			"claude-3-5-sonnet-20241022",
			"claude-3-opus-20240229",
			"claude-3-haiku-20240307",
		},
		SupportsStream: true,
	},
	{
		ID:          "groq",
		Name:        "groq",
		DisplayName: "Groq",
		BaseURL:     "https://api.groq.com/openai/v1",
		AuthType:    "bearer",
		AuthHeader:  "Authorization",
		DefaultModels: []string{
			"llama-3.3-70b-versatile",
			"llama-3.1-70b-versatile",
			"mixtral-8x7b-32768",
		},
		SupportsStream: true,
	},
	{
		ID:          "ollama",
		Name:        "ollama",
		DisplayName: "Ollama (Local)",
		BaseURL:     "http://localhost:11434",
		AuthType:    "none",
		AuthHeader:  "",
		DefaultModels: []string{
			"llama3.2",
			"llama3.1",
			"mistral",
			"codellama",
		},
		SupportsStream: true,
	},
}

func GetProvider(id string) *Provider {
	for i := range Providers {
		if Providers[i].ID == id {
			return &Providers[i]
		}
	}
	return nil
}

func GetProviderByModel(modelName string) *Provider {
	for _, p := range Providers {
		for _, m := range p.DefaultModels {
			if m == modelName {
				return &p
			}
		}
	}
	return nil
}
