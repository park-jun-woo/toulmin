//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ScanItems — split section lines into a nested list-item tree plus interstitial prose
package scan

import "strings"

// ScanItems partitions a section's lines into list-item lines and prose lines,
// then builds a nested Item tree (by indentation) from the item lines. A line is
// an item when its trimmed text starts with `- ` or `N. `; everything else
// (blank lines, non-list text) is returned as prose so no bytes are lost.
func ScanItems(lines []Line) (items []Item, prose []Line) {
	var itemLines []Line
	for _, ln := range lines {
		if _, _, _, ok := parseMarker(strings.TrimSpace(ln.Text)); ok {
			itemLines = append(itemLines, ln)
			continue
		}
		prose = append(prose, ln)
	}
	return buildItemTree(itemLines), prose
}
