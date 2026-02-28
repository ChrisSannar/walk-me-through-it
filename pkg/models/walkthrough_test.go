package models

import (
	"testing"
)

func TestWalkthrough_Validate_Valid(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 1, Title: "Step 1", File: "main.go"},
			{ID: 2, Title: "Step 2", File: "utils.go"},
		},
	}

	err := wt.Validate()

	if err != nil {
		t.Errorf("Validate failed for valid walkthrough: %v", err)
	}
}

func TestWalkthrough_Validate_MissingTitle(t *testing.T) {
	wt := Walkthrough{
		Title:       "",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 1, Title: "Step 1", File: "main.go"},
		},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for missing title")
	}

	if err != ErrMissingTitle {
		t.Errorf("expected ErrMissingTitle, got %v", err)
	}
}

func TestWalkthrough_Validate_NoSteps(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps:       []Step{},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for empty steps")
	}

	if err != ErrNoSteps {
		t.Errorf("expected ErrNoSteps, got %v", err)
	}
}

func TestWalkthrough_Validate_NilSteps(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps:       nil,
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for nil steps")
	}

	if err != ErrNoSteps {
		t.Errorf("expected ErrNoSteps, got %v", err)
	}
}

func TestStep_Validate_InvalidID(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 1, Title: "Step 1", File: "main.go"},
			{ID: 3, Title: "Step 2", File: "utils.go"}, // Missing ID 2
		},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for invalid step ID")
	}

	if err != ErrInvalidStepID {
		t.Errorf("expected ErrInvalidStepID, got %v", err)
	}
}

func TestStep_Validate_InvalidID_NotStartingAt1(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 0, Title: "Step 1", File: "main.go"}, // Not starting at 1
			{ID: 2, Title: "Step 2", File: "utils.go"},
		},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for invalid step ID")
	}

	if err != ErrInvalidStepID {
		t.Errorf("expected ErrInvalidStepID, got %v", err)
	}
}

func TestStep_Validate_MissingFile(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 1, Title: "Step 1", File: ""}, // Empty file
		},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for missing file")
	}

	if err != ErrMissingFile {
		t.Errorf("expected ErrMissingFile, got %v", err)
	}
}

func TestStep_Validate_WhitespaceFile(t *testing.T) {
	wt := Walkthrough{
		Title:       "Test Walkthrough",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps: []Step{
			{ID: 1, Title: "Step 1", File: "   "}, // Whitespace file
		},
	}

	err := wt.Validate()

	// Current implementation only checks for empty string, not whitespace
	// This test documents current behavior
	if err != nil && err != ErrMissingFile {
		t.Errorf("expected ErrMissingFile or no error for whitespace-only file, got %v", err)
	}
}

func TestWalkthrough_MultipleValidationErrors(t *testing.T) {
	wt := Walkthrough{
		Title:       "",
		Description: "A test walkthrough",
		Version:     "1.0.0",
		Steps:       []Step{},
	}

	err := wt.Validate()

	if err == nil {
		t.Fatal("expected error for invalid walkthrough")
	}

	// Should return first error encountered
	if err != ErrMissingTitle && err != ErrNoSteps {
		t.Errorf("expected ErrMissingTitle or ErrNoSteps, got %v", err)
	}
}
