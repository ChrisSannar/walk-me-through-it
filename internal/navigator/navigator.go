package navigator

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chrissannar/walk-me-through-it/pkg/models"
)

// Navigator handles walkthrough document parsing and navigation
type Navigator struct {
	walkthrough *models.Walkthrough
	currentStep int
}

// NewNavigator creates a new navigator instance
func NewNavigator() *Navigator {
	return &Navigator{
		currentStep: 0,
	}
}

// LoadWalkthrough loads a walkthrough from a JSON file
func (n *Navigator) LoadWalkthrough(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read walkthrough file: %w", err)
	}

	var wt models.Walkthrough
	if err := json.Unmarshal(data, &wt); err != nil {
		return fmt.Errorf("failed to parse walkthrough JSON: %w", err)
	}

	if err := wt.Validate(); err != nil {
		return fmt.Errorf("invalid walkthrough: %w", err)
	}

	n.walkthrough = &wt
	n.currentStep = 0
	return nil
}

// CurrentStep returns the current step
func (n *Navigator) CurrentStep() (*models.Step, error) {
	if n.walkthrough == nil {
		return nil, fmt.Errorf("no walkthrough loaded")
	}
	if n.currentStep >= len(n.walkthrough.Steps) {
		return nil, fmt.Errorf("no more steps")
	}
	return &n.walkthrough.Steps[n.currentStep], nil
}

// Next advances to the next step
func (n *Navigator) Next() (*models.Step, error) {
	if n.walkthrough == nil {
		return nil, fmt.Errorf("no walkthrough loaded")
	}
	if n.currentStep+1 >= len(n.walkthrough.Steps) {
		return nil, fmt.Errorf("already at last step")
	}
	n.currentStep++
	return &n.walkthrough.Steps[n.currentStep], nil
}

// Previous goes back to the previous step
func (n *Navigator) Previous() (*models.Step, error) {
	if n.walkthrough == nil {
		return nil, fmt.Errorf("no walkthrough loaded")
	}
	if n.currentStep <= 0 {
		return nil, fmt.Errorf("already at first step")
	}
	n.currentStep--
	return &n.walkthrough.Steps[n.currentStep], nil
}

// HasNext returns true if there are more steps
func (n *Navigator) HasNext() bool {
	if n.walkthrough == nil {
		return false
	}
	return n.currentStep < len(n.walkthrough.Steps)-1
}

// HasPrevious returns true if there are previous steps
func (n *Navigator) HasPrevious() bool {
	return n.currentStep > 0
}

// Progress returns current step number and total steps
func (n *Navigator) Progress() (current, total int) {
	if n.walkthrough == nil {
		return 0, 0
	}
	return n.currentStep + 1, len(n.walkthrough.Steps)
}

// Walkthrough returns the loaded walkthrough
func (n *Navigator) Walkthrough() *models.Walkthrough {
	return n.walkthrough
}
