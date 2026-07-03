//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeBacktick — no-prefix, unterminated, and success (including spaced Hangul names)
package ko

import "testing"

func TestTakeBacktick(t *testing.T) {
	t.Run("NoPrefix", func(t *testing.T) {
		name, rest, ok := takeBacktick("  참조한다")
		if ok {
			t.Fatal("expected ok=false for text with no leading backtick")
		}
		if name != "" {
			t.Errorf("name = %q, want empty", name)
		}
		if rest != "참조한다" {
			t.Errorf("rest = %q, want %q", rest, "참조한다")
		}
	})

	t.Run("Unterminated", func(t *testing.T) {
		name, rest, ok := takeBacktick("`미완")
		if ok {
			t.Fatal("expected ok=false for an unterminated backtick")
		}
		if name != "" {
			t.Errorf("name = %q, want empty", name)
		}
		if rest != "`미완" {
			t.Errorf("rest = %q, want %q", rest, "`미완")
		}
	})

	t.Run("SuccessSpacedHangulName", func(t *testing.T) {
		name, rest, ok := takeBacktick("  `IP 차단` 이면")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if name != "IP 차단" {
			t.Errorf("name = %q, want %q", name, "IP 차단")
		}
		if rest != " 이면" {
			t.Errorf("rest = %q, want %q", rest, " 이면")
		}
	})
}
