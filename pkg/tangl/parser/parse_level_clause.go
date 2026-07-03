//ff:func feature=tangl type=parser control=sequence
//ff:what parseLevelClause — parse a "level fail" / "level review" node suffix
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseLevelClause parses level_clause: "level" ("fail" | "review"). The
// leading comma is consumed by the caller; s begins at the "level" keyword.
// It returns the parsed Level and the remaining text.
func parseLevelClause(s string, path string, line int) (ast.Level, string, error) {
	rest, ok := takeKeyword(s, "level")
	if !ok {
		return ast.LevelUnset, s, errAt(path, line, "expected 'level' after ','")
	}
	if tail, ok := takeKeyword(rest, "fail"); ok {
		return ast.LevelFail, tail, nil
	}
	if tail, ok := takeKeyword(rest, "review"); ok {
		return ast.LevelReview, tail, nil
	}
	return ast.LevelUnset, s, errAt(path, line, "expected 'fail' or 'review' after 'level': %q", strings.TrimSpace(rest))
}
