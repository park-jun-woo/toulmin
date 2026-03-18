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

func TestParseYAMLInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte(`{{{not yaml`), 0644)
	_, err := ParseYAML(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseYAMLFileNotFound(t *testing.T) {
	_, err := ParseYAML("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
