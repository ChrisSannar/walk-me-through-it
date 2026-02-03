package main

import (
	initcmd "github.com/chrissannar/walk-me-through-it/cmd/init"
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
	// TODO: Implement TUI launch
	return nil
}
