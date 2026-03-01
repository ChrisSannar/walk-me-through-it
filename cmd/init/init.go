package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrissannar/walk-me-through-it/internal/config"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize wmti configuration",
		Long:  `Sets up initial configuration including .wmti/ directory and optional API key storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit()
		},
	}
}

func RunInit() error {
	fmt.Println("Initializing wmti configuration...")

	// Load or create config
	cfg, err := config.LoadOrCreate()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Get current working directory for project setup
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create .wmti/ directory
	wmtiDir := filepath.Join(cwd, ".wmti")
	if err := os.MkdirAll(wmtiDir, 0755); err != nil {
		return fmt.Errorf("failed to create .wmti directory: %w", err)
	}
	fmt.Printf("Created: %s\n", wmtiDir)

	// Create a sample walkthrough template
	samplePath := filepath.Join(wmtiDir, "example.wmti.json")
	if _, err := os.Stat(samplePath); os.IsNotExist(err) {
		sample := `{
  "title": "How wmti Works",
  "description": "A quick tour of how wmti walks you through codebases",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "JSON Walkthrough Files",
      "description": "wmti creates and reads these JSON files to walk through your project. Each file describes a walkthrough with a title, description, and a series of steps.",
      "file": ".wmti/example.wmti.json",
      "line_start": 2,
      "line_end": 4,
      "action": "read"
    },
    {
      "id": 2,
      "title": "Each Step Defines a View",
      "description": "Each 'step' in the steps array tells wmti which file to display, which lines to show, and what description to display in the sidebar.",
      "file": ".wmti/example.wmti.json",
      "line_start": 18,
      "line_end": 18,
      "action": "read"
    },
    {
      "id": 3,
      "title": "Cyclic Navigation",
      "description": "When you reach the last step, pressing Tab will loop back to the first step. Similarly, Shift+Tab cycles backwards.",
      "file": ".wmti/example.wmti.json",
      "line_start": 25,
      "line_end": 31,
      "action": "read"
    },
    {
      "id": 4,
      "title": "Creating Walkthroughs",
      "description": "You can create new walkthroughs by simply adding a JSON file to the .wmti/ folder. Each file must follow the *.wmti.json naming convention.",
      "file": ".wmti/example.wmti.json",
      "line_start": 0,
      "line_end": 0,
      "action": "read"
    }
  ]
}
`
		if err := os.WriteFile(samplePath, []byte(sample), 0644); err != nil {
			return fmt.Errorf("failed to create sample walkthrough: %w", err)
		}
		fmt.Printf("Created: %s\n", samplePath)
	}

	fmt.Println("\nConfiguration initialized successfully!")
	fmt.Printf("Config location: %s\n", cfg.ConfigPath())
	fmt.Printf("Project directory: %s\n", wmtiDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit .wmti/example.wmti.json to create your walkthrough")
	fmt.Println("  2. Or run 'wmti' to use the self-referential tutorial")

	return nil
}
