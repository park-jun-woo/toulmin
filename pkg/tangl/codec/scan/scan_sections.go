//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ScanSections — split lines into `## tangl:` sections, closing scope on any heading
package scan

// ScanSections walks lines and returns each `## tangl:<Name>` section along with
// the preamble lines before the first such section. A heading of any level
// closes a section's scope (matching the parser's findSectionEnd), so lines
// outside sections also arise between and after sections; the consumer fills
// that residual prose by sweeping the full line table. Section.Name is verbatim.
func ScanSections(lines []Line) (preamble []Line, secs []Section) {
	i := 0
	for i < len(lines) {
		name, ok := sectionName(lines[i].Text)
		if !ok && len(secs) == 0 {
			preamble = append(preamble, lines[i])
			i++
			continue
		}
		if !ok {
			i++
			continue
		}
		start := i + 1
		end := sectionEnd(lines, start)
		secs = append(secs, Section{
			Name:       name,
			HeaderSpan: lines[i].Span,
			Lines:      lines[start:end],
		})
		i = end
	}
	return preamble, secs
}
