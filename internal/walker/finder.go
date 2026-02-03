package walker

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindWalkthroughFiles searches for walkthrough JSON files in the given directory
// It looks in the current directory and .wmti/ subdirectory
func FindWalkthroughFiles(rootPath string) ([]string, error) {
	var files []string

	// Search patterns - only files following wmti_<title>.json naming convention
	patterns := []string{
		"wmti_*.json",       // Root level walkthrough files
		".wmti/wmti_*.json", // .wmti subdirectory
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(rootPath, pattern))
		if err != nil {
			return nil, fmt.Errorf("error searching for walkthrough files: %w", err)
		}
		files = append(files, matches...)
	}

	// Validate that files are valid walkthroughs by checking if they can be parsed
	var validFiles []string
	for _, file := range files {
		if isValidWalkthroughFile(file) {
			validFiles = append(validFiles, file)
		}
	}

	return validFiles, nil
}

// isValidWalkthroughFile checks if a file is a valid walkthrough JSON
func isValidWalkthroughFile(path string) bool {
	// For now, just check if it's a valid JSON file
	// In the future, we could validate the schema
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	// Check if it starts with '{' (basic JSON object check)
	if len(data) == 0 {
		return false
	}

	// Skip whitespace and check for opening brace
	for i := 0; i < len(data) && i < 100; i++ {
		if data[i] == '{' {
			return true
		}
		if data[i] != ' ' && data[i] != '\t' && data[i] != '\n' && data[i] != '\r' {
			return false
		}
	}

	return false
}

// GetDefaultWalkthroughPath returns the default walkthrough path if it exists
func GetDefaultWalkthroughPath(rootPath string) string {
	defaultPath := filepath.Join(rootPath, ".wmti", "walkthrough.json")
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath
	}
	return ""
}
