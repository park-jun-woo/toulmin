//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyCaseChild_Errors — tests applyCaseChild error propagation from require/node/attack/exec sub-parsers
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplyCaseChild_Errors(t *testing.T) {
	t.Run("RequireError", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "`amount` is required as", Line: 2}, "test.md")
		if err == nil {
			t.Fatal("expected an error for a missing type after 'as'")
		}
		if !strings.Contains(err.Error(), "expected type after 'as'") {
			t.Errorf("expected type-after-as error, got %v", err)
		}
	})

	t.Run("NodeError", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "`n1` is a foo", Line: 4}, "test.md")
		if err == nil {
			t.Fatal("expected an error for an unrecognized role expression")
		}
		if !strings.Contains(err.Error(), "expected role") {
			t.Errorf("expected role error, got %v", err)
		}
	})

	t.Run("AttackError", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "don't `target`", Line: 6}, "test.md")
		if err == nil {
			t.Fatal("expected an error for a missing 'when' in don't edge")
		}
		if !strings.Contains(err.Error(), "expected 'when' in don't edge") {
			t.Errorf("expected when error, got %v", err)
		}
	})

	t.Run("ExecError", func(t *testing.T) {
		c := &ast.Case{}
		err := applyCaseChild(c, item{Text: "do `act1`", Line: 8}, "test.md")
		if err == nil {
			t.Fatal("expected an error for a missing 'when' in do edge")
		}
		if !strings.Contains(err.Error(), "expected 'when' in do edge") {
			t.Errorf("expected when error, got %v", err)
		}
	})
}
