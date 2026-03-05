package tui

import (
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelSelectQGoesBack(t *testing.T) {
	m := Model{
		state:              StateModelSelect,
		modelIsAdding:      false,
		modelList:          []string{"model1", "model2", "+ Add new model"},
		modelSelectedIndex: 0,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	updatedModel, _ := m.Update(msg)

	m = updatedModel.(Model)

	if m.state == StateSelecting {
		t.Logf("SUCCESS: q in StateModelSelect goes back to StateSelecting")
	} else if m.state == StateModelSelect {
		t.Logf("STILL StateModelSelect - might quit")
	} else {
		t.Logf("State: %v", m.state)
	}

	if m.state != StateSelecting {
		t.Errorf("Expected state to be StateSelecting, got %v", m.state)
	}
}

func TestDeleteConfirmNStays(t *testing.T) {
	m := Model{
		state:             StateSelecting,
		deleteConfirmPath: "/some/path/test.wmti.json",
		list:              mustCreateList(),
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")}
	updatedModel, _ := m.Update(msg)

	m = updatedModel.(Model)

	if m.deleteConfirmPath != "" {
		t.Errorf("Expected deleteConfirmPath to be cleared after 'n', got %v", m.deleteConfirmPath)
	}
	if m.state != StateSelecting {
		t.Errorf("Expected state to stay StateSelecting after 'n', got %v", m.state)
	}
}

func TestDeleteConfirmYCancels(t *testing.T) {
	tmpFile := "/tmp/wmti_test_delete.wmti.json"
	os.WriteFile(tmpFile, []byte("{}"), 0644)
	defer os.Remove(tmpFile)

	m := Model{
		state:             StateSelecting,
		deleteConfirmPath: tmpFile,
		list:              mustCreateList(),
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")}
	updatedModel, _ := m.Update(msg)

	m = updatedModel.(Model)

	if m.deleteConfirmPath != "" {
		t.Errorf("Expected deleteConfirmPath to be cleared after delete, got %v", m.deleteConfirmPath)
	}
	if _, err := os.Stat(tmpFile); err == nil {
		t.Errorf("Expected file to be deleted after 'y', but it still exists")
	}
}

func mustCreateList() list.Model {
	return list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 24)
}
