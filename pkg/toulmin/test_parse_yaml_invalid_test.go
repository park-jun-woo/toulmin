//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLInvalid — tests error for invalid YAML content
package toulmin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAMLInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte(`{{{not yaml`), 0644)

	_, err := ParseYAML(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
