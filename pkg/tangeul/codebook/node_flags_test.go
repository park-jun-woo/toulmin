//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestNodeFlags — each Node flag occupies its §5.1 bit (F_FUNC b0 … F_WITH b6)
package codebook

import "testing"

func TestNodeFlags(t *testing.T) {
	cases := []struct {
		got NodeFlags
		bit uint
	}{
		{NodeFlagFunc, 0},
		{NodeFlagQual, 1},
		{NodeFlagCheck, 2},
		{NodeFlagClaim, 3},
		{NodeFlagVerdict, 4},
		{NodeFlagStrike, 5},
		{NodeFlagWith, 6},
	}
	for _, c := range cases {
		if uint16(c.got) != uint16(1)<<c.bit {
			t.Errorf("NodeFlag: got 0x%04X, want bit %d (0x%04X)", uint16(c.got), c.bit, uint16(1)<<c.bit)
		}
	}
}
