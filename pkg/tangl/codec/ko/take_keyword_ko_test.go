//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeKeywordKo — prefix mismatch, word-boundary rejection, space/tab/empty-rest matches
package ko

import "testing"

func TestTakeKeywordKo(t *testing.T) {
	t.Run("NoPrefixMatch", func(t *testing.T) {
		rest, ok := takeKeywordKo("  또는 `X`", "그리고")
		if ok {
			t.Fatal("expected no match when prefix does not match")
		}
		if rest != "또는 `X`" {
			t.Fatalf("rest = %q, want %q", rest, "또는 `X`")
		}
	})

	t.Run("PrefixMatchButNoWordBoundary", func(t *testing.T) {
		rest, ok := takeKeywordKo("그리고도 `X`", "그리고")
		if ok {
			t.Fatal("expected no match when following byte is not a boundary")
		}
		if rest != "그리고도 `X`" {
			t.Fatalf("rest = %q, want %q", rest, "그리고도 `X`")
		}
	})

	t.Run("MatchWithSpaceBoundary", func(t *testing.T) {
		rest, ok := takeKeywordKo("다음 중 `X`", "다음 중")
		if !ok {
			t.Fatal("expected match for compound keyword")
		}
		if rest != " `X`" {
			t.Fatalf("rest = %q, want %q", rest, " `X`")
		}
	})

	t.Run("MatchWithTabBoundary", func(t *testing.T) {
		rest, ok := takeKeywordKo("그리고\t`X`", "그리고")
		if !ok {
			t.Fatal("expected match with tab boundary")
		}
		if rest != "\t`X`" {
			t.Fatalf("rest = %q, want %q", rest, "\t`X`")
		}
	})

	t.Run("MatchWithEmptyRest", func(t *testing.T) {
		rest, ok := takeKeywordKo("  다음 중", "다음 중")
		if !ok {
			t.Fatal("expected match when keyword consumes the whole string")
		}
		if rest != "" {
			t.Fatalf("rest = %q, want empty", rest)
		}
	})

	t.Run("MatchWithCommaBoundary", func(t *testing.T) {
		// A modifier keyword abutting a separator comma (left by strings.Join)
		// still matches; the leading comma survives in rest for skipParticlesKo.
		rest, ok := takeKeywordKo("하는, 가중치 0.8 의", "하는")
		if !ok {
			t.Fatal("expected match when the following byte is a comma")
		}
		if rest != ", 가중치 0.8 의" {
			t.Fatalf("rest = %q, want %q", rest, ", 가중치 0.8 의")
		}
	})
}
