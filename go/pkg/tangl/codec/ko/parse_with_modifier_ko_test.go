//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseWithModifierKo — 기준절 확정, pkg.fn/판정 후퇴, 그리고-연쇄, 밀착 쉼표
package ko

import "testing"

func TestParseWithModifierKo(t *testing.T) {
	// (a) 기준절 확정 — single term followed by '을 기준으로 하는'.
	t.Run("SingleTermWith", func(t *testing.T) {
		terms, rest, ok := parseWithModifierKo("`서빙 온도` 를 기준으로 하는 `센서`.`temperatureLow` 쓰는")
		if !ok {
			t.Fatal("expected the 기준 modifier to be recognized")
		}
		if len(terms) != 1 || terms[0] != "서빙 온도" {
			t.Fatalf("terms=%v", terms)
		}
		if rest != " `센서`.`temperatureLow` 쓰는" {
			t.Fatalf("rest=%q", rest)
		}
	})

	// (a) 그리고-chained terms.
	t.Run("ChainedTermsWith", func(t *testing.T) {
		terms, rest, ok := parseWithModifierKo("`신용 기준` 그리고 `연체 이력` 을 기준으로 하는 `심사`.`Evaluate` 쓰는")
		if !ok {
			t.Fatal("expected the chained 기준 modifier to be recognized")
		}
		if len(terms) != 2 || terms[0] != "신용 기준" || terms[1] != "연체 이력" {
			t.Fatalf("terms=%v", terms)
		}
		if rest != " `심사`.`Evaluate` 쓰는" {
			t.Fatalf("rest=%q", rest)
		}
	})

	// (b) pkg.fn ref — '.' right after the first backtick short-circuits to false,
	// leaving s fully unconsumed for the caller's '쓰는' handling.
	t.Run("PkgFnRefBacktracks", func(t *testing.T) {
		in := "`센서`.`noCupAvailable` 쓰는"
		terms, rest, ok := parseWithModifierKo(in)
		if ok || terms != nil || rest != in {
			t.Fatalf("expected clean backtrack, got terms=%v rest=%q ok=%v", terms, rest, ok)
		}
	})

	// (c) '`Y` 판정 쓰는' — backtick run not followed by 기준으로, so backtrack.
	t.Run("CheckingBacktracks", func(t *testing.T) {
		in := "`하위 판정` 판정 쓰는"
		terms, rest, ok := parseWithModifierKo(in)
		if ok || terms != nil || rest != in {
			t.Fatalf("expected clean backtrack, got terms=%v rest=%q ok=%v", terms, rest, ok)
		}
	})

	// Not a backtick at all (e.g. leads with 가중치).
	t.Run("NoBacktickBacktracks", func(t *testing.T) {
		in := "가중치 0.8 의"
		terms, rest, ok := parseWithModifierKo(in)
		if ok || terms != nil || rest != in {
			t.Fatalf("expected clean backtrack, got terms=%v rest=%q ok=%v", terms, rest, ok)
		}
	})

	// Comma abutting 하는 (separator left by strings.Join) is tolerated — the
	// leading comma stays in rest for the caller's skipParticlesKo.
	t.Run("AbuttingSeparatorComma", func(t *testing.T) {
		terms, rest, ok := parseWithModifierKo("`차단 목록` 을 기준으로 하는, 가중치 0.8 의")
		if !ok {
			t.Fatal("expected the 기준 modifier before a separator comma")
		}
		if len(terms) != 1 || terms[0] != "차단 목록" {
			t.Fatalf("terms=%v", terms)
		}
		if rest != ", 가중치 0.8 의" {
			t.Fatalf("rest=%q", rest)
		}
	})
}
