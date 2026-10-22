//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSupersedeItem — tests parseSupersedeItem for non-match, missing target, trailing text, success, and don't non-collision
package parser

import (
	"strings"
	"testing"
)

func TestParseSupersedeItem(t *testing.T) {
	t.Run("NotSupersede", func(t *testing.T) {
		_, ok, err := parseSupersedeItem(item{Text: "`x` is a general rule", Line: 1}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a non-supersede statement")
		}
	})

	t.Run("NoLeadingBacktick", func(t *testing.T) {
		_, ok, err := parseSupersedeItem(item{Text: "don't `a` when `b`", Line: 2}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a don't edge (no collision with supersede)")
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		_, ok, err := parseSupersedeItem(item{Text: "`x` supersedes notabacktick", Line: 3}, "test.md")
		if !ok {
			t.Fatal("expected ok=true once supersedes keyword matched")
		}
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted target after 'supersedes'") {
			t.Fatalf("expected target error, got %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, ok, err := parseSupersedeItem(item{Text: "`x` supersedes `y` extra", Line: 4}, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Fatalf("expected trailing text error, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		s, ok, err := parseSupersedeItem(item{Text: "`x` supersedes `y`", Line: 5}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if s.Superseder != "x" || s.Superseded != "y" || s.Line != 5 {
			t.Fatalf("expected Superseder/Superseded/Line set, got %+v", s)
		}
	})
}
