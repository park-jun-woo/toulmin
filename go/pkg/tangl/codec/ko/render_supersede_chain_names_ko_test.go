//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersedeChainNamesKo — chain names are backtick-quoted and joined by " > "
package ko

import "testing"

func TestRenderSupersedeChainNamesKo(t *testing.T) {
	if got := renderSupersedeChainNamesKo(nil); got != "" {
		t.Fatalf("empty = %q", got)
	}
	if got, want := renderSupersedeChainNamesKo([]string{"가", "나", "다"}), "`가` > `나` > `다`"; got != want {
		t.Fatalf("chain = %q, want %q", got, want)
	}
}
