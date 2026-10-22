//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeQuantifier — tests takeQuantifier for each quantifier keyword match and the no-match branch
package parser

import "testing"

func TestTakeQuantifier(t *testing.T) {
	t.Run("MatchAny", func(t *testing.T) {
		kw, rest, ok := takeQuantifier("any of `coupons`")
		if !ok || kw != "any" {
			t.Fatalf("expected any match, got kw=%q ok=%v", kw, ok)
		}
		if rest != " of `coupons`" {
			t.Errorf("expected rest=%q, got %q", " of `coupons`", rest)
		}
	})

	t.Run("MatchAll", func(t *testing.T) {
		kw, _, ok := takeQuantifier("all of `coupons`")
		if !ok || kw != "all" {
			t.Fatalf("expected all match, got kw=%q ok=%v", kw, ok)
		}
	})

	t.Run("MatchCount", func(t *testing.T) {
		kw, _, ok := takeQuantifier("count of `coupons`")
		if !ok || kw != "count" {
			t.Fatalf("expected count match, got kw=%q ok=%v", kw, ok)
		}
	})

	t.Run("MatchSum", func(t *testing.T) {
		kw, _, ok := takeQuantifier("sum of `coupons`")
		if !ok || kw != "sum" {
			t.Fatalf("expected sum match, got kw=%q ok=%v", kw, ok)
		}
	})

	t.Run("NoMatch", func(t *testing.T) {
		kw, rest, ok := takeQuantifier("`amount` equals 5")
		if ok || kw != "" {
			t.Fatalf("expected no match, got kw=%q ok=%v", kw, ok)
		}
		if rest != "`amount` equals 5" {
			t.Errorf("expected original text as rest, got %q", rest)
		}
	})
}
