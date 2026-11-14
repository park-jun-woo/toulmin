//ff:func feature=tangl type=parser control=sequence
//ff:what TestCompareOpsKo — Op strings are the 8 canonical operators and phrases are ordered longest-first
package ko

import (
	"testing"
	"unicode/utf8"
)

func TestCompareOpsKo(t *testing.T) {
	canonical := map[string]bool{
		"is not empty": true, "is empty": true, "is greater than": true,
		"is less than": true, "is at most": true, "is at least": true,
		"is in": true, "equals": true,
	}
	t.Run("OpsAreCanonical", func(t *testing.T) {
		for _, o := range compareOpsKo {
			if !canonical[o.Op] {
				t.Errorf("Op %q is not one of the eight canonical operators", o.Op)
			}
		}
	})

	t.Run("LongestFirst", func(t *testing.T) {
		prev := utf8.RuneCountInString(compareOpsKo[0].Ko)
		for _, o := range compareOpsKo[1:] {
			n := utf8.RuneCountInString(o.Ko)
			if n > prev {
				t.Errorf("phrase %q (%d runes) is longer than its predecessor (%d)", o.Ko, n, prev)
			}
			prev = n
		}
	})
}
