//ff:func feature=tangl type=parser control=selection
//ff:what parseReportKo — parse the '(실패를|검토를) 보고하는' adnominal level modifier
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseReportKo parses the optional 보고 modifier: '실패를 보고하는' → LevelFail,
// '검토를 보고하는' → LevelReview. The fused 를 is an ignored particle attached to
// the keyword (정서법), so the level word's leading token is 실패를/실패 (or 검토를/검토)
// — both allomorphs select the same level, and skipParticlesKo peels a detached 를,
// so '실패 를 보고하는' parses identically. ok is false (s unconsumed) when s does not
// open with a level word; a level word without a following '보고하는' is an error.
func parseReportKo(s string, line int, path string) (ast.Level, string, bool, error) {
	head := strings.TrimLeft(s, " \t")
	tok := head
	if i := strings.IndexAny(head, " \t"); i >= 0 {
		tok = head[:i]
	}
	var lvl ast.Level
	switch tok {
	case "실패를", "실패":
		lvl = ast.LevelFail
	case "검토를", "검토":
		lvl = ast.LevelReview
	default:
		return ast.LevelUnset, s, false, nil
	}
	rest := head[len(tok):]
	r, ok := takeKeywordKo(skipParticlesKo(rest), "보고하는")
	if !ok {
		return ast.LevelUnset, "", false, errAtKo(path, line, "expected '보고하는' after level word, got %q", rest)
	}
	return lvl, r, true, nil
}
