//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestVFixFromFloat — spec golden bit patterns (±1.0, 0.0, 0.8) plus out-of-range rejection
package word

import "testing"

func TestVFixFromFloat(t *testing.T) {
	// Exact spec example values (§3.1): stored int16 bit pattern must match,
	// and a nil error slot means the value is expected to be rejected.
	cases := []struct {
		in   float64
		want uint16
		ok   bool
	}{
		{1.0, 0x2710, true},
		{0.0, 0x0000, true},
		{-1.0, 0xD8F0, true},
		{0.8, 0x1F40, true},
		{1.0001, 0, false},
		{-2.0, 0, false},
	}
	for _, c := range cases {
		v, err := VFixFromFloat(c.in)
		if c.ok && (err != nil || uint16(v) != c.want) {
			t.Fatalf("VFixFromFloat(%v): got 0x%04X err=%v want 0x%04X", c.in, uint16(v), err, c.want)
		}
		if !c.ok && err == nil {
			t.Fatalf("VFixFromFloat(%v): expected out-of-range error", c.in)
		}
	}
}
