//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectFlags — tests collectFlags aggregates document and rule flag scans
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCollectFlags(t *testing.T) {
	t.Run("all false", func(t *testing.T) {
		flags := collectFlags(&ast.Document{})
		if flags != (genFlags{}) {
			t.Errorf("expected zero flags for an empty document, got %+v", flags)
		}
	})

	t.Run("flat compare rule", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r1", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}}},
		}
		flags := collectFlags(doc)
		if !flags.NeedsCompareVal || !flags.NeedsCtxVal || !flags.NeedsWalk {
			t.Errorf("expected CompareVal/CtxVal/Walk for a flat compare, got %+v", flags)
		}
		if flags.NeedsAny || flags.NeedsHolds || flags.NeedsSpecVal || flags.NeedsItoa {
			t.Errorf("unexpected extra flags for a flat compare: %+v", flags)
		}
	})

	t.Run("exists with binder path and spec ref", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r2", Cond: ast.Exists{
				Source: ast.Path{Steps: []string{"coupons"}},
				Binder: "c",
				Cond: ast.Compare{
					Lhs: ast.Path{Steps: []string{"c", "minPrice"}},
					Op:  "is at most",
					Rhs: ast.SpecRef{Steps: []string{"discount", "max"}},
				},
			}}},
		}
		flags := collectFlags(doc)
		if !flags.NeedsAny || !flags.NeedsItoa || !flags.NeedsWalk || !flags.NeedsSpecVal || !flags.NeedsCtxVal {
			t.Errorf("expected Any/Itoa/Walk/SpecVal/CtxVal, got %+v", flags)
		}
	})

	t.Run("logic call and not", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r3", Cond: ast.Logic{Op: "or", Terms: []ast.Expr{
				ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}},
				ast.Not{Term: ast.Compare{Field: "x", Op: "is empty"}},
			}}}},
		}
		flags := collectFlags(doc)
		if !flags.NeedsOr || !flags.NeedsHolds || !flags.NeedsNot || !flags.NeedsCompareVal {
			t.Errorf("expected Or/Holds/Not/CompareVal, got %+v", flags)
		}
		if flags.NeedsAnd {
			t.Errorf("unexpected NeedsAnd: %+v", flags)
		}
	})
}
