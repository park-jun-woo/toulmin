//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersedeForms — a pair edge renders "supersedes"; a chain head renders the ordered chain
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderSupersedeForms(t *testing.T) {
	if got, want := renderSupersede(ast.Supersede{Superseder: "A", Superseded: "B"}), "`A` supersedes `B`"; got != want {
		t.Fatalf("pair = %q, want %q", got, want)
	}
	chain := ast.Supersede{Superseder: "A", Superseded: "B", Chain: []string{"A", "B", "C"}}
	if got, want := renderSupersede(chain), "supersede in order: `A` > `B` > `C`"; got != want {
		t.Fatalf("chain = %q, want %q", got, want)
	}
	short := ast.Supersede{Superseder: "A", Superseded: "B", Chain: []string{"A", "B"}}
	if got := renderSupersede(short); got != "`A` supersedes `B`" {
		t.Fatalf("two-name chain must render as a pair, got %q", got)
	}
}
