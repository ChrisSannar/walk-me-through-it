package main

import (
	"fmt"

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
		return runTUI()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initcmd.NewInitCommand())
}

func runTUI() error {
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	return nil
}
