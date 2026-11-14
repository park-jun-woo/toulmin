//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRequireKo — bare require, typed require, non-require passthrough, malformed type
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseRequireKo(t *testing.T) {
	t.Run("Bare", func(t *testing.T) {
		req, ok, err := parseRequireKo(scan.Item{Text: "`사용자` 는 필수다", No: 3}, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if req.Field != "사용자" || req.Type != "" || req.Line != 3 {
			t.Fatalf("got %+v", req)
		}
	})

	t.Run("Typed", func(t *testing.T) {
		req, ok, err := parseRequireKo(scan.Item{Text: "`금액` 은 필수다, 타입 Currency", No: 4}, "t.md")
		if err != nil || !ok || req.Type != "Currency" {
			t.Fatalf("got %+v ok=%v err=%v", req, ok, err)
		}
	})

	t.Run("NotRequire", func(t *testing.T) {
		_, ok, err := parseRequireKo(scan.Item{Text: "`인증` 은 일반 규칙이다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("NoBacktick", func(t *testing.T) {
		_, ok, _ := parseRequireKo(scan.Item{Text: "사용자 는 필수다", No: 1}, "t.md")
		if ok {
			t.Fatal("expected ok=false without a backtick field")
		}
	})

	t.Run("MalformedTail", func(t *testing.T) {
		if _, ok, err := parseRequireKo(scan.Item{Text: "`x` 는 필수다, Currency", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected an error for a tail without 타입 (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("NonParticleBeforeKeyword", func(t *testing.T) {
		_, ok, err := parseRequireKo(scan.Item{Text: "`x` 는 강조 필수다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough for non-particle text before 필수다, ok=%v err=%v", ok, err)
		}
	})

	t.Run("MissingTypeAfterKeyword", func(t *testing.T) {
		if _, ok, err := parseRequireKo(scan.Item{Text: "`x` 는 필수다, 타입", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected an error when '타입' has no type after it (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("TrailingTextAfterType", func(t *testing.T) {
		if _, ok, err := parseRequireKo(scan.Item{Text: "`x` 는 필수다, 타입 Currency 여분텍스트", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected an error for trailing text after the type ref (ok=%v err=%v)", ok, err)
		}
	})
}
