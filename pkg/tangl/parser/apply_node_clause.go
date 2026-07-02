//ff:func feature=tangl type=parser control=sequence
//ff:what applyNodeClause — parse a node's trailing "using ..." or "checking ..." clause
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// applyNodeClause parses the optional trailing clause of a node declaration:
// either "using FUNC_REF [with_clause] [qual_clause]" or "checking NAME".
func applyNodeClause(n *ast.Node, rest string, line int, path string) error {
	if using, ok := takeKeyword(rest, "using"); ok {
		ref, with, qual, tail, err := parseUsingClause(using, path, line)
		if err != nil {
			return err
		}
		n.Using, n.With, n.Qualified = &ref, with, qual
		rest = tail
	} else if checking, ok := takeKeyword(rest, "checking"); ok {
		caseName, tail, err := parseCheckingClause(checking, path, line)
		if err != nil {
			return err
		}
		n.Checking = caseName
		rest = tail
	} else {
		return errAt(path, line, "unexpected trailing text: %q", rest)
	}
	if strings.TrimSpace(rest) != "" {
		return errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return nil
}
