package navigator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrissannar/walk-me-through-it/pkg/models"
)

// Navigator handles walkthrough document parsing and navigation
type Navigator struct {
	walkthrough *models.Walkthrough
	currentStep int
	rootPath    string
}

// NewNavigator creates a new navigator instance
func NewNavigator() *Navigator {
	return &Navigator{
		currentStep: 0,
	}
}

func NewNavigatorWithRoot(rootPath string) *Navigator {
	return &Navigator{
		currentStep: 0,
		rootPath:    rootPath,
	}
}

// LoadWalkthrough loads a walkthrough from a JSON file
func (n *Navigator) LoadWalkthrough(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read walkthrough file: %w", err)
	}

	// Set rootPath from walkthrough file location if not set
	if n.rootPath == "" {
		n.rootPath = filepath.Dir(path)
	}

	var wt models.Walkthrough
	if err := json.Unmarshal(data, &wt); err != nil {
		return fmt.Errorf("failed to parse walkthrough JSON: %w", err)
	}

	if err := wt.Validate(); err != nil {
		return fmt.Errorf("invalid walkthrough: %w", err)
	}

	// Validate all step file paths are within root
	for _, step := range wt.Steps {
		if err := n.validateStepFilePath(step.File); err != nil {
			return fmt.Errorf("invalid step %d: %w", step.ID, err)
		}
	}

	n.walkthrough = &wt
	n.currentStep = 0
	return nil
}

func (n *Navigator) validateStepFilePath(filePath string) error {
	cleanPath := filepath.Clean(filePath)
	fullPath := filepath.Join(n.rootPath, cleanPath)

	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		resolvedPath = fullPath
	}

	rootWithSep := n.rootPath
	if !strings.HasSuffix(rootWithSep, string(filepath.Separator)) {
		rootWithSep += string(filepath.Separator)
	}

	if !strings.HasPrefix(resolvedPath+string(filepath.Separator), rootWithSep) &&
		resolvedPath != n.rootPath {
		return fmt.Errorf("path escapes root: %s", filePath)
	}

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

// NextCycle advances to the next step, cycling back to first if at the end
func (n *Navigator) NextCycle() (*models.Step, error) {
	if n.walkthrough == nil {
		return nil, fmt.Errorf("no walkthrough loaded")
	}
	if n.currentStep+1 >= len(n.walkthrough.Steps) {
		// Cycle back to first step
		n.currentStep = 0
	} else {
		n.currentStep++
	}
	return &n.walkthrough.Steps[n.currentStep], nil
}

// PreviousCycle goes to the previous step, cycling to last if at the beginning
func (n *Navigator) PreviousCycle() (*models.Step, error) {
	if n.walkthrough == nil {
		return nil, fmt.Errorf("no walkthrough loaded")
	}
	if n.currentStep <= 0 {
		// Cycle to last step
		n.currentStep = len(n.walkthrough.Steps) - 1
	} else {
		n.currentStep--
	}
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
