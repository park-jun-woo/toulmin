//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAML — tests YAML parsing into GraphDef
package toulmin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(`
graph: access
rules:
  - name: isAuthenticated
    role: warrant
  - name: isIPBlocked
    role: rebuttal
    qualifier: 0.8
defeats:
  - from: isIPBlocked
    to: isAuthenticated
`), 0644)

	def, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def.Graph != "access" {
		t.Errorf("graph name: expected 'access', got %q", def.Graph)
	}
	if len(def.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(def.Rules))
	}
	if def.Rules[1].Qualifier != 0.8 {
		t.Errorf("qualifier: expected 0.8, got %f", def.Rules[1].Qualifier)
	}
	if len(def.Defeats) != 1 {
		t.Fatalf("expected 1 defeat, got %d", len(def.Defeats))
	}
}
