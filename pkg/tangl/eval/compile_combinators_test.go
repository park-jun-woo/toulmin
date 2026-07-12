//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileCombinators — tests and/or/either/not witness propagation and the compare leaf
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompileCombinators(t *testing.T) {
	truthy := ast.Compare{Field: "x", Op: "equals", Value: "1"}
	falsy := ast.Compare{Field: "y", Op: "equals", Value: "2"}
	empty := fakeResolver{}

	t.Run("AndFirstFalse", func(t *testing.T) {
		b, ev := runRule(t, ast.Logic{Op: "and", Terms: []ast.Expr{falsy, truthy}}, empty, testCtx())
		if b || asFact(t, ev).Where != "y" {
			t.Fatalf("want false with y witness, got %v %+v", b, ev)
		}
	})
	t.Run("AndAllTrueLastWitness", func(t *testing.T) {
		a := ast.Compare{Field: "x", Op: "equals", Value: "1"}
		bb := ast.Compare{Field: "y", Op: "equals", Value: "9"}
		b, ev := runRule(t, ast.Logic{Op: "and", Terms: []ast.Expr{a, bb}}, empty, testCtx())
		if !b || asFact(t, ev).Where != "y" {
			t.Fatalf("want true with last(y) witness, got %v %+v", b, ev)
		}
	})
	t.Run("OrFirstTrue", func(t *testing.T) {
		b, ev := runRule(t, ast.Logic{Op: "or", Terms: []ast.Expr{falsy, truthy}}, empty, testCtx())
		if !b || asFact(t, ev).Where != "x" {
			t.Fatalf("want true with x witness, got %v %+v", b, ev)
		}
	})
	t.Run("EitherFirstTrue", func(t *testing.T) {
		b, ev := runRule(t, ast.Either{Terms: []ast.Expr{falsy, truthy}}, empty, testCtx())
		if !b || asFact(t, ev).Where != "x" {
			t.Fatalf("want true with x witness, got %v %+v", b, ev)
		}
	})
	t.Run("NotPassThrough", func(t *testing.T) {
		b, ev := runRule(t, ast.Not{Term: falsy}, empty, testCtx())
		if !b || asFact(t, ev).Where != "y" {
			t.Fatalf("want true with unchanged y witness, got %v %+v", b, ev)
		}
	})
	t.Run("CompareLeafFact", func(t *testing.T) {
		_, ev := runRule(t, truthy, empty, testCtx())
		f := asFact(t, ev)
		if f.Where != "x" || f.Expected != "equals 1" || f.Actual != "1" {
			t.Fatalf("compare fact mismatch: %+v", f)
		}
	})
	t.Run("CompareLeafFactNoRhs", func(t *testing.T) {
		b, ev := runRule(t, ast.Compare{Field: "x", Op: "is not empty"}, empty, testCtx())
		f := asFact(t, ev)
		if !b || f.Where != "x" || f.Expected != "is not empty" || f.Actual != "1" {
			t.Fatalf("no-rhs compare fact mismatch: %v %+v", b, f)
		}
		b2, ev2 := runRule(t, ast.Compare{Field: "missing", Op: "is empty"}, empty, testCtx())
		f2 := asFact(t, ev2)
		if !b2 || f2.Where != "missing" || f2.Expected != "is empty" || f2.Actual != nil {
			t.Fatalf("no-rhs is-empty fact mismatch: %v %+v", b2, f2)
		}
	})
}
