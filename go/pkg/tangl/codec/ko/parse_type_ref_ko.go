//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseTypeRefKo — parse a '타입' clause's type token (built-in or Ref)
package ko

import "strings"

// parseTypeRefKo parses the TYPE_REF following '타입': either a bare built-in type
// keyword (e.g. Currency, Integer — never translated) or a backtick-quoted,
// optionally qualified reference (e.g. `센서`.`Temperature`), stringified as
// "alias.Name". It mirrors pkg/tangl/parser.parseTypeRef.
func parseTypeRefKo(s string) (string, string, bool) {
	s = strings.TrimLeft(s, " \t")
	if strings.HasPrefix(s, "`") {
		ref, rest, ok := parseRefKo(s)
		if !ok {
			return "", s, false
		}
		if ref.Alias != "" {
			return ref.Alias + "." + ref.Name, rest, true
		}
		return ref.Name, rest, true
	}
	i := 0
	for i < len(s) && s[i] != ' ' && s[i] != '\t' {
		i++
	}
	if i == 0 {
		return "", s, false
	}
	return s[:i], s[i:], true
}
