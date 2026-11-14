//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeKeyword — tests takeKeyword for prefix-mismatch, word-boundary-mismatch, exact-match, and trailing-boundary branches
package parser

import "testing"

func TestTakeKeyword(t *testing.T) {
	t.Run("NoPrefixMatch", func(t *testing.T) {
		rest, ok := takeKeyword("  hello world", "case")
		if ok {
			t.Fatal("expected no match when prefix does not match")
		}
		if rest != "hello world" {
			t.Fatalf("expected rest=%q, got %q", "hello world", rest)
		}
	})

	t.Run("PrefixMatchButNoWordBoundary", func(t *testing.T) {
		rest, ok := takeKeyword("don't stop", "do")
		if ok {
			t.Fatal("expected no match when following byte is not a boundary")
		}
		if rest != "don't stop" {
			t.Fatalf("expected rest=%q, got %q", "don't stop", rest)
		}
	})

	t.Run("MatchWithSpaceBoundary", func(t *testing.T) {
		rest, ok := takeKeyword("case of X", "case")
		if !ok {
			t.Fatal("expected match")
		}
		if rest != " of X" {
			t.Fatalf("expected rest=%q, got %q", " of X", rest)
		}
	})

	t.Run("MatchWithTabBoundary", func(t *testing.T) {
		rest, ok := takeKeyword("case\tof X", "case")
		if !ok {
			t.Fatal("expected match")
		}
		if rest != "\tof X" {
			t.Fatalf("expected rest=%q, got %q", "\tof X", rest)
		}
	})

	t.Run("MatchWithEmptyRest", func(t *testing.T) {
		rest, ok := takeKeyword("  case", "case")
		if !ok {
			t.Fatal("expected match when keyword consumes entire remaining string")
		}
		if rest != "" {
			t.Fatalf("expected empty rest, got %q", rest)
		}
	})
}
