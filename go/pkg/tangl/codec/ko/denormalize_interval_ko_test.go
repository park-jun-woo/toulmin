//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestDenormalizeIntervalKo — denormalizeIntervalKo inverts s/m/h units and passes clock times through
package ko

import "testing"

// TestDenormalizeIntervalKo checks unit denormalization (s→초, m→분, h→시간) and a
// clock time passing through verbatim, round-tripping normalizeIntervalKo.
func TestDenormalizeIntervalKo(t *testing.T) {
	cases := []struct{ iv, want string }{
		{"30s", "30초"},
		{"5m", "5분"},
		{"1h", "1시간"},
		{"7:00", "7:00"},
	}
	for _, c := range cases {
		if got := denormalizeIntervalKo(c.iv); got != c.want {
			t.Errorf("denormalizeIntervalKo(%q) = %q, want %q", c.iv, got, c.want)
		}
	}
	// The unit forms round-trip back through the encoder's normalizer.
	for _, c := range cases[:3] {
		if back, ok := normalizeIntervalKo(c.want); !ok || back != c.iv {
			t.Errorf("normalizeIntervalKo(%q) = %q,%v, want %q", c.want, back, ok, c.iv)
		}
	}
}
