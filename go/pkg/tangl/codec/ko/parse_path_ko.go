//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parsePathKo — parse a 의-chained backtick data path (`a` 의 `b` 의 `c`)
package ko

// parsePathKo parses a 의-chained backtick data path: one or more backtick names
// joined by the 의 connective (`a` 의 `b` 의 `c`), returning the step names, the
// remaining text, and whether a well-formed path was found. It is the KO mirror of
// pkg/tangl/parser.parsePath (whose separator is a dot). The 의 lookahead backtracks:
// when 의 is not followed by a backtick the connective is left unconsumed (rest still
// begins at the 의), so `점수화` 의 결과 keeps its 의 for the caller — the 의 there is
// the value-call marker's ignored particle, not a path link.
func parsePathKo(s string) (steps []string, rest string, ok bool) {
	first, r, ok := takeBacktick(s)
	if !ok {
		return nil, s, false
	}
	steps = []string{first}
	rest = r
	for {
		r2, ok2 := takeKeywordKo(rest, "의")
		if !ok2 {
			return steps, rest, true
		}
		next, r3, ok3 := takeBacktick(r2)
		if !ok3 {
			return steps, rest, true
		}
		steps = append(steps, next)
		rest = r3
	}
}
