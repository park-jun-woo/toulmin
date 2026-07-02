//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeItem_Clause — tests parseNodeItem for the matched no-clause/clause-error/clause-success branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseNodeItem_Clause(t *testing.T) {
	t.Run("NoClause", func(t *testing.T) {
		it := item{Text: "`node1` is a general rule", Line: 5}
		n, matched, err := parseNodeItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Name != "node1" || n.Role != ast.GeneralRule {
			t.Errorf("unexpected node: %+v", n)
		}
	})

	t.Run("ClauseError", func(t *testing.T) {
		it := item{Text: "`node1` is a general rule bogus clause", Line: 6}
		_, matched, err := parseNodeItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error from applyNodeClause")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ClauseSuccess", func(t *testing.T) {
		it := item{Text: "`node1` is a general rule checking `case1`", Line: 7}
		n, matched, err := parseNodeItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Name != "node1" || n.Role != ast.GeneralRule || n.Checking != "case1" {
			t.Errorf("unexpected node: %+v", n)
		}
	})
}
