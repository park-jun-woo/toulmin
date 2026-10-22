//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCondListKo — renderCondListKo appends every flattened term's rendered lines and propagates a term's render error
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCondListKo checks the loop over flattenLogicKo(e): a binary Logic
// flattens to two terms whose rendered lines are concatenated in order (looping more
// than once), and a term that renderTermKo cannot render (here, a raw non-binary
// Logic, which is neither a Compare/Not leaf nor an Either) makes the loop return the
// error immediately instead of any partial lines.
func TestRenderCondListKo(t *testing.T) {
	a := ast.Compare{Field: "a", Op: "equals", Value: "true"}
	b := ast.Compare{Field: "b", Op: "equals", Value: "true"}

	t.Run("Ok", func(t *testing.T) {
		e := ast.Logic{Op: "and", Terms: []ast.Expr{a, b}}
		got, err := renderCondListKo(e, 0)
		if err != nil {
			t.Fatalf("renderCondListKo: %v", err)
		}
		want := []string{"- `a` 면", "- 그리고 `b` 면"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("renderCondListKo = %+v, want %+v", got, want)
		}
	})

	t.Run("Err", func(t *testing.T) {
		e := ast.Logic{Op: "and", Terms: []ast.Expr{a, b, a}}
		if _, err := renderCondListKo(e, 0); err == nil {
			t.Error("renderCondListKo accepted an unsupported condition term")
		}
	})
}
