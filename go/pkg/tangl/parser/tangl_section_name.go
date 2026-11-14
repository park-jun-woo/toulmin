//ff:func feature=tangl type=parser control=sequence
//ff:what tanglSectionName — extract the section name from a `## tangl:X` heading
package parser

import "strings"

// tanglSectionName extracts the section name from a `## tangl:<Name>` heading
// line. It reports ok=false for any other heading or non-heading line.
func tanglSectionName(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "## tangl:") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(trimmed, "## tangl:"))
	if name == "" {
		return "", false
	}
	return name, true
}
