//ff:func feature=tangl type=validator control=sequence
//ff:what TestInternalReachableCases — tests internalReachableCases for empty, duplicate-seed, and BFS-revisit branches
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestInternalReachableCases(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		reached := internalReachableCases(doc)
		if len(reached) != 0 {
			t.Fatalf("expected empty set, got %v", reached)
		}
	})

	t.Run("DuplicateSeedAndBFSRevisit", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Runs: []string{"a", "a"}},
			},
			Cases: []ast.Case{
				{Name: "a", Execs: []ast.Exec{
					{Kind: ast.RunExec, Case: "c"},
				}},
				{Name: "c", Execs: []ast.Exec{
					// "a" is already reached: exercises the BFS continue branch.
					{Kind: ast.RunExec, Case: "a"},
					{Kind: ast.RunExec, Case: "d"},
				}},
				{Name: "d"},
			},
		}
		reached := internalReachableCases(doc)
		want := map[string]bool{"a": true, "c": true, "d": true}
		if len(reached) != len(want) {
			t.Fatalf("expected %v, got %v", want, reached)
		}
		for name := range want {
			if !reached[name] {
				t.Errorf("expected %q to be reached, got %v", name, reached)
			}
		}
	})
}
