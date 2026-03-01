package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	initcmd "github.com/chrissannar/walk-me-through-it/cmd/init"
	"github.com/chrissannar/walk-me-through-it/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wmti",
	Short: "Walk Me Through It - Interactive codebase walkthrough",
	Long: `Walk Me Through It (wmti) is a terminal-based IDE-like experience
that helps developers understand unfamiliar codebases by following
structured walkthrough documents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isInitialized() {
			fmt.Println("wmti not initialized. Running init...")
			if err := initcmd.RunInit(); err != nil {
				return fmt.Errorf("auto-init failed: %w", err)
			}
			fmt.Println()
		}
		return runTUI()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initcmd.NewInitCommand())
}

func isInitialized() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	wmtiDir := filepath.Join(cwd, ".wmti")
	info, err := os.Stat(wmtiDir)
	if err != nil || !info.IsDir() {
		return false
	}
	entries, err := os.ReadDir(wmtiDir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".wmti.json") {
			return true
		}
	}
	return false
}

func runTUI() error {
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	return nil
}
