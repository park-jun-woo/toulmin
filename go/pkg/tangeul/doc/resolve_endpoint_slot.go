//ff:func feature=tangl type=util control=selection
//ff:what resolveEndpointSlot — resolve an endpoint body slot (0 Require·1 Run·2 Check)
package doc

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// resolveEndpointSlot returns the slotIdx-th element of the endpoint body slice
// selected by slotKind (0 Require, 1 Run, 2 Check). Runs/Checks are []string
// without Line, so the encoder assigns slotIdx per kind and Blocks carry the
// original interleaving (README contract §2).
func resolveEndpointSlot(ep *ast.Endpoint, slotKind, slotIdx int) (any, error) {
	switch slotKind {
	case 0:
		if slotIdx < 0 || slotIdx >= len(ep.Requires) {
			return nil, fmt.Errorf("tangeul/doc: Endpoint Require index %d out of range (%d)", slotIdx, len(ep.Requires))
		}
		return ep.Requires[slotIdx], nil
	case 1:
		if slotIdx < 0 || slotIdx >= len(ep.Runs) {
			return nil, fmt.Errorf("tangeul/doc: Endpoint Run index %d out of range (%d)", slotIdx, len(ep.Runs))
		}
		return ep.Runs[slotIdx], nil
	case 2:
		if slotIdx < 0 || slotIdx >= len(ep.Checks) {
			return nil, fmt.Errorf("tangeul/doc: Endpoint Check index %d out of range (%d)", slotIdx, len(ep.Checks))
		}
		return ep.Checks[slotIdx], nil
	default:
		return nil, fmt.Errorf("tangeul/doc: unknown Endpoint slotKind %d", slotKind)
	}
}
