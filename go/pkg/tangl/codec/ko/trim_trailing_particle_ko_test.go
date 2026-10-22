//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestTrimTrailingParticleKo — strips trailing 와/과/을, keeps a bare value, empty input
package ko

import "testing"

func TestTrimTrailingParticleKo(t *testing.T) {
	cases := [][2]string{
		{`"VIP" 와`, `"VIP"`},
		{`"VIP" 과`, `"VIP"`},
		{"650 을", "650"},
		{"650 를", "650"},
		{"650", "650"},
		{"", ""},
	}
	for _, c := range cases {
		if got := trimTrailingParticleKo(c[0]); got != c[1] {
			t.Errorf("trimTrailingParticleKo(%q) = %q, want %q", c[0], got, c[1])
		}
	}
}
