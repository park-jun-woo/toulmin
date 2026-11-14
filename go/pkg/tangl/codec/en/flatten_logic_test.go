//ff:func feature=tangl type=parser control=sequence
//ff:what TestFlattenLogic — flattenLogic covers the Logic/non-Logic and Terms-length branches
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestFlattenLogic exercises the three ways flattenLogic's guard can resolve:
// e is not a Logic at all (type assertion fails), e is a Logic but does not
// have exactly two Terms (len check fails), and e is a two-term Logic (the
// recursive case), which in turn recurses down to a non-Logic base case.
func TestFlattenLogic(t *testing.T) {
	leaf := ast.Compare{Field: "a", Op: "equals", Value: "1"}

	// Not a Logic at all: type assertion fails, single unprefixed term.
	got := flattenLogic(leaf)
	want := []condTerm{{Prefix: "", Term: leaf}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("non-Logic: got %+v, want %+v", got, want)
	}

	// A Logic whose Terms length is not 2: guard's length check fails.
	single := ast.Logic{Op: "and", Terms: []ast.Expr{leaf}}
	got = flattenLogic(single)
	want = []condTerm{{Prefix: "", Term: single}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Logic len!=2 (1 term): got %+v, want %+v", got, want)
	}

	three := ast.Logic{Op: "and", Terms: []ast.Expr{leaf, leaf, leaf}}
	got = flattenLogic(three)
	want = []condTerm{{Prefix: "", Term: three}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Logic len!=2 (3 terms): got %+v, want %+v", got, want)
	}

	// A two-term Logic: recursive case, chained onto a non-Logic base case.
	b := ast.Compare{Field: "b", Op: "equals", Value: "2"}
	two := ast.Logic{Op: "or", Terms: []ast.Expr{leaf, b}}
	got = flattenLogic(two)
	want = []condTerm{{Prefix: "", Term: leaf}, {Prefix: "or", Term: b}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Logic len==2: got %+v, want %+v", got, want)
	}

	// Nested two-term Logics: the recursive call itself hits the recursive
	// branch again before bottoming out at the base case.
	c := ast.Compare{Field: "c", Op: "equals", Value: "3"}
	nested := ast.Logic{Op: "and", Terms: []ast.Expr{two, c}}
	got = flattenLogic(nested)
	want = []condTerm{
		{Prefix: "", Term: leaf},
		{Prefix: "or", Term: b},
		{Prefix: "and", Term: c},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nested Logic: got %+v, want %+v", got, want)
	}
}
