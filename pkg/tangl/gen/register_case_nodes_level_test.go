//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseNodes_Level — tests registerCaseNodes Level emission for fail, review, and unset
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRegisterCaseNodes_Level(t *testing.T) {
	t.Run("level fail node emits Level call", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "nodeF", Role: ast.CounterRule, Level: ast.LevelFail}}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), ".Level(toulmin.Fail)") {
			t.Errorf("expected .Level(toulmin.Fail), got %q", w.String())
		}
	})

	t.Run("level review node emits Level call", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "nodeR", Role: ast.ExceptRule, Level: ast.LevelReview}}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), ".Level(toulmin.Review)") {
			t.Errorf("expected .Level(toulmin.Review), got %q", w.String())
		}
	})

	t.Run("unset level emits no Level call", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "nodeU", Role: ast.GeneralRule}}}
		registerCaseNodes(&w, gc, c)
		if strings.Contains(w.String(), ".Level(") {
			t.Errorf("expected no .Level call for unset level, got %q", w.String())
		}
	})
}
