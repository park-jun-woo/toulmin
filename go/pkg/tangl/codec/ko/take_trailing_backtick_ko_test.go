//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeTrailingBacktick — last backtick name split, no trailing backtick, single unterminated backtick
package ko

import "testing"

func TestTakeTrailingBacktick(t *testing.T) {
	t.Run("SplitsLastName", func(t *testing.T) {
		name, prefix, ok := takeTrailingBacktick("`신용점수` 가 650 이상이면 `신용 충분`")
		if !ok || name != "신용 충분" || prefix != "`신용점수` 가 650 이상이면 " {
			t.Fatalf("name=%q prefix=%q ok=%v", name, prefix, ok)
		}
	})

	t.Run("NestedHead", func(t *testing.T) {
		name, prefix, ok := takeTrailingBacktick("다음이면 `우대 대상`")
		if !ok || name != "우대 대상" || prefix != "다음이면 " {
			t.Fatalf("name=%q prefix=%q ok=%v", name, prefix, ok)
		}
	})

	t.Run("NoTrailingBacktick", func(t *testing.T) {
		if _, _, ok := takeTrailingBacktick("`x` 이다"); ok {
			t.Fatal("expected ok=false when text does not end with a backtick")
		}
	})

	t.Run("SingleBacktick", func(t *testing.T) {
		if _, _, ok := takeTrailingBacktick("no name here `"); ok {
			t.Fatal("expected ok=false for a lone unterminated backtick")
		}
	})
}
