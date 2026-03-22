//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLQualifierZero — tests explicit qualifier 0.0 is preserved
package graphdef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAMLQualifierZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(`
graph: test
rules:
  - name: W
    role: warrant
    qualifier: 0.0
`), 0644)
	def, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *def.Rules[0].Qualifier != 0.0 {
		t.Errorf("expected explicit 0.0, got %f", *def.Rules[0].Qualifier)
	}
}
