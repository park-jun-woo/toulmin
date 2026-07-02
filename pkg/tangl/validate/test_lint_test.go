//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestLint — Internal-reachable unguarded do edges warn, others stay silent
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestLint checks that the americano fixture (which has a tangl:Internal
// tick loop reaching a `do log.CoffeeStatus when 'brewing complete'` with no
// `once` guard) produces at least one warning, while the access and transfer
// fixtures (no tangl:Internal section at all) produce none.
func TestLint(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/americano.md")
	if err != nil {
		t.Fatalf("Parse(americano): %v", err)
	}
	warnings := Lint(doc)
	if len(warnings) == 0 {
		t.Error("Lint(americano) = [], want at least one unguarded-do warning")
	}

	for _, p := range []string{"../parser/testdata/access.md", "../parser/testdata/transfer.md"} {
		doc, err := parser.Parse(p)
		if err != nil {
			t.Fatalf("Parse(%s): %v", p, err)
		}
		if w := Lint(doc); len(w) != 0 {
			t.Errorf("Lint(%s) = %v, want no warnings (no tangl:Internal section)", p, w)
		}
	}
}
