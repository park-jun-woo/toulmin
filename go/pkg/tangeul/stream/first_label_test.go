//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFirstLabel — a populated label list, a missing TID and a present-but-empty list
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestFirstLabel covers both arms of firstLabel: a TID with one or more
// labels returns the first, and a TID missing from the map or mapped to an
// empty slice returns "".
func TestFirstLabel(t *testing.T) {
	labels := map[toulmin.TID][]string{
		1: {"balance sufficient", "type"},
		2: {},
	}
	cases := []struct {
		name string
		tid  toulmin.TID
		want string
	}{
		{"present", 1, "balance sufficient"},
		{"emptySlice", 2, ""},
		{"missing", 3, ""},
	}
	for _, c := range cases {
		if got := firstLabel(labels, c.tid); got != c.want {
			t.Errorf("%s: firstLabel = %q, want %q", c.name, got, c.want)
		}
	}
}
