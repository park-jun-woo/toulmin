//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseInternalItem — parse one "on <event>" / "every <interval>" block
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseInternalItem parses on_decl / every_decl and their nested run/check
// declarations.
func parseInternalItem(it item, path string) (ast.Internal, error) {
	in, err := parseInternalHead(it, path)
	if err != nil {
		return ast.Internal{}, err
	}
	for _, child := range it.Children {
		caseName, kind, ok, err := parseRunCheckItem(child, path)
		if err != nil {
			return ast.Internal{}, err
		}
		if !ok {
			return ast.Internal{}, errAt(path, child.Line, "unrecognized internal statement: %q", child.Text)
		}
		if kind == "run" {
			in.Runs = append(in.Runs, caseName)
		} else {
			in.Checks = append(in.Checks, caseName)
		}
	}
	return in, nil
}
