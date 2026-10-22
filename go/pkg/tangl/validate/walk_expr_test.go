//ff:func feature=tangl type=validator control=sequence
//ff:what TestWalkExpr — tests walkExpr preorder coverage over nested Logic/Not/Either/Exists nodes and walkOperand nil-safety
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestWalkExpr(t *testing.T) {
	t.Run("NilRoot", func(t *testing.T) {
		count := 0
		walkExpr(nil, func(ast.Expr) { count++ })
		if count != 0 {
			t.Fatalf("expected 0 visits for nil root, got %d", count)
		}
	})

	t.Run("PreorderCoverage", func(t *testing.T) {
		tree := ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Not{Term: ast.Compare{Field: "a", Op: "is empty"}},
			ast.Either{Terms: []ast.Expr{
				ast.Exists{
					Source: ast.Path{Steps: []string{"xs"}},
					Binder: "x",
					Cond:   ast.Compare{Field: "x", Op: "is empty"},
				},
			}},
		}}
		var kinds []string
		walkExpr(tree, func(e ast.Expr) {
			switch e.(type) {
			case ast.Logic:
				kinds = append(kinds, "logic")
			case ast.Not:
				kinds = append(kinds, "not")
			case ast.Either:
				kinds = append(kinds, "either")
			case ast.Exists:
				kinds = append(kinds, "exists")
			case ast.Compare:
				kinds = append(kinds, "compare")
			}
		})
		// logic, not, compare, either, exists, compare
		if len(kinds) != 6 || kinds[0] != "logic" || kinds[3] != "either" || kinds[4] != "exists" {
			t.Fatalf("unexpected preorder visit sequence: %v", kinds)
		}
	})

	t.Run("WalkOperandNil", func(t *testing.T) {
		count := 0
		walkOperand(nil, func(ast.Operand) { count++ })
		if count != 0 {
			t.Fatalf("expected 0 visits for nil operand, got %d", count)
		}
	})

	t.Run("WalkOperandVisit", func(t *testing.T) {
		count := 0
		walkOperand(ast.Path{Steps: []string{"a"}}, func(ast.Operand) { count++ })
		if count != 1 {
			t.Fatalf("expected 1 visit, got %d", count)
		}
	})
}
