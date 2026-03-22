//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraphYAMLMissingFunc — tests error for missing function in YAML-loaded graph
package toulmin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGraphYAMLMissingFunc(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(`
graph: test
rules:
  - name: missing
    role: warrant
`), 0644)

	_, err := LoadGraphYAML(path, map[string]any{}, nil)
	if err == nil {
		t.Fatal("expected error for missing function")
	}
}
