//ff:func feature=analyzer type=analyzer control=sequence
//ff:what writeGoFile — helper that writes Go source to a temp file for testing
package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func writeGoFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
