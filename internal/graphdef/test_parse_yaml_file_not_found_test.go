//ff:func feature=graph type=parser control=sequence
//ff:what TestParseYAMLFileNotFound — tests error returned for nonexistent file
package graphdef

import (
	"testing"
)

func TestParseYAMLFileNotFound(t *testing.T) {
	_, err := ParseYAML("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
