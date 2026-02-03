package api

// API handles optional AI API integration for generating walkthroughs
type API struct {
	apiKey string
}

// NewAPI creates a new API client
func NewAPI(apiKey string) *API {
	return &API{
		apiKey: apiKey,
	}
}

// GenerateWalkthrough generates a walkthrough document from a codebase
// This is a placeholder for future AI integration
func (a *API) GenerateWalkthrough(codebasePath string) ([]byte, error) {
	// TODO: Implement AI API integration
	return nil, nil
}
