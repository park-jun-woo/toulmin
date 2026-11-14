//ff:func feature=tangl type=codegen control=sequence
//ff:what TestIndexRulesByName — rules are keyed by name and a later duplicate wins
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestIndexRulesByName(t *testing.T) {
	got := indexRulesByName([]ast.InlineRule{{Name: "a", Line: 1}, {Name: "b", Line: 2}, {Name: "a", Line: 3}})
	if len(got) != 2 || got["a"].Line != 3 || got["b"].Line != 2 {
		t.Fatalf("indexRulesByName = %#v", got)
	}
	if got := indexRulesByName(nil); len(got) != 0 {
		t.Fatalf("nil rules = %#v, want empty", got)
	}
}
