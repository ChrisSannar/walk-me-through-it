package models

// Walkthrough represents the root structure of a walkthrough document
type Walkthrough struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Steps       []Step `json:"steps"`
}

// Step represents a single step in the walkthrough
type Step struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file"`
	LineStart   int    `json:"line_start"`
	LineEnd     int    `json:"line_end"`
	Action      string `json:"action"`
}

// Validate checks if the walkthrough document is valid
func (w *Walkthrough) Validate() error {
	if w.Title == "" {
		return ErrMissingTitle
	}
	if len(w.Steps) == 0 {
		return ErrNoSteps
	}
	for i, step := range w.Steps {
		if step.ID != i+1 {
			return ErrInvalidStepID
		}
		if step.File == "" {
			return ErrMissingFile
		}
	}
	return nil
}
