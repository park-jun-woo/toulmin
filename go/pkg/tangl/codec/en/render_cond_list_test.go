//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCondList — renderCondList flattens e and renders one list line per flattened term
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCondList checks the loop over flattenLogic's result: a single
// unprefixed term (one iteration), a two-term and/or chain (two iterations),
// and a nested chain (three iterations, exercising append across iterations).
func TestRenderCondList(t *testing.T) {
	a := ast.Compare{Field: "a", Op: "equals", Value: "1"}
	b := ast.Compare{Field: "b", Op: "equals", Value: "2"}
	c := ast.Compare{Field: "c", Op: "equals", Value: "3"}

	got := renderCondList(a, 0)
	want := []string{"- `a` equals 1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("single term: got %v, want %v", got, want)
	}

	two := ast.Logic{Op: "and", Terms: []ast.Expr{a, b}}
	got = renderCondList(two, 0)
	want = []string{"- `a` equals 1", "- and `b` equals 2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("two terms: got %v, want %v", got, want)
	}

	nested := ast.Logic{Op: "or", Terms: []ast.Expr{two, c}}
	got = renderCondList(nested, 1)
	want = []string{
		"  - `a` equals 1",
		"  - and `b` equals 2",
		"  - or `c` equals 3",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nested terms: got %v, want %v", got, want)
	}
}
