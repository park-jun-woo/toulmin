//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseSection — parse all statement lines within a TANGL section
package parser

import "strings"

// parseSection parses lines within a TANGL section until the next heading or end of file.
func parseSection(lines []string, i int, sectionName string, f *File, errs []string) (int, []string) {
	for i < len(lines) {
		trimmedRaw := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmedRaw, "## ") {
			break
		}
		adv, newErrs := parseSectionEntry(lines, i, sectionName, trimmedRaw, f, errs)
		errs = newErrs
		i += adv
	}
	return i, errs
}
