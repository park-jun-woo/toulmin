//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderGroundsKo — names join with 그리고 and end in a 받침-selected 을/를 근거로 하는
package ko

import "testing"

func TestRenderGroundsKo(t *testing.T) {
	if got := renderGroundsKo(nil); got != "" {
		t.Fatalf("no grounds = %q, want empty", got)
	}
	if got, want := renderGroundsKo([]string{"시도 기록"}), "`시도 기록` 을 근거로 하는"; got != want {
		t.Fatalf("one ground = %q, want %q", got, want)
	}
	if got, want := renderGroundsKo([]string{"기록", "로그"}), "`기록` 그리고 `로그` 를 근거로 하는"; got != want {
		t.Fatalf("two grounds = %q, want %q", got, want)
	}
}
