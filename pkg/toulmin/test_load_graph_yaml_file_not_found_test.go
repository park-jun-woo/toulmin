//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraphYAMLFileNotFound — tests error for missing YAML file
package toulmin

import (
	"testing"
)

func TestLoadGraphYAMLFileNotFound(t *testing.T) {
	_, err := LoadGraphYAML("/nonexistent/path.yaml", nil, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
