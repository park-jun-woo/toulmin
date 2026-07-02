//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionItem — tests parseDefinitionItem for name/means/const branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseDefinitionItem(t *testing.T) {
	t.Run("NoBacktickName", func(t *testing.T) {
		it := item{Text: "credit means something", Line: 1}
		_, err := parseDefinitionItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted term name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted term name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoMeans", func(t *testing.T) {
		it := item{Text: "`credit` is something", Line: 2}
		_, err := parseDefinitionItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'means'")
		}
		if !strings.Contains(err.Error(), "expected 'means' after term name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ConstError", func(t *testing.T) {
		it := item{Text: "`credit` means 650 as `credit", Line: 8}
		_, err := parseDefinitionItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseConstValue")
		}
		if !strings.Contains(err.Error(), "expected reference after 'as'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ConstSuccess", func(t *testing.T) {
		it := item{Text: "`credit` means 650 as `credit`.`Threshold`", Line: 9}
		def, err := parseDefinitionItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Kind != ast.ConstDef {
			t.Errorf("expected ConstDef, got %v", def.Kind)
		}
		if def.Value != "650" {
			t.Errorf("expected Value=650, got %q", def.Value)
		}
		if def.SpecRef == nil || def.SpecRef.Alias != "credit" || def.SpecRef.Name != "Threshold" {
			t.Errorf("unexpected SpecRef: %+v", def.SpecRef)
		}
	})
}
