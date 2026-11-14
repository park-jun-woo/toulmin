//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRequireItem_Type — tests parseRequireItem for trailing-text/as/type branches
package parser

import (
	"strings"
	"testing"
)

func TestParseRequireItem_Type(t *testing.T) {
	t.Run("NoAsTrailingText", func(t *testing.T) {
		it := item{Text: "`field` is required extra", Line: 5}
		_, matched, err := parseRequireItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for trailing text without 'as'")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("AsNoType", func(t *testing.T) {
		it := item{Text: "`field` is required as ", Line: 6}
		_, matched, err := parseRequireItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for missing type after 'as'")
		}
		if !strings.Contains(err.Error(), "expected type after 'as'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TypeTrailingText", func(t *testing.T) {
		it := item{Text: "`field` is required as Int extra", Line: 7}
		_, matched, err := parseRequireItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for trailing text after type")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TypeSuccess", func(t *testing.T) {
		it := item{Text: "`field` is required as Int", Line: 8}
		req, matched, err := parseRequireItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.Field != "field" || req.Type != "Int" {
			t.Errorf("unexpected require: %+v", req)
		}
	})
}
