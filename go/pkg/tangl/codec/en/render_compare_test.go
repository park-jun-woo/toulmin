//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCompare — renderCompare covers all eight operators incl. the value-less forms
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCompare checks each comparison operator, including "is empty" and
// "is not empty" which carry no value.
func TestRenderCompare(t *testing.T) {
	cases := []struct {
		c    ast.Compare
		want string
	}{
		{ast.Compare{Field: "f", Op: "is empty"}, "`f` is empty"},
		{ast.Compare{Field: "f", Op: "is not empty"}, "`f` is not empty"},
		{ast.Compare{Field: "f", Op: "equals", Value: `"x"`}, "`f` equals \"x\""},
		{ast.Compare{Field: "f", Op: "is in", Value: "`list`"}, "`f` is in `list`"},
		{ast.Compare{Field: "f", Op: "is greater than", Value: "5"}, "`f` is greater than 5"},
		{ast.Compare{Field: "f", Op: "is less than", Value: "2"}, "`f` is less than 2"},
		{ast.Compare{Field: "f", Op: "is at most", Value: "0"}, "`f` is at most 0"},
		{ast.Compare{Field: "f", Op: "is at least", Value: "1000"}, "`f` is at least 1000"},
	}
	for _, tc := range cases {
		if got := renderCompare(tc.c); got != tc.want {
			t.Errorf("renderCompare(%+v) = %q, want %q", tc.c, got, tc.want)
		}
	}
}
