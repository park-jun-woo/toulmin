//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestControlCode — each ControlCode constant equals the §2.4 spec table code
package codebook

import "testing"

func TestControlCode(t *testing.T) {
	cases := []struct {
		got  ControlCode
		want uint16
	}{
		{ControlStreamHeader, 0x00},
		{ControlSegmentBoundary, 0x01},
		{ControlCodebookDecl, 0x02},
		{ControlSection, 0x03},
		{ControlCaseBegin, 0x04},
		{ControlCaseEnd, 0x05},
		{ControlProse, 0x06},
	}
	for _, c := range cases {
		if uint16(c.got) != c.want {
			t.Errorf("ControlCode: got 0x%02X, want 0x%02X", uint16(c.got), c.want)
		}
	}
}
