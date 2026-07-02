//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseFieldItem — tests parseFieldItem for has/name/as/type/trailing branches
package parser

import (
	"strings"
	"testing"
)

func TestParseFieldItem(t *testing.T) {
	t.Run("NoHas", func(t *testing.T) {
		it := item{Text: "`age` as Int", Line: 1}
		_, err := parseFieldItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'has'")
		}
		if !strings.Contains(err.Error(), "expected 'has `field` as Type'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoBacktickName", func(t *testing.T) {
		it := item{Text: "has age as Int", Line: 2}
		_, err := parseFieldItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted field name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted field name after 'has'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoAs", func(t *testing.T) {
		it := item{Text: "has `age` Int", Line: 3}
		_, err := parseFieldItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'as'")
		}
		if !strings.Contains(err.Error(), "expected 'as Type' after field name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoType", func(t *testing.T) {
		it := item{Text: "has `age` as ", Line: 4}
		_, err := parseFieldItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing type")
		}
		if !strings.Contains(err.Error(), "expected type after 'as'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		it := item{Text: "has `age` as Int extra", Line: 5}
		_, err := parseFieldItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		it := item{Text: "has `age` as Int", Line: 6}
		f, err := parseFieldItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Name != "age" || f.Type != "Int" {
			t.Errorf("expected {Name:age Type:Int}, got %+v", f)
		}
	})
}
