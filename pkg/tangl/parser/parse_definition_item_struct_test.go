//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionItem_Struct — tests parseDefinitionItem for StructDef branches (no children, child success, child error)
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseDefinitionItem_Struct(t *testing.T) {
	t.Run("NoChildren", func(t *testing.T) {
		it := item{Text: "`credit` means", Line: 3}
		def, err := parseDefinitionItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Kind != ast.StructDef {
			t.Errorf("expected StructDef, got %v", def.Kind)
		}
		if def.Name != "credit" {
			t.Errorf("expected Name=credit, got %q", def.Name)
		}
		if len(def.Fields) != 0 {
			t.Errorf("expected 0 fields, got %d", len(def.Fields))
		}
	})

	t.Run("WithChildSuccess", func(t *testing.T) {
		it := item{
			Text: "`credit` means",
			Line: 4,
			Children: []item{
				{Text: "has `age` as Int", Line: 5},
			},
		}
		def, err := parseDefinitionItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Kind != ast.StructDef {
			t.Errorf("expected StructDef, got %v", def.Kind)
		}
		if len(def.Fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(def.Fields))
		}
		if def.Fields[0].Name != "age" || def.Fields[0].Type != "Int" {
			t.Errorf("unexpected field: %+v", def.Fields[0])
		}
	})

	t.Run("WithChildError", func(t *testing.T) {
		it := item{
			Text: "`credit` means",
			Line: 6,
			Children: []item{
				{Text: "invalid child text", Line: 7},
			},
		}
		_, err := parseDefinitionItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseFieldItem")
		}
		if !strings.Contains(err.Error(), "expected 'has `field` as Type'") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
