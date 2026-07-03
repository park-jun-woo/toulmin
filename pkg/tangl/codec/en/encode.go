//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what Encode — absorb English TANGEUL markdown into a fully tiled doc.Doc
package en

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// Encode reuses parser.ParseSource for the AST, then rescans src with the shared
// scanner to build a stream-order Block tiling: a SectionBlock per `## tangl:`
// heading and a StmtBlock per statement (linked to its AST element by source line),
// with every remaining byte captured as verbatim ProseBlocks. src becomes
// Doc.Source in full and Doc.Locale is stamped "en". The result must pass
// doc.VerifyTiling — the encoder's own gate — which is returned to the caller.
func (ENCodec) Encode(src []byte, path string) (*doc.Doc, error) {
	a, err := parser.ParseSource(string(src), path)
	if err != nil {
		return nil, err
	}
	lines := scan.ScanLines(src)
	_, secs := scan.ScanSections(lines)
	refByLine := buildRefIndex(a)
	if err := matchBodyRefs(a, secs, refByLine); err != nil {
		return nil, err
	}
	anchors, err := collectAnchors(secs, refByLine)
	if err != nil {
		return nil, err
	}
	d := &doc.Doc{
		AST:    a,
		Source: src,
		Blocks: mergeProse(anchors, len(src)),
		Locale: "en",
	}
	return d, doc.VerifyTiling(d)
}
