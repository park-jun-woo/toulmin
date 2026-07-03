//ff:func feature=tangl type=util control=sequence
//ff:what nodeKey — the in-case node lookup key (case index + node name)
package stream

import "strconv"

// nodeKey builds the tidTable.nodeByCase key for an in-case rule node, scoping
// the node name to its case index. Node names are case-local (TANGL §execution
// edges), so attacks and `when <node>` triggers must resolve within the same
// case; the NUL separator keeps distinct (index, name) pairs from aliasing.
func nodeKey(caseIdx int, name string) string {
	return strconv.Itoa(caseIdx) + "\x00" + name
}
