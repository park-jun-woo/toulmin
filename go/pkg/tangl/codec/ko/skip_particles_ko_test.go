//ff:func feature=tangl type=parser control=sequence
//ff:what TestSkipParticlesKo — consumes leading allomorphs and commas, stops at names/keywords, preserves 이다-ordering
package ko

import "testing"

func TestSkipParticlesKo(t *testing.T) {
	t.Run("ConsumesLeadingParticleKeepsBacktick", func(t *testing.T) {
		got := skipParticlesKo("는 `path` 에서 참조한다")
		if got != "`path` 에서 참조한다" {
			t.Fatalf("got %q, want %q", got, "`path` 에서 참조한다")
		}
	})

	t.Run("BothAllomorphsConsumed", func(t *testing.T) {
		pairs := [][2]string{
			{"은 `X`", "`X`"}, {"는 `X`", "`X`"},
			{"이 `X`", "`X`"}, {"가 `X`", "`X`"},
			{"을 `X`", "`X`"}, {"를 `X`", "`X`"},
			{"와 `X`", "`X`"}, {"과 `X`", "`X`"},
			{"로 `X`", "`X`"}, {"으로 `X`", "`X`"},
			{"면 `X`", "`X`"}, {"이면 `X`", "`X`"},
		}
		for _, p := range pairs {
			if got := skipParticlesKo(p[0]); got != p[1] {
				t.Errorf("skipParticlesKo(%q) = %q, want %q", p[0], got, p[1])
			}
		}
	})

	t.Run("NonParticlesNotConsumed", func(t *testing.T) {
		for _, in := range []string{"까지 `C`", "보다 크면", "중 하나면"} {
			if got := skipParticlesKo(in); got != in {
				t.Errorf("skipParticlesKo(%q) = %q, want unchanged", in, got)
			}
		}
	})

	t.Run("ConsumesCommaThenParticle", func(t *testing.T) {
		got := skipParticlesKo(", 타입 `T`")
		if got != "타입 `T`" {
			t.Fatalf("got %q, want %q", got, "타입 `T`")
		}
	})

	t.Run("ConsumesParticleAbuttingComma", func(t *testing.T) {
		// '의,' (particle fused to a separator comma, left by strings.Join) splits
		// into 의 + ',' so the particle is peeled and the leading comma dropped.
		got := skipParticlesKo("의, 실패를 보고하는")
		if got != "실패를 보고하는" {
			t.Fatalf("got %q, want %q", got, "실패를 보고하는")
		}
	})

	t.Run("IdaOrderingNotSplit", func(t *testing.T) {
		// standalone 이 is dropped, but the 계사 키워드 이다 is not a particle so it
		// survives even if skipParticlesKo were applied to it — the suffix keyword
		// is consumed first by takeSuffixKo, so no 이/이다 collision arises.
		if got := skipParticlesKo("이 `문서`"); got != "`문서`" {
			t.Errorf("skipParticlesKo(%q) = %q, want %q", "이 `문서`", got, "`문서`")
		}
		if got := skipParticlesKo("이다"); got != "이다" {
			t.Errorf("skipParticlesKo(%q) = %q, want unchanged", "이다", got)
		}
	})

	t.Run("StopsAtKeyword", func(t *testing.T) {
		if got := skipParticlesKo("참조한다"); got != "참조한다" {
			t.Fatalf("got %q, want unchanged", got)
		}
	})
}
