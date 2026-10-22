//ff:func feature=tangl type=util control=selection
//ff:what assembleDoc — test-only: dispatch a fixture name to its flat Doc builder
package tangeul

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// assembleDoc pairs a parsed AST with a flat, hand-tiled Block layout for the named
// step-1 fixture, producing a Doc ready for Save. It is a test-only stand-in for
// the tangl-en encoder: the Blocks carry no prose (each is a 4-byte placeholder
// span that tiles the Source exactly so doc.VerifyTiling holds), which is enough
// for AST/validate/effects/gen parity but deliberately not for mode-1 byte
// identity. The builders are per fixture (not a generic AST walk) to stay flat and
// obviously test-scoped rather than a second encoder.
func assembleDoc(name string, a *ast.Document) *doc.Doc {
	switch name {
	case "transfer":
		return assembleTransferDoc(a)
	case "access":
		return assembleAccessDoc(a)
	default:
		return nil
	}
}
