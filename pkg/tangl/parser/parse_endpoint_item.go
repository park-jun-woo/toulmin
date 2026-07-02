//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseEndpointItem — parse one "provides `name`" endpoint block
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseEndpointItem parses "provides `<name>`" and its nested required
// fields and run/check declarations.
func parseEndpointItem(it item, path string) (ast.Endpoint, error) {
	rest, ok := takeKeyword(it.Text, "provides")
	if !ok {
		return ast.Endpoint{}, errAt(path, it.Line, "expected 'provides `name`', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Endpoint{}, errAt(path, it.Line, "expected backtick-quoted endpoint name")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Endpoint{}, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	ep := ast.Endpoint{Name: name, Line: it.Line}
	for _, child := range it.Children {
		if req, ok, err := parseRequireItem(child, path); err != nil {
			return ast.Endpoint{}, err
		} else if ok {
			ep.Requires = append(ep.Requires, req)
			continue
		}
		caseName, kind, ok, err := parseRunCheckItem(child, path)
		if err != nil {
			return ast.Endpoint{}, err
		}
		if !ok {
			return ast.Endpoint{}, errAt(path, child.Line, "unrecognized provides statement: %q", child.Text)
		}
		if kind == "run" {
			ep.Runs = append(ep.Runs, caseName)
		} else {
			ep.Checks = append(ep.Checks, caseName)
		}
	}
	return ep, nil
}
