//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLFileNotFound — tests error for missing YAML file
package toulmin

import "testing"

func TestParseYAMLFileNotFound(t *testing.T) {
	_, err := ParseYAML("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
