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
			return runInit()
		},
	}
}

func runInit() error {
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
  "title": "Example Walkthrough",
  "description": "A template walkthrough to get you started",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "First Step",
      "description": "Description of what this step explains",
      "file": "path/to/file.go",
      "line_start": 1,
      "line_end": 10,
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
