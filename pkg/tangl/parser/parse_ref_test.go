//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRef — tests parseRef for no-backtick, bare-name, qualified-success, and qualified-failure branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRef(t *testing.T) {
	t.Run("NoBacktick", func(t *testing.T) {
		ref, rest, ok := parseRef("notbacktick")
		if ok {
			t.Fatal("expected ok=false for text with no backtick")
		}
		if ref != (ast.Ref{}) {
			t.Errorf("expected zero-value Ref, got %+v", ref)
		}
		if rest != "notbacktick" {
			t.Errorf("expected rest='notbacktick', got %q", rest)
		}
	})

	t.Run("BareName", func(t *testing.T) {
		ref, rest, ok := parseRef("`credit` extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if ref.Alias != "" || ref.Name != "credit" {
			t.Errorf("expected {Name:credit}, got %+v", ref)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("QualifiedSuccess", func(t *testing.T) {
		ref, rest, ok := parseRef("`credit`.`Threshold` extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if ref.Alias != "credit" || ref.Name != "Threshold" {
			t.Errorf("expected {Alias:credit Name:Threshold}, got %+v", ref)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("QualifiedFailure", func(t *testing.T) {
		ref, rest, ok := parseRef("`credit`.notbacktick")
		if ok {
			t.Fatal("expected ok=false for a malformed qualified reference")
		}
		if ref != (ast.Ref{}) {
			t.Errorf("expected zero-value Ref, got %+v", ref)
		}
		if rest != "`credit`.notbacktick" {
			t.Errorf("expected original string returned as rest, got %q", rest)
		}
	})
}
