//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseSupersedeChainNames — parse "`A` > `B` > ..." into at least 3 distinct names
package parser

import "strings"

// parseSupersedeChainNames reads backtick-quoted names separated by '>' until
// the text is exhausted. A missing name, a duplicate name, a separator other
// than '>', or fewer than 3 names is an error at line.
func parseSupersedeChainNames(s, path string, line int) ([]string, error) {
	var names []string
	seen := map[string]bool{}
	rest := strings.TrimSpace(s)
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, errAt(path, line, "expected backtick-quoted name in supersede chain")
		}
		if seen[name] {
			return nil, errAt(path, line, "duplicate name %q in supersede chain", name)
		}
		seen[name] = true
		names = append(names, name)
		rest = strings.TrimSpace(tail)
		if rest == "" {
			break
		}
		if !strings.HasPrefix(rest, ">") {
			return nil, errAt(path, line, "expected '>' in supersede chain, got %q", rest)
		}
		rest = strings.TrimSpace(rest[1:])
	}
	if len(names) < 3 {
		return nil, errAt(path, line, "supersede chain requires at least 3 names")
	}
	return names, nil
}
