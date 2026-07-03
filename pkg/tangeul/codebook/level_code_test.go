//ff:func feature=tangl type=model control=iteration dimension=1
//ff:what TestLevelCode — each LevelCode constant holds its §5.5 value within the 2-bit range
package codebook

import "testing"

func TestLevelCode(t *testing.T) {
	cases := []struct {
		got  LevelCode
		want uint16
	}{
		{LevelUnset, 0},
		{LevelFail, 1},
		{LevelReview, 2},
	}
	for _, c := range cases {
		if uint16(c.got) != c.want {
			t.Errorf("LevelCode: got %d, want %d", uint16(c.got), c.want)
		}
		if uint16(c.got) > 3 {
			t.Errorf("LevelCode %d exceeds the 2-bit range", uint16(c.got))
		}
	}
}
