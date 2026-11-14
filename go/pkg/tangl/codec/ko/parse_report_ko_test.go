//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseReportKo — 실패를/검토를→level, 붙임·띄움 병기 동등성, 미개시 후퇴, 보고하는 누락 에러
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseReportKo(t *testing.T) {
	t.Run("FailFused", func(t *testing.T) {
		lvl, rest, ok, err := parseReportKo("실패를 보고하는 `정책`.`f` 쓰는", 1, "t.md")
		if err != nil || !ok || lvl != ast.LevelFail {
			t.Fatalf("lvl=%v ok=%v err=%v", lvl, ok, err)
		}
		if rest != " `정책`.`f` 쓰는" {
			t.Fatalf("rest=%q", rest)
		}
	})

	t.Run("ReviewFused", func(t *testing.T) {
		lvl, _, ok, err := parseReportKo("검토를 보고하는", 1, "t.md")
		if err != nil || !ok || lvl != ast.LevelReview {
			t.Fatalf("lvl=%v ok=%v err=%v", lvl, ok, err)
		}
	})

	// 붙여 쓴 '실패를' 과 띄어 쓴 '실패 를' 은 동일한 AST(이형태 병기 축).
	t.Run("FusedAndDetachedEquivalent", func(t *testing.T) {
		a, ar, aok, aerr := parseReportKo("실패를 보고하는", 1, "t.md")
		b, br, bok, berr := parseReportKo("실패 를 보고하는", 1, "t.md")
		if aerr != nil || berr != nil || !aok || !bok {
			t.Fatalf("fused=(%v,%v) detached=(%v,%v)", aok, aerr, bok, berr)
		}
		if a != b || ar != br {
			t.Fatalf("mismatch: fused=(%v,%q) detached=(%v,%q)", a, ar, b, br)
		}
		if a != ast.LevelFail {
			t.Fatalf("lvl=%v, want Fail", a)
		}
	})

	t.Run("NoLevelWordBacktracks", func(t *testing.T) {
		in := "`정책`.`f` 쓰는"
		lvl, rest, ok, err := parseReportKo(in, 1, "t.md")
		if err != nil || ok || lvl != ast.LevelUnset || rest != in {
			t.Fatalf("expected clean backtrack, got lvl=%v rest=%q ok=%v err=%v", lvl, rest, ok, err)
		}
	})

	t.Run("MissingBogohaneunErrors", func(t *testing.T) {
		if _, _, ok, err := parseReportKo("실패를 쓰는", 1, "t.md"); err == nil || ok {
			t.Fatalf("expected error for a level word without '보고하는' (ok=%v err=%v)", ok, err)
		}
	})
}
