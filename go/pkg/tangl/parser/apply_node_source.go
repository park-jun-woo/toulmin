//ff:func feature=tangl type=parser control=sequence
//ff:what applyNodeSource — parse a node's optional "using FUNC_REF ..." or "checking CASE_REF" source clause
package parser

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// applyNodeSource parses the optional predicate source of a node declaration —
// "using FUNC_REF [with_clause] [qual_clause] [grounded_clause]" (legacy
// migration input) or "checking CASE_REF" (verdict composition) — onto n and
// returns the unconsumed remainder. Text that starts with neither keyword is
// returned unchanged for the qualifier/level clauses.
func applyNodeSource(n *ast.Node, rest string, line int, path string) (string, error) {
	if using, ok := takeKeyword(rest, "using"); ok {
		ref, with, qual, grounds, tail, err := parseUsingClause(using, path, line)
		if err != nil {
			return "", err
		}
		n.Using, n.With, n.Qualified, n.Grounds = &ref, with, qual, grounds
		return tail, nil
	}
	if checking, ok := takeKeyword(rest, "checking"); ok {
		ref, tail, err := parseCheckingClause(checking, path, line)
		if err != nil {
			return "", err
		}
		n.Checking = &ref
		return tail, nil
	}
	return rest, nil
}
