//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseSupersedes — tests registerCaseSupersedes for empty edges and multiple edges emitted in order
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRegisterCaseSupersedes(t *testing.T) {
	nodes := map[string]nodeInfo{
		"a": {Var: "nA"},
		"b": {Var: "nB"},
		"c": {Var: "nC"},
	}

	t.Run("no supersedes emits trailing newline only", func(t *testing.T) {
		var w strings.Builder
		registerCaseSupersedes(&w, ast.Case{}, nodes)
		if w.String() != "\n" {
			t.Errorf("got %q, want %q", w.String(), "\n")
		}
	})

	t.Run("multiple supersedes emitted in order", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{
			Supersedes: []ast.Supersede{
				{Superseder: "a", Superseded: "b"},
				{Superseder: "b", Superseded: "c"},
			},
		}
		registerCaseSupersedes(&w, c, nodes)
		want := "\tnA.Supersedes(nB)\n\tnB.Supersedes(nC)\n\n"
		if w.String() != want {
			t.Errorf("got %q, want %q", w.String(), want)
		}
	})

	t.Run("DuplicatePairEmittedOnce", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{Supersedes: []ast.Supersede{
			{Superseder: "a", Superseded: "b"},
			{Superseder: "a", Superseded: "b"},
		}}
		nodes := map[string]nodeInfo{"a": {Var: "nA"}, "b": {Var: "nB"}}
		registerCaseSupersedes(&w, c, nodes)
		if got, want := w.String(), "\tnA.Supersedes(nB)\n\n"; got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})
}
