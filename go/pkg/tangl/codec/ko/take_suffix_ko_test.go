//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeSuffixKo — trailing keyword consumption, longest-match, boundary rejection inside 같으면, backtick/start boundaries
package ko

import "testing"

func TestTakeSuffixKo(t *testing.T) {
	t.Run("StripsTrailingKeywordAndSpace", func(t *testing.T) {
		got, ok := takeSuffixKo("`X` 면 `F` 실행한다", "실행한다")
		if !ok {
			t.Fatal("expected suffix match")
		}
		if got != "`X` 면 `F`" {
			t.Fatalf("got %q, want %q", got, "`X` 면 `F`")
		}
	})

	t.Run("LongestCompoundMatch", func(t *testing.T) {
		got, ok := takeSuffixKo("`X` 면 `F` 한 번만 실행한다", "한 번만 실행한다")
		if !ok {
			t.Fatal("expected compound suffix match")
		}
		if got != "`X` 면 `F`" {
			t.Fatalf("got %q, want %q", got, "`X` 면 `F`")
		}
	})

	t.Run("RejectsInsideCompoundWord", func(t *testing.T) {
		// 면 inside 같으면 must not match: the preceding byte is not a boundary.
		if _, ok := takeSuffixKo("`f` 가 `g` 와 같으면", "면"); ok {
			t.Fatal("expected no match for 면 inside 같으면")
		}
	})

	t.Run("NoSuffix", func(t *testing.T) {
		if _, ok := takeSuffixKo("`X` 무효다", "실행한다"); ok {
			t.Fatal("expected no match when kw is not a suffix")
		}
	})

	t.Run("IdaConsumedWholeNotSplit", func(t *testing.T) {
		got, ok := takeSuffixKo("이 문서는 `아메리카노` 이다", "이다")
		if !ok {
			t.Fatal("expected 이다 to be consumed as a whole keyword")
		}
		if got != "이 문서는 `아메리카노`" {
			t.Fatalf("got %q, want %q", got, "이 문서는 `아메리카노`")
		}
	})

	t.Run("BacktickBoundary", func(t *testing.T) {
		got, ok := takeSuffixKo("`X`이다", "이다")
		if !ok {
			t.Fatal("expected match with backtick boundary")
		}
		if got != "`X`" {
			t.Fatalf("got %q, want %q", got, "`X`")
		}
	})

	t.Run("StartOfStringBoundary", func(t *testing.T) {
		got, ok := takeSuffixKo("경우에", "경우에")
		if !ok {
			t.Fatal("expected match when kw is the whole string")
		}
		if got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})
}
