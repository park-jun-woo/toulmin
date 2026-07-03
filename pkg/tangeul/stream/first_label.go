//ff:func feature=tangl type=util control=sequence
//ff:what firstLabel — the primary (first) label string recorded for a node TID, or empty
package stream

// firstLabel returns the first label recorded for tid — a node's name — or the
// empty string when it has none.
func firstLabel(labels map[uint16][]string, tid uint16) string {
	if ls := labels[tid]; len(ls) > 0 {
		return ls[0]
	}
	return ""
}
