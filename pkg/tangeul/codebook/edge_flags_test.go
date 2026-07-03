//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestEdgeFlags — each Edge flag occupies its §6.1 bit (F_ONCE b0 … F_ORDERED b3)
package codebook

import "testing"

func TestEdgeFlags(t *testing.T) {
	cases := []struct {
		got EdgeFlags
		bit uint
	}{
		{EdgeFlagOnce, 0},
		{EdgeFlagGate, 1},
		{EdgeFlagStrike, 2},
		{EdgeFlagOrdered, 3},
	}
	for _, c := range cases {
		if uint16(c.got) != uint16(1)<<c.bit {
			t.Errorf("EdgeFlag: got 0x%04X, want bit %d (0x%04X)", uint16(c.got), c.bit, uint16(1)<<c.bit)
		}
	}
}
