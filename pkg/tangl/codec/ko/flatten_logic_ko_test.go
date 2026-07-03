//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestFlattenLogicKo — flattenLogicKo flattens a right-leaning binary Logic chain, picks 그리고/또는 by Op, and passes through non-binary-Logic terms
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestFlattenLogicKo checks the two flattenLogicKo branches: a non-Logic leaf (and a
// Logic whose Terms length isn't 2) fall through to the single-item base case, while a
// binary Logic recurses on its left term and appends its right term with the
// "and"→그리고 / "or"→또는 prefix, including a chain deep enough to exercise the
// recursion more than once.
func TestFlattenLogicKo(t *testing.T) {
	a := ast.Compare{Field: "a", Op: "equals", Value: "1"}
	b := ast.Compare{Field: "b", Op: "equals", Value: "2"}
	c := ast.Compare{Field: "c", Op: "equals", Value: "3"}

	cases := []struct {
		name string
		e    ast.Expr
		want []condTermKo
	}{
		{"Leaf", a, []condTermKo{{Prefix: "", Term: a}}},
		{
			"NonBinaryLogic",
			ast.Logic{Op: "and", Terms: []ast.Expr{a, b, c}},
			[]condTermKo{{Prefix: "", Term: ast.Logic{Op: "and", Terms: []ast.Expr{a, b, c}}}},
		},
		{
			"And",
			ast.Logic{Op: "and", Terms: []ast.Expr{a, b}},
			[]condTermKo{{Prefix: "", Term: a}, {Prefix: "그리고", Term: b}},
		},
		{
			"Or",
			ast.Logic{Op: "or", Terms: []ast.Expr{a, b}},
			[]condTermKo{{Prefix: "", Term: a}, {Prefix: "또는", Term: b}},
		},
		{
			"ChainedRecursion",
			ast.Logic{Op: "and", Terms: []ast.Expr{
				ast.Logic{Op: "or", Terms: []ast.Expr{a, b}},
				c,
			}},
			[]condTermKo{
				{Prefix: "", Term: a},
				{Prefix: "또는", Term: b},
				{Prefix: "그리고", Term: c},
			},
		},
	}
	for _, c := range cases {
		if got := flattenLogicKo(c.e); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: flattenLogicKo(...) = %+v, want %+v", c.name, got, c.want)
		}
	}
}
