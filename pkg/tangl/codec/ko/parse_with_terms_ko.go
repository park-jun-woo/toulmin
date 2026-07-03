//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseWithTermsKo — parse a '기준 `a` 그리고 `b`' term list into Definitions names
package ko

import "strings"

// parseWithTermsKo parses the backtick term list following 기준: one or more
// backtick-quoted Definitions term names chained by 그리고, returning the names and
// the unconsumed remainder (e.g. a following ', 가중 ...' clause). It mirrors
// pkg/tangl/parser.parseWithClause.
func parseWithTermsKo(s string, line int, path string) ([]string, string, error) {
	var terms []string
	rest := s
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, "", errAtKo(path, line, "expected backtick-quoted term after '기준'/'그리고'")
		}
		terms = append(terms, name)
		rest = strings.TrimLeft(tail, " \t")
		if r, ok := takeKeywordKo(rest, "그리고"); ok {
			rest = r
			continue
		}
		break
	}
	return terms, rest, nil
}
