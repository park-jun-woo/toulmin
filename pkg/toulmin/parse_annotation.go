//ff:func feature=annotation type=parser control=iteration dimension=1
//ff:what ParseAnnotation — parses //tm: comment lines into RuleMeta
package toulmin

import "strings"

// ParseAnnotation parses //tm: prefixed comment lines into a RuleMeta.
// The returned RuleMeta has Fn=nil; the caller must set it before registering.
func ParseAnnotation(lines []string) RuleMeta {
	meta := RuleMeta{Qualifier: 1.0}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//tm:backing ") {
			meta.Backing = strings.Trim(trimmed[len("//tm:backing "):], "\"")
		}
	}
	return meta
}
