//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLQualifierDefault — tests default qualifier value is 1.0 when omitted
package graphdef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAMLQualifierDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(`
graph: test
rules:
  - name: W
    role: warrant
`), 0644)
	def, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *def.Rules[0].Qualifier != 1.0 {
		t.Errorf("expected default 1.0, got %f", *def.Rules[0].Qualifier)
	}
}
