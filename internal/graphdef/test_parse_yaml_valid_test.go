//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLValid — tests valid YAML parsing into GraphDef
package graphdef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAMLValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(`
graph: example
rules:
  - name: W
    role: warrant
    qualifier: 0.8
  - name: R
    role: rebuttal
defeats:
  - from: R
    to: W
`), 0644)
	def, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def.Graph != "example" {
		t.Errorf("graph name: expected 'example', got %q", def.Graph)
	}
	if len(def.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(def.Rules))
	}
	if *def.Rules[0].Qualifier != 0.8 {
		t.Errorf("W qualifier: expected 0.8, got %f", *def.Rules[0].Qualifier)
	}
	if len(def.Defeats) != 1 {
		t.Fatalf("expected 1 defeat, got %d", len(def.Defeats))
	}
}
