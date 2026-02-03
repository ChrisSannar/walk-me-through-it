package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWalkthroughFiles(t *testing.T) {
	// Test with current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	// Go to project root
	root := filepath.Join(cwd, "..", "..")

	files, err := FindWalkthroughFiles(root)
	if err != nil {
		t.Fatalf("FindWalkthroughFiles failed: %v", err)
	}

	t.Logf("Found %d files:", len(files))
	for _, f := range files {
		t.Logf("  - %s", f)
	}
}
