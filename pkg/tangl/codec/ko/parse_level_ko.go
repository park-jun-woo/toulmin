//ff:func feature=tangl type=parser control=sequence
//ff:what parseLevelKo — parse the '수준' node clause value (실패→Fail, 검토→Review)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseLevelKo parses the crisp reporting level following 수준: '실패' → LevelFail,
// '검토' → LevelReview. It returns the level and the unconsumed remainder (empty when
// the clause is the last, since 수준 is rendered last). An unrecognized value is an
// error. It mirrors the level branch of pkg/tangl/parser.applyNodeClause.
func parseLevelKo(s string, line int, path string) (ast.Level, string, error) {
	if rest, ok := takeKeywordKo(s, "실패"); ok {
		return ast.LevelFail, rest, nil
	}
	if rest, ok := takeKeywordKo(s, "검토"); ok {
		return ast.LevelReview, rest, nil
	}
	return ast.LevelUnset, "", errAtKo(path, line, "expected 실패 or 검토 after '수준', got %q", s)
}
