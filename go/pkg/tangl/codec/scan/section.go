//ff:type feature=tangl type=parser
//ff:what Section — one `## tangl:<Name>` section: verbatim name, header span, and content lines
package scan

// Section is one `## tangl:<Name>` block found by ScanSections. Name is the
// raw heading name (locale canonicalization is the codec's job). HeaderSpan
// covers the heading line's raw bytes, and Lines are the content lines up to
// (but excluding) the next heading of any level or EOF.
type Section struct {
	Name       string
	HeaderSpan Span
	Lines      []Line
}
