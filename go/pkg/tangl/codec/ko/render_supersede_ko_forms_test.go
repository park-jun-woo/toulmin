//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersedeKoForms — pair edges use a 받침-selected 면/이면; chain heads render 순서로 배제한다
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderSupersedeKoForms(t *testing.T) {
	if got, want := renderSupersedeKo(ast.Supersede{Superseder: "A", Superseded: "B"}), "`A` 면 `B` 배제한다"; got != want {
		t.Fatalf("pair = %q, want %q", got, want)
	}
	chain := ast.Supersede{Superseder: "A", Superseded: "B", Chain: []string{"A", "B", "C"}}
	if got, want := renderSupersedeKo(chain), "`A` > `B` > `C` 순서로 배제한다"; got != want {
		t.Fatalf("chain = %q, want %q", got, want)
	}
}
