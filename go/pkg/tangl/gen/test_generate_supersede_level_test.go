//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGenerateSupersedeLevel — generates a supersede+level fixture and checks gofmt-clean output with the new API calls
package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestGenerateSupersedeLevel parses a fixture that exercises a supersede
// edge and both crisp levels, then checks the generated code re-formats
// cleanly and contains the Supersedes and Level API calls.
func TestGenerateSupersedeLevel(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/supersede_level.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("re-format.Source: %v\n%s", err, src)
	}
	for _, want := range []string{
		".Supersedes(",
		".Level(toulmin.Fail)",
		".Level(toulmin.Review)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated code missing %q; got:\n%s", want, src)
		}
	}
}
