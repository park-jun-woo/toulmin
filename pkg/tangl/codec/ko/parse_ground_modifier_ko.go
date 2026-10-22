//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseGroundModifierKo — parse a Korean grounded-on adnominal name list
package ko

import "strings"

func parseGroundModifierKo(s string) ([]string, string, bool) {
	var names []string
	rest := s
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, s, false
		}
		if len(names) == 0 && strings.HasPrefix(tail, ".") {
			return nil, s, false
		}
		names = append(names, name)
		rest = skipParticlesKo(tail)
		if next, ok := takeKeywordKo(rest, "그리고"); ok {
			rest = next
			continue
		}
		break
	}
	if next, ok := takeKeywordKo(rest, "근거로"); ok {
		if tail, ok := takeKeywordKo(next, "하는"); ok {
			return names, tail, true
		}
	}
	return nil, s, false
}
