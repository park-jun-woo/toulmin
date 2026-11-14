//ff:func feature=tangl type=parser control=sequence
//ff:what applyNodeClause — parse a node's trailing "[using ...|checking ...] [qualified F] [, level ...]" clauses
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// applyNodeClause parses the optional trailing clauses of a node declaration,
// node_decl = NAME "is" ("a"|"an") ROLE_EXPR ["using" ... | "checking" CASE_REF]
// [qual_clause] [level_clause]: an optional using/checking source
// (applyNodeSource), then an optional "qualified F" (when the using clause has
// not already consumed one), then an optional level_clause (", level fail" /
// ", level review"). A canonical Case binding of a Rules-declared predicate
// therefore reads "`x` is a counter rule qualified 0.8, level fail".
func applyNodeClause(n *ast.Node, rest string, line int, path string) error {
	rest, err := applyNodeSource(n, rest, line, path)
	if err != nil {
		return err
	}
	rest = strings.TrimSpace(rest)
	if r, ok := takeKeyword(rest, "qualified"); ok && n.Qualified == nil {
		q, tail, err := parseQualifiedClause(r, path, line)
		if err != nil {
			return err
		}
		n.Qualified, rest = &q, strings.TrimSpace(tail)
	}
	if r, ok := takeKeyword(rest, ","); ok {
		lvl, tail, err := parseLevelClause(r, path, line)
		if err != nil {
			return err
		}
		n.Level, rest = lvl, tail
	}
	if strings.TrimSpace(rest) != "" {
		return errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return nil
}
