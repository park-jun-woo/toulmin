//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseWithTermsKo — single term, 그리고-chained terms, a following clause remainder, and the missing-backtick error
package ko

import "testing"

func TestParseWithTermsKo(t *testing.T) {
	t.Run("SingleTerm", func(t *testing.T) {
		terms, rest, err := parseWithTermsKo("`a`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(terms) != 1 || terms[0] != "a" || rest != "" {
			t.Fatalf("terms=%v rest=%q", terms, rest)
		}
	})

	t.Run("ChainedTerms", func(t *testing.T) {
		terms, rest, err := parseWithTermsKo("`a` 그리고 `b`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(terms) != 2 || terms[0] != "a" || terms[1] != "b" || rest != "" {
			t.Fatalf("terms=%v rest=%q", terms, rest)
		}
	})

	t.Run("RemainderAfterTerm", func(t *testing.T) {
		terms, rest, err := parseWithTermsKo("`a`, 가중 0.5", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(terms) != 1 || terms[0] != "a" || rest != ", 가중 0.5" {
			t.Fatalf("terms=%v rest=%q", terms, rest)
		}
	})

	t.Run("MissingBacktick", func(t *testing.T) {
		if _, _, err := parseWithTermsKo("그냥텍스트", 1, "t.md"); err == nil {
			t.Fatal("expected error for a missing backtick-quoted term")
		}
	})
}
