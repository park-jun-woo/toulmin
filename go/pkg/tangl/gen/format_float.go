//ff:func feature=tangl type=codegen control=sequence
//ff:what formatFloat — renders a float64 as the shortest round-tripping Go literal
package gen

import "strconv"

// formatFloat renders v as the shortest Go float literal that round-trips
// it, used for Qualifier() weights and certainty-gate thresholds.
func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}
