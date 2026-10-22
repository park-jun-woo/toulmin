//ff:func feature=tangl type=parser control=sequence
//ff:what sectionName — extract the verbatim name from a `## tangl:<Name>` heading
package scan

import "strings"

// sectionName extracts the raw name from a `## tangl:<Name>` heading line,
// matching the parser's tanglSectionName rule. It reports ok=false for any
// other heading or non-heading line. The name is kept verbatim; locale
// canonicalization is the codec's job.
func sectionName(text string) (string, bool) {
	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "## tangl:") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(trimmed, "## tangl:"))
	if name == "" {
		return "", false
	}
	return name, true
}
