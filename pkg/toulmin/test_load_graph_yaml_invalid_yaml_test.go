//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraphYAMLInvalidYAML — tests error for invalid YAML content
package toulmin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGraphYAMLInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte(`{{{not yaml`), 0644)

	_, err := LoadGraphYAML(path, map[string]any{}, nil)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
