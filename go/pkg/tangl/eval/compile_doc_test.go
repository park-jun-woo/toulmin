//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileDoc — tests doc-local bare resolution, unresolved errors, cycles, determinism
package eval

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompileDoc(t *testing.T) {
	empty := fakeResolver{}

	t.Run("DocLocalBarePredicateFirst", func(t *testing.T) {
		// A references bare rule B; B is a doc rule the host does not know.
		doc := &ast.Document{
			Path: "d.md",
			Rules: []ast.InlineRule{
				{Name: "B", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}, Line: 1},
				{Name: "A", Cond: ast.Call{Fn: ast.Ref{Name: "B"}}, Line: 2},
			},
		}
		out, err := CompileDoc(doc, empty)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		if len(out) != 2 || out["A"] == nil || out["B"] == nil {
			t.Fatalf("expected A and B closures, got %v", out)
		}
		c := toulmin.NewContext()
		c.Set("x", "1")
		b, ev := out["A"](c, nil)
		f, ok := ev.(toulmin.Fact)
		if !b || !ok || f.Where != "x" {
			t.Fatalf("A should be true carrying B's leaf witness, got %v %+v", b, ev)
		}
	})

	t.Run("UnresolvedBareCallError", func(t *testing.T) {
		doc := &ast.Document{Path: "d.md", Rules: []ast.InlineRule{
			{Name: "A", Cond: ast.Call{Fn: ast.Ref{Name: "ghost"}}, Line: 1},
		}}
		if _, err := CompileDoc(doc, empty); err == nil {
			t.Fatal("expected error for unresolved bare call")
		}
	})
	t.Run("UnresolvedAliasedCallError", func(t *testing.T) {
		doc := &ast.Document{Path: "d.md", Rules: []ast.InlineRule{
			{Name: "A", Cond: ast.Call{Fn: ast.Ref{Alias: "z", Name: "f"}}, Line: 1},
		}}
		if _, err := CompileDoc(doc, empty); err == nil {
			t.Fatal("expected error for unresolved aliased call")
		}
	})
	t.Run("NameCycleError", func(t *testing.T) {
		doc := &ast.Document{Path: "d.md", Rules: []ast.InlineRule{
			{Name: "A", Cond: ast.Call{Fn: ast.Ref{Name: "B"}}, Line: 1},
			{Name: "B", Cond: ast.Call{Fn: ast.Ref{Name: "A"}}, Line: 2},
		}}
		if _, err := CompileDoc(doc, empty); err == nil {
			t.Fatal("expected cycle error")
		}
	})

	t.Run("Determinism", func(t *testing.T) {
		doc := &ast.Document{Path: "d.md", Rules: []ast.InlineRule{
			{Name: "B", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}, Line: 1},
			{Name: "A", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
				ast.Call{Fn: ast.Ref{Name: "B"}},
				ast.Compare{Field: "y", Op: "is at least", Value: "5"},
			}}, Line: 2},
		}}
		out, err := CompileDoc(doc, empty)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		c := toulmin.NewContext()
		c.Set("x", "1")
		c.Set("y", "9")
		b1, e1 := out["A"](c, nil)
		b2, e2 := out["A"](c, nil)
		if b1 != b2 || !reflect.DeepEqual(e1, e2) {
			t.Fatalf("non-deterministic: (%v,%+v) vs (%v,%+v)", b1, e1, b2, e2)
		}
		if !b1 {
			t.Fatalf("expected A true, got %v", b1)
		}
	})
}
