//ff:func feature=tangl type=util control=sequence
//ff:what checkGenParity — assert gen.Generate is byte-identical and gofmt-clean for one fixture
package tangeul

import (
	"go/format"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
)

// checkGenParity generates Go source from both the text-parse AST and the reloaded
// golden .tangeul AST of the named fixture, failing if they differ or if the output
// is not gofmt-parseable. gen emits only import paths (never Document.Path) and
// ignores Line, so no normalization is needed.
func checkGenParity(t *testing.T, name string) {
	t.Helper()
	old, loaded := astPair(t, name)
	oldCode, err := gen.Generate(old)
	if err != nil {
		t.Fatalf("%s: Generate(old): %v", name, err)
	}
	newCode, err := gen.Generate(loaded)
	if err != nil {
		t.Fatalf("%s: Generate(new): %v", name, err)
	}
	if oldCode != newCode {
		t.Fatalf("%s: gen output differs between text-parse and reloaded AST", name)
	}
	if _, err := format.Source([]byte(newCode)); err != nil {
		t.Fatalf("%s: gen output is not gofmt-parseable: %v", name, err)
	}
}
