//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraphYAML — tests loading graph definition from YAML file
package toulmin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGraphYAML(t *testing.T) {
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

	funcs := map[string]any{
		"isAuthenticated": func(claim, ground, backing any) (bool, any) { return true, nil },
		"isIPBlocked":     func(claim, ground, backing any) (bool, any) { return true, nil },
	}

	g, err := LoadGraphYAML(path, funcs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("evaluate error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].Verdict >= 1.0 || results[0].Verdict <= -1.0 {
		t.Errorf("unexpected verdict: %f", results[0].Verdict)
	}
}
