//ff:func feature=tangl type=parser control=sequence
//ff:what parseUsingClause — parse "FUNC_REF [with_clause] [qual_clause]"
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseUsingClause parses the function reference following "using", plus
// its optional "with" and "qualified" clauses.
func parseUsingClause(s string, path string, line int) (ast.Ref, []string, *float64, string, error) {
	ref, rest, ok := parseRef(s)
	if !ok {
		return ast.Ref{}, nil, nil, "", errAt(path, line, "expected function reference after 'using'")
	}
	var with []string
	rest = strings.TrimSpace(rest)
	if r2, ok2 := takeKeyword(rest, "with"); ok2 {
		terms, tail, err := parseWithClause(r2, path, line)
		if err != nil {
			return ast.Ref{}, nil, nil, "", err
		}
		with, rest = terms, tail
	}
	var qual *float64
	rest = strings.TrimSpace(rest)
	if r3, ok3 := takeKeyword(rest, "qualified"); ok3 {
		q, tail, err := parseQualifiedClause(r3, path, line)
		if err != nil {
			return ast.Ref{}, nil, nil, "", err
		}
		qual, rest = &q, tail
	}
	return ref, with, qual, rest, nil
}
