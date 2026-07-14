package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDirFailsOnInvalidFrontmatter(t *testing.T) {
	dir := t.TempDir()
	postPath := filepath.Join(dir, "broken.md")
	if err := os.WriteFile(postPath, []byte("---\ntitle: [broken\n---\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := NewParser().ParseDir(dir)
	if err == nil {
		t.Fatal("ParseDir() error = nil, want invalid frontmatter error")
	}
	if !strings.Contains(err.Error(), postPath) {
		t.Fatalf("ParseDir() error = %q, want post path %q", err, postPath)
	}
}
