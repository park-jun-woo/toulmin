//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLInvalidYAML — tests error returned for malformed YAML input
package graphdef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAMLInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte(`{{{not yaml`), 0644)
	_, err := ParseYAML(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
