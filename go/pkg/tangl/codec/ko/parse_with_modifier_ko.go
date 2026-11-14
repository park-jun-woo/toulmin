//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseWithModifierKo — tentatively parse '`a` 그리고 `b` 을 기준으로 하는', backtracking on miss
package ko

import "strings"

// parseWithModifierKo tentatively parses a leading 기준 adnominal modifier. ok is
// false — with s unconsumed — when the leading backtick run is not followed by
// '기준으로 하는' (i.e. it is a '<fn> 쓰는' or '`Y` 판정 쓰는' clause instead), or when
// s does not start with a backtick at all; a '.' right after the first backtick
// name short-circuits to false (a pkg.fn ref is never a Definitions term list).
// Malformed text also returns false so the caller's later '쓰는' check reports it.
func parseWithModifierKo(s string) ([]string, string, bool) {
	var terms []string
	rest := s
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, s, false
		}
		if len(terms) == 0 && strings.HasPrefix(tail, ".") {
			return nil, s, false
		}
		terms = append(terms, name)
		rest = skipParticlesKo(tail)
		if r, ok := takeKeywordKo(rest, "그리고"); ok {
			rest = r
			continue
		}
		break
	}
	if r, ok := takeKeywordKo(rest, "기준으로"); ok {
		if r2, ok2 := takeKeywordKo(r, "하는"); ok2 {
			return terms, r2, true
		}
	}
	return nil, s, false
}
