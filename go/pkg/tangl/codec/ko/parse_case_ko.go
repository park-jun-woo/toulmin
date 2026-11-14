//ff:func feature=tangl type=parser control=sequence
//ff:what parseCaseKo — parse a case head '`X` 의 경우에' into an ast.Case (name + line, no body)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseCaseKo parses the tangl:사안 case head "`X` 의 경우에", returning an ast.Case
// with its Name and Line but no body — the caller walks the item's children and
// dispatches each via applyCaseChildKo so it can assign body slot coordinates. The
// sentence-final 경우에 is consumed as one keyword (the 에 is not a standalone
// particle), then the single backtick case name is read.
func parseCaseKo(it scan.Item, path string) (ast.Case, error) {
	body, ok := takeSuffixKo(it.Text, "경우에")
	if !ok {
		return ast.Case{}, errAtKo(path, it.No, "expected '`X` 의 경우에', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Case{}, errAtKo(path, it.No, "expected backtick-quoted case name")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Case{}, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return ast.Case{Name: name, Line: it.No}, nil
}
