package walker

import (
	"bufio"
	"fmt"
	"os"
)

// Walker handles file system operations for reading source files
type Walker struct {
	rootPath string
}

// NewWalker creates a new walker instance
func NewWalker(rootPath string) *Walker {
	return &Walker{
		rootPath: rootPath,
	}
}

// ReadFileLines reads a specific range of lines from a file
func (w *Walker) ReadFileLines(filePath string, start, end int) ([]string, error) {
	fullPath := w.rootPath + "/" + filePath

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum >= start && lineNum <= end {
			lines = append(lines, scanner.Text())
		}
		if lineNum > end {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return lines, nil
}

// FileExists checks if a file exists
func (w *Walker) FileExists(filePath string) bool {
	fullPath := w.rootPath + "/" + filePath
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}
