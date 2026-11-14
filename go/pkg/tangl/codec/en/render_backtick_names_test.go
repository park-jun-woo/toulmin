//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderBacktickNames — names are backtick-quoted and joined by " and "; none renders empty
package en

import "testing"

func TestRenderBacktickNames(t *testing.T) {
	if got := renderBacktickNames(nil); got != "" {
		t.Fatalf("no names = %q, want empty", got)
	}
	if got := renderBacktickNames([]string{"a"}); got != "`a`" {
		t.Fatalf("one name = %q", got)
	}
	if got, want := renderBacktickNames([]string{"a", "b c", "d"}), "`a` and `b c` and `d`"; got != want {
		t.Fatalf("three names = %q, want %q", got, want)
	}
}
