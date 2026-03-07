package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ModelConfig holds configuration for a model
type ModelConfig struct {
	Name     string `mapstructure:"name"`
	Provider string `mapstructure:"provider"`
}

// Config holds the application configuration
type Config struct {
	Models     []ModelConfig `mapstructure:"models"`
	Selected   string        `mapstructure:"selected"`
	ConfigDir  string
	ConfigFile string
}

// LoadOrCreate loads existing config or creates a new one
func LoadOrCreate() (*Config, error) {
	cfg := &Config{}

	// Determine config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cfg.ConfigDir = filepath.Join(homeDir, ".config", "wmti")
	cfg.ConfigFile = filepath.Join(cfg.ConfigDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cfg.ConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfg.ConfigDir)

	// Set defaults
	viper.SetDefault("models", []ModelConfig{})
	viper.SetDefault("selected", "")

	// Try to read existing config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found; create it
		if err := viper.WriteConfigAs(cfg.ConfigFile); err != nil {
			return nil, fmt.Errorf("failed to write config: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// ConfigPath returns the path to the config file
func (c *Config) ConfigPath() string {
	return c.ConfigFile
}

// Save saves the current configuration
func (c *Config) Save() error {
	viper.Set("models", c.Models)
	viper.Set("selected", c.Selected)
	return viper.WriteConfig()
}

// GetModelList returns just the model names for backwards compatibility
func (c *Config) GetModelList() []string {
	names := make([]string, len(c.Models))
	for i, m := range c.Models {
		names[i] = m.Name
	}
	return names
}

// GetSelectedModelConfig returns the config for the selected model
func (c *Config) GetSelectedModelConfig() *ModelConfig {
	for _, m := range c.Models {
		if m.Name == c.Selected {
			return &m
		}
	}
	return nil
}
