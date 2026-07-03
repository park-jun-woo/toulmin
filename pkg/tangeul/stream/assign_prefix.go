//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what assignPrefix — assign Text TIDs to prefix payloads, recovering a Rules EXPR by position
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// assignPrefix walks a statement's prefix payloads in order. When isRules is set,
// the first payload is the condition EXPR (decoded, not counted); every other
// payload is a Text whose TID is issued from the shared counter, matching the
// Writer so label/means/with references resolve. It returns the Rules condition,
// or nil.
func assignPrefix(d *dec, prefix []rawPacket, isRules bool) (ast.Expr, error) {
	var cond ast.Expr
	for idx, rp := range prefix {
		got, err := decodePrefixItem(d, rp, isRules && idx == 0)
		if err != nil {
			return nil, err
		}
		if got != nil {
			cond = got
		}
	}
	return cond, nil
}
