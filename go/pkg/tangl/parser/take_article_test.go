//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeArticle — tests takeArticle for 'a'/'an'/unmatched branches
package parser

import "testing"

func TestTakeArticle(t *testing.T) {
	t.Run("A", func(t *testing.T) {
		rest, ok := takeArticle("a general rule")
		if !ok {
			t.Fatal("expected ok=true for 'a'")
		}
		if rest != " general rule" {
			t.Errorf("expected rest=' general rule', got %q", rest)
		}
	})

	t.Run("An", func(t *testing.T) {
		rest, ok := takeArticle("an except rule")
		if !ok {
			t.Fatal("expected ok=true for 'an'")
		}
		if rest != " except rule" {
			t.Errorf("expected rest=' except rule', got %q", rest)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		rest, ok := takeArticle("the general rule")
		if ok {
			t.Fatal("expected ok=false for an unrecognized article")
		}
		if rest != "the general rule" {
			t.Errorf("expected rest='the general rule', got %q", rest)
		}
	})
}
