//ff:type feature=tangl type=parser
//ff:what section — the raw content lines of one `## tangl:X` section
package parser

// section holds the raw content lines belonging to one `## tangl:<Name>`
// heading, up to (but excluding) the next heading of any level or EOF.
type section struct {
	Name       string
	HeaderLine int      // 1-based line number of the "## tangl:Name" heading
	Lines      []string // content lines (may be empty)
	LineOffset int      // 1-based line number of Lines[0]
}
