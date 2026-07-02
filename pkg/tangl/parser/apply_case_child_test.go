//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyCaseChild — tests applyCaseChild for require/node/attack/exec success branches, plus the unrecognized fallback
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplyCaseChild(t *testing.T) {
	t.Run("RequireSuccess", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "`amount` is required", Line: 1}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Requires) != 1 || c.Requires[0].Field != "amount" {
			t.Fatalf("expected Requires to contain amount, got %+v", c.Requires)
		}
	})

	t.Run("NodeSuccess", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "`n1` is a general rule", Line: 3}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Nodes) != 1 || c.Nodes[0].Name != "n1" {
			t.Fatalf("expected Nodes to contain n1, got %+v", c.Nodes)
		}
	})

	t.Run("AttackSuccess", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "don't `target` when `attacker`", Line: 5}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Attacks) != 1 || c.Attacks[0].Target != "target" || c.Attacks[0].Attacker != "attacker" {
			t.Fatalf("expected Attacks to contain target/attacker, got %+v", c.Attacks)
		}
	})

	t.Run("ExecSuccess", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "do `act1` when `node1`", Line: 7}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Execs) != 1 || c.Execs[0].Kind != ast.DoExec || c.Execs[0].Node != "node1" {
			t.Fatalf("expected Execs to contain a do exec on node1, got %+v", c.Execs)
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "some unrecognized statement", Line: 9}, "test.md")
		if err == nil {
			t.Fatal("expected an error for an unrecognized case statement")
		}
		if !strings.Contains(err.Error(), "unrecognized case statement") {
			t.Errorf("expected unrecognized statement error, got %v", err)
		}
	})
}
