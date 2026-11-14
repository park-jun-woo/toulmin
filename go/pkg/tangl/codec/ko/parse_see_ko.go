//ff:func feature=tangl type=parser control=sequence
//ff:what parseSeeKo — parse '`p` 는 `path` 에서 참조한다' into an ast.See
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseSeeKo parses the tangl:참조 statement "`alias` 는 `package` 에서 참조한다".
// The sentence-final 참조한다 is consumed first, then the two backtick names are
// read left to right with intervening particles (는/에서) skipped.
func parseSeeKo(it scan.Item, path string) (ast.See, error) {
	body, ok := takeSuffixKo(it.Text, "참조한다")
	if !ok {
		return ast.See{}, errAtKo(path, it.No, "expected '`p` 는 `path` 에서 참조한다', got %q", it.Text)
	}
	alias, rest, ok := takeBacktick(body)
	if !ok {
		return ast.See{}, errAtKo(path, it.No, "expected backtick-quoted alias")
	}
	rest = skipParticlesKo(rest)
	pkg, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.See{}, errAtKo(path, it.No, "expected backtick-quoted package path")
	}
	rest = skipParticlesKo(rest)
	if strings.TrimSpace(rest) != "" {
		return ast.See{}, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return ast.See{Alias: alias, Package: pkg, Line: it.No}, nil
}
