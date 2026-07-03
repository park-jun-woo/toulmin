//ff:func feature=tangl type=util control=sequence
//ff:what termSpecRef — resolve a const term's `as` spec reference from its sameAs edge
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// termSpecRef returns the `as <ref>` spec reference of a const definition, or nil
// when there is none. The sameAs edge points at a Func node that carries the full
// reference spelling.
func termSpecRef(d *dec, rel []packet.EdgePacket) (*ast.Ref, error) {
	e, ok := findRel(rel, codebook.PropSameAs)
	if !ok {
		return nil, nil
	}
	ref, err := resolveFuncRef(d, uint16(e.DstTID))
	if err != nil {
		return nil, err
	}
	return &ref, nil
}
