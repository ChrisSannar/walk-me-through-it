package navigator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewNavigator(t *testing.T) {
	n := NewNavigator()
	if n == nil {
		t.Fatal("NewNavigator returned nil")
	}
	if n.currentStep != 0 {
		t.Errorf("expected currentStep 0, got %d", n.currentStep)
	}
	if n.walkthrough != nil {
		t.Error("expected walkthrough to be nil")
	}
}

func TestLoadWalkthrough_ValidFile(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	err := n.LoadWalkthrough(testFile)

	if err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	if n.walkthrough == nil {
		t.Fatal("walkthrough should not be nil")
	}

	if n.walkthrough.Title == "" {
		t.Error("title should not be empty")
	}

	if len(n.walkthrough.Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(n.walkthrough.Steps))
	}

	if n.currentStep != 0 {
		t.Errorf("expected currentStep 0, got %d", n.currentStep)
	}
}

func TestLoadWalkthrough_InvalidJSON(t *testing.T) {
	n := NewNavigator()

	tmpFile := filepath.Join(t.TempDir(), "invalid.json")
	err := os.WriteFile(tmpFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	err = n.LoadWalkthrough(tmpFile)

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	if n.walkthrough != nil {
		t.Error("walkthrough should be nil after failed load")
	}
}

func TestLoadWalkthrough_MissingFile(t *testing.T) {
	n := NewNavigator()

	err := n.LoadWalkthrough("/nonexistent/path.json")

	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCurrentStep_Empty(t *testing.T) {
	n := NewNavigator()

	_, err := n.CurrentStep()

	if err == nil {
		t.Fatal("expected error when no walkthrough loaded")
	}
}

func TestCurrentStep_AfterLoad(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	step, err := n.CurrentStep()

	if err != nil {
		t.Fatalf("CurrentStep failed: %v", err)
	}

	if step.ID != 1 {
		t.Errorf("expected step ID 1, got %d", step.ID)
	}

	if step.File == "" {
		t.Error("step file should not be empty")
	}
}

func TestNext_Advances(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	step, err := n.Next()

	if err != nil {
		t.Fatalf("Next failed: %v", err)
	}

	if step.ID != 2 {
		t.Errorf("expected step ID 2, got %d", step.ID)
	}

	if n.currentStep != 1 {
		t.Errorf("expected currentStep 1, got %d", n.currentStep)
	}
}

func TestNext_AtEnd(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// Advance to the end
	n.currentStep = len(n.walkthrough.Steps) - 1

	_, err := n.Next()

	if err == nil {
		t.Fatal("expected error when at last step")
	}
}

func TestPrevious_GoesBack(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// Move to second step
	n.currentStep = 1

	step, err := n.Previous()

	if err != nil {
		t.Fatalf("Previous failed: %v", err)
	}

	if step.ID != 1 {
		t.Errorf("expected step ID 1, got %d", step.ID)
	}

	if n.currentStep != 0 {
		t.Errorf("expected currentStep 0, got %d", n.currentStep)
	}
}

func TestPrevious_AtStart(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// Already at first step
	_, err := n.Previous()

	if err == nil {
		t.Fatal("expected error when at first step")
	}
}

func TestNextCycle_Loops(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// Move to last step
	n.currentStep = len(n.walkthrough.Steps) - 1
	lastStepID := n.walkthrough.Steps[n.currentStep].ID

	step, err := n.NextCycle()

	if err != nil {
		t.Fatalf("NextCycle failed: %v", err)
	}

	// Should loop back to first step
	if step.ID != 1 {
		t.Errorf("expected step ID 1 after cycling, got %d", step.ID)
	}

	// Verify last step was the original last
	if lastStepID != 3 {
		t.Errorf("expected last step ID 3, got %d", lastStepID)
	}
}

func TestPreviousCycle_Loops(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// Already at first step
	step, err := n.PreviousCycle()

	if err != nil {
		t.Fatalf("PreviousCycle failed: %v", err)
	}

	// Should loop to last step
	if step.ID != 3 {
		t.Errorf("expected step ID 3 after cycling, got %d", step.ID)
	}
}

func TestHasNext(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// At first step, should have next
	if !n.HasNext() {
		t.Error("HasNext should return true at first step")
	}

	// Move to last step
	n.currentStep = len(n.walkthrough.Steps) - 1
	if n.HasNext() {
		t.Error("HasNext should return false at last step")
	}
}

func TestHasPrevious(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// At first step, should not have previous
	if n.HasPrevious() {
		t.Error("HasPrevious should return false at first step")
	}

	// Move to second step
	n.currentStep = 1
	if !n.HasPrevious() {
		t.Error("HasPrevious should return true at second step")
	}
}

func TestProgress_ReturnsCorrectValues(t *testing.T) {
	n := NewNavigator()

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	// At first step
	current, total := n.Progress()
	if current != 1 {
		t.Errorf("expected current 1, got %d", current)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}

	// Move to second step
	n.currentStep = 1
	current, total = n.Progress()
	if current != 2 {
		t.Errorf("expected current 2, got %d", current)
	}
}

func TestProgress_Empty(t *testing.T) {
	n := NewNavigator()

	current, total := n.Progress()

	if current != 0 {
		t.Errorf("expected current 0, got %d", current)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestWalkthrough_ReturnsValue(t *testing.T) {
	n := NewNavigator()

	if n.Walkthrough() != nil {
		t.Error("expected nil walkthrough before load")
	}

	testFile := filepath.Join("..", "..", "testdata", "sample-walkthrough.json")
	if err := n.LoadWalkthrough(testFile); err != nil {
		t.Fatalf("LoadWalkthrough failed: %v", err)
	}

	wt := n.Walkthrough()
	if wt == nil {
		t.Error("walkthrough should not be nil after load")
	}

	if wt.Title != "Understanding the Authentication Flow" {
		t.Errorf("expected title 'Understanding the Authentication Flow', got '%s'", wt.Title)
	}
}
