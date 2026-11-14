//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what Encode — absorb Korean TANGEUL markdown into a fully tiled doc.Doc
package ko

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// Encode scans src with the shared scanner and parses it with the KO predicate-final
// parser, building both the ast.Document and a stream-order Block tiling in one pass:
// a SectionBlock per `## tangl:` heading, a StmtBlock per statement (with its AST
// StmtRef), and — via mergeProseKo — a verbatim ProseBlock for every remaining byte.
// src becomes Doc.Source in full and Doc.Locale is stamped "ko". The required 주제 and
// 사안 sections are enforced, and the result must pass doc.VerifyTiling — the encoder's
// own gate — which is returned to the caller.
func (KoCodec) Encode(src []byte, path string) (*doc.Doc, error) {
	lines := scan.ScanLines(src)
	_, secs := scan.ScanSections(lines)
	a := &ast.Document{Path: path}
	anchors, err := collectAnchorsKo(a, secs, path)
	if err != nil {
		return nil, err
	}
	if a.Subject == "" {
		return nil, errAtKo(path, 1, "tangl:주제 section is required")
	}
	if len(a.Cases) == 0 {
		return nil, errAtKo(path, 1, "tangl:사안 section is required")
	}
	d := &doc.Doc{
		AST:    a,
		Source: src,
		Blocks: mergeProseKo(anchors, len(src)),
		Locale: "ko",
	}
	return d, doc.VerifyTiling(d)
}
