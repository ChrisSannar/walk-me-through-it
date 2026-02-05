package walker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Security constants
const (
	maxFileSize = 10 * 1024 * 1024 // 10MB limit
)

// Allowed file extensions for security
var allowedExtensions = []string{
	".go", ".js", ".ts", ".jsx", ".tsx",
	".py", ".rb", ".php",
	".java", ".kt", ".scala",
	".c", ".cpp", ".h", ".hpp",
	".rs", ".swift",
	".md", ".txt", ".json", ".yaml", ".yml",
	".html", ".css", ".scss", ".sass",
	".sh", ".bash", ".zsh",
}

// Walker handles file system operations for reading source files
type Walker struct {
	rootPath         string
	lastAuditMessage string
}

// NewWalker creates a new walker instance
func NewWalker(rootPath string) *Walker {
	return &Walker{
		rootPath: rootPath,
	}
}

// ReadFileLines reads a specific range of lines from a file
func (w *Walker) ReadFileLines(filePath string, start, end int) ([]string, error) {
	// Validate and sanitize path
	safePath, err := w.sanitizePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Check file extension
	if !w.isAllowedExtension(safePath) {
		return nil, fmt.Errorf("file type not allowed: %s", filepath.Ext(safePath))
	}

	// Check file size before reading
	fileInfo, err := os.Stat(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	if fileInfo.Size() > maxFileSize {
		return nil, fmt.Errorf("file exceeds maximum size of %d bytes", maxFileSize)
	}

	// Log access for audit
	w.lastAuditMessage = fmt.Sprintf("[AUDIT] Reading file: %s (lines %d-%d, size: %d bytes)", safePath, start, end, fileInfo.Size())
	log.Println(w.lastAuditMessage)

	// Open file with explicit read-only flag
	file, err := os.OpenFile(safePath, os.O_RDONLY, 0)
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
	safePath, err := w.sanitizePath(filePath)
	if err != nil {
		return false
	}
	_, err = os.Stat(safePath)
	return !os.IsNotExist(err)
}

// sanitizePath prevents directory traversal attacks
func (w *Walker) sanitizePath(filePath string) (string, error) {
	// Clean the path to resolve any . or .. components
	cleanPath := filepath.Clean(filePath)

	// Join with root path
	fullPath := filepath.Join(w.rootPath, cleanPath)

	// Resolve any symlinks
	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		// If the file doesn't exist yet, we can't resolve symlinks
		// but we can still check if the path is within bounds
		resolvedPath = fullPath
	}

	// Ensure the resolved path is within the project root
	// Add trailing separator to root to prevent partial matches
	rootWithSep := w.rootPath
	if !strings.HasSuffix(rootWithSep, string(filepath.Separator)) {
		rootWithSep += string(filepath.Separator)
	}

	if !strings.HasPrefix(resolvedPath+string(filepath.Separator), rootWithSep) &&
		resolvedPath != w.rootPath {
		return "", fmt.Errorf("path escapes project root: %s", filePath)
	}

	return resolvedPath, nil
}

// isAllowedExtension checks if the file has an allowed extension
func (w *Walker) isAllowedExtension(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// GetAllowedExtensions returns the list of allowed file extensions
func GetAllowedExtensions() []string {
	return allowedExtensions
}

// GetLastAuditMessage returns the last audit message
func (w *Walker) GetLastAuditMessage() string {
	return w.lastAuditMessage
}
