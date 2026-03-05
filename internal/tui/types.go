package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/chrissannar/walk-me-through-it/internal/highlighter"
	"github.com/chrissannar/walk-me-through-it/internal/navigator"
	"github.com/chrissannar/walk-me-through-it/internal/walker"
)

// AppState represents the current state of the application
type AppState int

const (
	StateSelecting AppState = iota
	StateEnteringPath
	StateLoading
	StateViewing
	StateLoadingStep
	StateConfirmDelete
	StateModelSelect
)

// walkthroughItem represents a walkthrough file for the list
type walkthroughItem struct {
	path string
	name string
}

func (w walkthroughItem) FilterValue() string { return w.name }
func (w walkthroughItem) Title() string       { return w.name }
func (w walkthroughItem) Description() string { return w.path }

// Model represents the TUI state
type Model struct {
	state        AppState
	width        int
	height       int
	list         list.Model
	textInput    textinput.Model
	navigator    *navigator.Navigator
	walker       *walker.Walker
	rootPath     string
	selectedFile string
	fileContent  []string
	lastAuditMsg string
	err          error
	styles       *Styles
	highlighter  *highlighter.Highlighter

	// Model selection
	selectedModel      string
	modelList          []string
	modelSelectedIndex int
	modelTextInput     textinput.Model
	modelIsAdding      bool
	modelAddingName    string
	modelAskingFor     string // "name" or "key"
	modelDeleteConfirm bool

	// Delete confirmation
	deleteConfirmPath string
	deleteConfirmFor  string // "model" or "walkthrough"

	// Line ranges for display
	highlightStart int
	highlightEnd   int
	displayStart   int
	displayEnd     int
}

// Messages
type walkthroughFilesMsg struct {
	files []string
}

type errMsg struct {
	err error
}

type walkthroughLoadedMsg struct {
	nav *navigator.Navigator
}

type fileContentMsg struct {
	lines          []string
	highlightStart int
	highlightEnd   int
	displayStart   int
	displayEnd     int
}
