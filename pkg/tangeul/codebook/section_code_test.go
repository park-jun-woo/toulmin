//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestSectionCode — each SectionCode constant equals the §2.4 subcode (1..7)
package codebook

import "testing"

func TestSectionCode(t *testing.T) {
	cases := []struct {
		got  SectionCode
		want uint16
	}{
		{SectionSubject, 1},
		{SectionSee, 2},
		{SectionDefinitions, 3},
		{SectionRules, 4},
		{SectionCases, 5},
		{SectionProvides, 6},
		{SectionInternal, 7},
	}
	for _, c := range cases {
		if uint16(c.got) != c.want {
			t.Errorf("SectionCode: got %d, want %d", uint16(c.got), c.want)
		}
	}
}
