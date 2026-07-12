//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileCompare — tests lhs/rhs operand forms, value-call errors, and Fact fields
package eval

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompileCompare(t *testing.T) {
	boom := errors.New("boom")
	res := fakeResolver{values: map[string]func(toulmin.Context, toulmin.Specs) (any, error){
		"\x00bad": func(toulmin.Context, toulmin.Specs) (any, error) { return nil, boom },
	}}

	t.Run("FieldVsLiteral", func(t *testing.T) {
		b, ev := runRule(t, ast.Compare{Field: "x", Op: "equals", Value: "1"}, res, testCtx())
		f := asFact(t, ev)
		if !b || f.Where != "x" || f.Expected != "equals 1" || f.Actual != "1" {
			t.Fatalf("field-vs-literal fact mismatch: %v %+v", b, f)
		}
	})
	t.Run("LhsVsRhsOperands", func(t *testing.T) {
		c := ast.Compare{Lhs: ast.Path{Steps: []string{"x"}}, Op: "equals", Rhs: ast.Path{Steps: []string{"x"}}}
		b, ev := runRule(t, c, res, testCtx())
		f := asFact(t, ev)
		if !b || f.Where != "x" || f.Expected != "equals x" || f.Actual != "1" {
			t.Fatalf("operand-vs-operand fact mismatch: %v %+v", b, f)
		}
	})
	t.Run("LhsValueCallError", func(t *testing.T) {
		c := ast.Compare{Lhs: ast.Call{Fn: ast.Ref{Name: "bad"}}, Op: "is empty"}
		b, ev := runRule(t, c, res, testCtx())
		f := asFact(t, ev)
		if !b || f.Expected != "is empty" || f.Actual != boom {
			t.Fatalf("erroring lhs call should carry the error as Actual: %v %+v", b, f)
		}
	})
	t.Run("MissingFieldNilActual", func(t *testing.T) {
		b, ev := runRule(t, ast.Compare{Field: "ghost", Op: "equals", Value: "1"}, res, testCtx())
		f := asFact(t, ev)
		if b || f.Actual != nil {
			t.Fatalf("missing field should be false with nil Actual: %v %+v", b, f)
		}
	})
	t.Run("UnresolvedLhsCall", func(t *testing.T) {
		c := ast.Compare{Lhs: ast.Call{Fn: ast.Ref{Name: "nope"}}, Op: "is empty"}
		if _, err := compileCompare(c, res); err == nil {
			t.Fatal("want unresolved lhs value error")
		}
	})
	t.Run("UnresolvedRhsCall", func(t *testing.T) {
		c := ast.Compare{Field: "x", Op: "equals", Rhs: ast.Call{Fn: ast.Ref{Name: "nope"}}}
		if _, err := compileCompare(c, res); err == nil {
			t.Fatal("want unresolved rhs value error")
		}
	})
}
