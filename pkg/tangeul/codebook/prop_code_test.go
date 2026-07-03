//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestPropCode — each PropCode constant equals the §2.3 spec table code
package codebook

import "testing"

func TestPropCode(t *testing.T) {
	cases := []struct {
		got  PropCode
		want uint16
	}{
		{PropAttacks, 0x0101},
		{PropDo, 0x0102},
		{PropUndo, 0x0103},
		{PropRun, 0x0104},
		{PropChecks, 0x0105},
		{PropRequires, 0x0106},
		{PropMeans, 0x0107},
		{PropHas, 0x0108},
		{PropRuns, 0x0109},
		{PropEvaluates, 0x010A},
		{PropTriggers, 0x010B},
		{PropUntil, 0x010C},
		{PropSameAs, 0x010D},
		{PropImportedFrom, 0x010E},
		{PropLabel, 0x010F},
		{PropProvenance, 0x0110},
		{PropSupersedes, 0x0111},
	}
	for _, c := range cases {
		if uint16(c.got) != c.want {
			t.Errorf("PropCode: got 0x%04X, want 0x%04X", uint16(c.got), c.want)
		}
	}
}
