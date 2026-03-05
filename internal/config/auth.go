package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Auth map[string]string

func AuthPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "share", "wmti", "auth.json"), nil
}

func LoadAuth() (Auth, error) {
	authPath, err := AuthPath()
	if err != nil {
		return nil, err
	}

	auth := Auth{}

	data, err := os.ReadFile(authPath)
	if err != nil {
		if os.IsNotExist(err) {
			return auth, nil
		}
		return nil, fmt.Errorf("failed to read auth file: %w", err)
	}

	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, fmt.Errorf("failed to parse auth file: %w", err)
	}

	return auth, nil
}

func SaveAuth(auth Auth) error {
	authPath, err := AuthPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(authPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	data, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth: %w", err)
	}

	if err := os.WriteFile(authPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}

	return nil
}

func GetAPIKey(modelName string) (string, error) {
	auth, err := LoadAuth()
	if err != nil {
		return "", err
	}
	return auth[modelName], nil
}

func SetAPIKey(modelName, apiKey string) error {
	auth, err := LoadAuth()
	if err != nil {
		return err
	}
	auth[modelName] = apiKey
	return SaveAuth(auth)
}

func DeleteAPIKey(modelName string) error {
	auth, err := LoadAuth()
	if err != nil {
		return err
	}
	delete(auth, modelName)
	return SaveAuth(auth)
}
