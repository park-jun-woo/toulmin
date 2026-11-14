//ff:func feature=tangl type=util control=sequence
//ff:what firstLabel — the primary (first) label string recorded for a node TID, or empty
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// firstLabel returns the first label recorded for tid — a node's name — or the
// empty string when it has none.
func firstLabel(labels map[toulmin.TID][]string, tid toulmin.TID) string {
	if ls := labels[tid]; len(ls) > 0 {
		return ls[0]
	}
	return ""
}
