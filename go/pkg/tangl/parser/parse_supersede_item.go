//ff:func feature=tangl type=parser control=sequence
//ff:what parseSupersedeItem — parse a "`superseder` supersedes `superseded`" edge
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseSupersedeItem parses edge_decl: NAME "supersedes" NAME. It returns
// ok=false (no error) when the item is not a supersede edge, so the caller
// can try the other case-body forms next.
func parseSupersedeItem(it item, path string) (ast.Supersede, bool, error) {
	superseder, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Supersede{}, false, nil
	}
	rest, ok = takeKeyword(rest, "supersedes")
	if !ok {
		return ast.Supersede{}, false, nil
	}
	superseded, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Supersede{}, true, errAt(path, it.Line, "expected backtick-quoted target after 'supersedes'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Supersede{}, true, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return ast.Supersede{Superseder: superseder, Superseded: superseded, Line: it.Line}, true, nil
}
