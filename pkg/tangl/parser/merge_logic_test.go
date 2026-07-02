//ff:func feature=tangl type=parser control=sequence
//ff:what TestMergeLogic — tests mergeLogic builds an ast.Logic wrapping the left and right terms
package parser

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestMergeLogic(t *testing.T) {
	left := ast.Compare{Field: "amount", Op: "is empty"}
	right := ast.Compare{Field: "status", Op: "equals", Value: "done"}

	got := mergeLogic(left, "and", right)

	want := ast.Logic{Op: "and", Terms: []ast.Expr{left, right}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}
