//ff:func feature=tangl type=parser control=sequence
//ff:what parseCheckingClause — parse the CASE_REF of a "checking" clause
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseCheckingClause parses the CASE_REF following "checking": either a
// qualified `alias`.`case` reference or a bare local `case` reference.
// It reuses parseRef (the same grammar as FUNC_REF) since CASE_REF and
// FUNC_REF are structurally identical.
func parseCheckingClause(s string, path string, line int) (ast.Ref, string, error) {
	ref, rest, ok := parseRef(s)
	if !ok {
		return ast.Ref{}, "", errAt(path, line, "expected backtick-quoted case name after 'checking'")
	}
	return ref, rest, nil
}
