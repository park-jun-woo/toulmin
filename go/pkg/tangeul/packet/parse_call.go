//ff:func feature=tangl type=util control=sequence
//ff:what parseCall — rebuild an ast.Call (predicate position) from an ExprCall record [alias][name]
package packet

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCall rebuilds an ast.Call used as a predicate term from an ExprCall
// record body (an alias label index and a name label index) at *pos, advancing
// *pos and resolving both through labels. It returns node count 1. An
// out-of-range index is a hard error.
func parseCall(recs []core.Word, pos *int, labels []string) (ast.Expr, int, error) {
	if len(recs) < *pos+2 {
		return nil, 0, errors.New("tangeul/packet: EXPR Call record truncated")
	}
	alias := int(recs[*pos])
	name := int(recs[*pos+1])
	*pos += 2
	if alias >= len(labels) || name >= len(labels) {
		return nil, 0, errors.New("tangeul/packet: EXPR Call label out of range")
	}
	return ast.Call{Fn: ast.Ref{Alias: labels[alias], Name: labels[name]}}, 1, nil
}
