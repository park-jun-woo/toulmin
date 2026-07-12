//ff:func feature=tangl type=parser control=sequence
//ff:what TestParsePathKo — single/multi 의-chain steps and the 의-before-non-backtick backtrack
package ko

import (
	"testing"
)

func TestParsePathKo(t *testing.T) {
	t.Run("SingleStep", func(t *testing.T) {
		steps, rest, ok := parsePathKo("`작성자` 가 비어 있으면")
		if !ok || len(steps) != 1 || steps[0] != "작성자" {
			t.Fatalf("got steps=%v ok=%v", steps, ok)
		}
		if rest != " 가 비어 있으면" {
			t.Fatalf("rest = %q", rest)
		}
	})

	t.Run("MultiStep", func(t *testing.T) {
		steps, rest, ok := parsePathKo("`작성자` 의 `신뢰점수` 가 0.7 이상이면")
		if !ok || len(steps) != 2 || steps[0] != "작성자" || steps[1] != "신뢰점수" {
			t.Fatalf("got steps=%v ok=%v", steps, ok)
		}
		if rest != " 가 0.7 이상이면" {
			t.Fatalf("rest = %q", rest)
		}
	})

	t.Run("BacktrackOnNonBacktickAfter의", func(t *testing.T) {
		// `점수화` 의 결과 — the 의 belongs to the value-call marker, not the path.
		steps, rest, ok := parsePathKo("`점수화` 의 결과 가")
		if !ok || len(steps) != 1 || steps[0] != "점수화" {
			t.Fatalf("got steps=%v ok=%v", steps, ok)
		}
		if rest != " 의 결과 가" {
			t.Fatalf("rest = %q, want the un-consumed 의", rest)
		}
	})

	t.Run("NotABacktick", func(t *testing.T) {
		if _, _, ok := parsePathKo("기준 `x`"); ok {
			t.Fatal("expected ok=false when the text does not start with a backtick")
		}
	})
}
