//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersedeChainNames — chain names are backtick-quoted and joined by " > " in order
package en

import "testing"

func TestRenderSupersedeChainNames(t *testing.T) {
	if got := renderSupersedeChainNames(nil); got != "" {
		t.Fatalf("empty chain = %q, want empty", got)
	}
	if got, want := renderSupersedeChainNames([]string{"A", "B", "C"}), "`A` > `B` > `C`"; got != want {
		t.Fatalf("chain = %q, want %q", got, want)
	}
}
