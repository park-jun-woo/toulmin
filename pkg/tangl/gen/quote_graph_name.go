//ff:func feature=tangl type=codegen control=sequence
//ff:what quoteGraphName — renders the "<subject>.<case>" graph name literal
package gen

import "strconv"

// quoteGraphName renders the "<subject>.<case>" graph name
// toulmin.NewGraph expects, as a quoted Go string literal.
func quoteGraphName(subject, caseName string) string {
	return strconv.Quote(subject + "." + caseName)
}
