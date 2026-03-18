//ff:func feature=annotation type=parser control=iteration dimension=1
//ff:what ParseAnnotation — parses //rule: comment lines into RuleMeta
package toulmin

import "strings"

// ParseAnnotation parses //rule: prefixed comment lines into a RuleMeta.
// The returned RuleMeta has Fn=nil; the caller must set it before registering.
func ParseAnnotation(lines []string) RuleMeta {
	var meta RuleMeta
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//rule:backing ") {
			meta.Backing = strings.Trim(trimmed[len("//rule:backing "):], "\"")
			continue
		}
		if strings.HasPrefix(trimmed, "//rule:what ") {
			meta.What = trimmed[len("//rule:what "):]
			continue
		}
		var rest string
		if strings.HasPrefix(trimmed, "//rule:warrant ") {
			rest = trimmed[len("//rule:warrant "):]
		} else if strings.HasPrefix(trimmed, "//rule:defeater ") {
			rest = trimmed[len("//rule:defeater "):]
			meta.Strength = Defeater
		} else {
			continue
		}
		parseParams(&meta, rest)
	}
	return meta
}
