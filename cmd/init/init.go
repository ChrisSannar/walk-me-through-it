package initcmd

import (
	"fmt"

	"github.com/chrissannar/walk-me-through-it/internal/config"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize wmti configuration",
		Long:  `Sets up initial configuration including optional API key storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}
}

func runInit() error {
	fmt.Println("Initializing wmti configuration...")

	cfg, err := config.LoadOrCreate()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	fmt.Println("Configuration initialized successfully!")
	fmt.Printf("Config location: %s\n", cfg.ConfigPath())

	return nil
}
