//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseTypeRef — parse an "as Type" clause's type token
package parser

// parseTypeRef parses the TYPE_REF following "as": either a bare built-in
// type keyword (e.g. Currency, Integer) or a backtick-quoted, optionally
// qualified reference (e.g. `credit`.`Threshold`), stringified as "alias.Name".
func parseTypeRef(s string) (string, string, bool) {
	s = trimStart(s)
	if len(s) > 0 && s[0] == '`' {
		ref, rest, ok := parseRef(s)
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
