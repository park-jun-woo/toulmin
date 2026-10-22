//ff:func feature=tangl type=parser control=sequence
//ff:what TestParsePath — tests parsePath for single-step, multi-step, no-backtick, and malformed-chain branches
package parser

import (
	"testing"
)

func TestParsePath(t *testing.T) {
	t.Run("SingleStep", func(t *testing.T) {
		p, rest, ok := parsePath("`author` rest")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if len(p.Steps) != 1 || p.Steps[0] != "author" {
			t.Fatalf("expected [author], got %+v", p.Steps)
		}
		if rest != " rest" {
			t.Errorf("expected rest=' rest', got %q", rest)
		}
	})

	t.Run("MultiStep", func(t *testing.T) {
		p, rest, ok := parsePath("`author`.`trustScore` is at least 0.7")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if len(p.Steps) != 2 || p.Steps[0] != "author" || p.Steps[1] != "trustScore" {
			t.Fatalf("expected [author trustScore], got %+v", p.Steps)
		}
		if rest != " is at least 0.7" {
			t.Errorf("expected rest=' is at least 0.7', got %q", rest)
		}
	})

	t.Run("ThreeStep", func(t *testing.T) {
		p, _, ok := parsePath("`a`.`b`.`c`")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if len(p.Steps) != 3 {
			t.Fatalf("expected 3 steps, got %+v", p.Steps)
		}
	})

	t.Run("NoBacktick", func(t *testing.T) {
		_, rest, ok := parsePath("author")
		if ok {
			t.Fatal("expected ok=false for text with no backtick")
		}
		if rest != "author" {
			t.Errorf("expected original text as rest, got %q", rest)
		}
	})

	t.Run("MalformedChain", func(t *testing.T) {
		_, rest, ok := parsePath("`author`.notbacktick")
		if ok {
			t.Fatal("expected ok=false for a malformed dot-chain")
		}
		if rest != "`author`.notbacktick" {
			t.Errorf("expected original text as rest, got %q", rest)
		}
	})
}
