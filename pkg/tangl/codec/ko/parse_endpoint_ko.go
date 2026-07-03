//ff:func feature=tangl type=parser control=sequence
//ff:what parseEndpointKo — parse an endpoint head '`X` 를 한다면' into an ast.Endpoint (name + line)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseEndpointKo parses the tangl:업무 endpoint head "`X` 를 한다면", returning an
// ast.Endpoint with its Name and Line but no body — the caller walks the item's
// children (필수다 requires, 진행한다 runs, 판정한다 checks) so it can assign body slot
// coordinates. 한다면 is one keyword (its 면 is not a standalone particle).
func parseEndpointKo(it scan.Item, path string) (ast.Endpoint, error) {
	body, ok := takeSuffixKo(it.Text, "한다면")
	if !ok {
		return ast.Endpoint{}, errAtKo(path, it.No, "expected '`X` 를 한다면', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Endpoint{}, errAtKo(path, it.No, "expected backtick-quoted endpoint name")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Endpoint{}, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return ast.Endpoint{Name: name, Line: it.No}, nil
}
