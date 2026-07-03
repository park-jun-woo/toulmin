//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what internLabel — return the index of s in the label table, appending it if new
package packet

// internLabel returns the index of s within *labels, appending s and returning
// its new index when absent. Interning deduplicates repeated field/value strings
// in an EXPR label table so identical spellings share one Text entry.
func internLabel(labels *[]string, s string) uint16 {
	for i, l := range *labels {
		if l == s {
			return uint16(i)
		}
	}
	*labels = append(*labels, s)
	return uint16(len(*labels) - 1)
}
