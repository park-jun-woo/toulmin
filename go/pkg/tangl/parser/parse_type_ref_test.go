//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseTypeRef — tests parseTypeRef for backtick-qualified, backtick-bare, backtick-error, plain, and empty branches
package parser

import "testing"

func TestParseTypeRef(t *testing.T) {
	t.Run("BacktickError", func(t *testing.T) {
		typ, rest, ok := parseTypeRef("`unterminated")
		if ok {
			t.Fatal("expected ok=false for an unterminated backtick reference")
		}
		if typ != "" {
			t.Errorf("expected empty typ, got %q", typ)
		}
		if rest != "`unterminated" {
			t.Errorf("expected rest='`unterminated', got %q", rest)
		}
	})

	t.Run("BacktickQualified", func(t *testing.T) {
		typ, rest, ok := parseTypeRef("`credit`.`Threshold` extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if typ != "credit.Threshold" {
			t.Errorf("expected typ='credit.Threshold', got %q", typ)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("BacktickBare", func(t *testing.T) {
		typ, rest, ok := parseTypeRef("`Currency` extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if typ != "Currency" {
			t.Errorf("expected typ='Currency', got %q", typ)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("PlainSuccess", func(t *testing.T) {
		typ, rest, ok := parseTypeRef("Int extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if typ != "Int" {
			t.Errorf("expected typ='Int', got %q", typ)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		typ, rest, ok := parseTypeRef("   ")
		if ok {
			t.Fatal("expected ok=false for an empty/whitespace-only string")
		}
		if typ != "" {
			t.Errorf("expected empty typ, got %q", typ)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})
}
