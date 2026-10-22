//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what assignPrefix — assign Text TIDs to prefix payloads, recovering a Rules EXPR by position
package stream

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// assignPrefix decodes a statement's prefix records: each TEXT is recorded under
// its Edge Id, and a PREDICATE-DECLARE (present only for a Rules condition rule)
// yields the condition tree, which is returned. The record codec, not its
// position, tells the two apart, so a host rule's prefix (TEXTs only) reads too.
func assignPrefix(d *dec, prefix []rawPacket) (ast.Expr, error) {
	var cond ast.Expr
	for _, rp := range prefix {
		got, err := decodePrefixItem(d, rp)
		if err != nil {
			return nil, err
		}
		if got != nil {
			cond = got
		}
	}
	return cond, nil
}
