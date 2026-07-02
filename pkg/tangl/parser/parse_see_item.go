//ff:func feature=tangl type=parser control=sequence
//ff:what parseSeeItem — parse a single `see `alias` from `package`` entry
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseSeeItem parses "see `<alias>` from `<package_path>`".
func parseSeeItem(it item, path string) (ast.See, error) {
	rest, ok := takeKeyword(it.Text, "see")
	if !ok {
		return ast.See{}, errAt(path, it.Line, "expected 'see `alias` from `package`', got %q", it.Text)
	}
	alias, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.See{}, errAt(path, it.Line, "expected backtick-quoted alias after 'see'")
	}
	rest, ok = takeKeyword(rest, "from")
	if !ok {
		return ast.See{}, errAt(path, it.Line, "expected 'from' after alias")
	}
	pkg, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.See{}, errAt(path, it.Line, "expected backtick-quoted package path after 'from'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.See{}, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return ast.See{Alias: alias, Package: pkg, Line: it.Line}, nil
}
