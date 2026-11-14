//ff:type feature=tangl type=model
//ff:what Doc — the in-memory TANGEUL container (AST + Source Text + Block tiling)
package doc

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// Doc is the in-memory representation of a TANGEUL stream: the parsed AST graph
// (the existing toolchain's input, kept unmodified per README contract §2), the
// authored Source Text captured verbatim (NFC not applied), and the stream-order
// Block tiling that carries prose spans and statement anchor spans.
//
// AST identity is preserved so validate/effects/gen keep consuming *ast.Document
// unchanged; prose losslessness lives entirely in Source + Blocks.
type Doc struct {
	AST    *ast.Document // graph (existing toolchain input)
	Source []byte        // Source Text — authored bytes verbatim
	Blocks []Block       // stream-order tiling: PROSE spans + statement anchor spans
	Locale string        // authoring locale tag ("en"/"ko"/""), stamped into the STREAM_HEADER (§11-11)
}
