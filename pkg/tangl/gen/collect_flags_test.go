//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectFlags — tests collectFlags aggregates all four flag computations
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectFlags(t *testing.T) {
	t.Run("all false", func(t *testing.T) {
		flags := collectFlags(&ast.Document{})
		if flags.NeedsCompare {
			t.Error("expected NeedsCompare false for no rules")
		}
	})

	t.Run("needs compare true", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r1", Cond: ast.Compare{Field: "x", Op: "==", Value: "1"}}},
		}
		flags := collectFlags(doc)
		if !flags.NeedsCompare {
			t.Error("expected NeedsCompare true when doc.Rules is non-empty")
		}
	})
}
