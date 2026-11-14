//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ruleAnchorsKo — parse the 규칙 section, anchoring each inline rule subtree as one StmtBlock
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// ruleAnchorsKo parses each non-struck 규칙 item into an ast.InlineRule and anchors
// the whole head-plus-nested-condition subtree as a single StmtBlock with ref
// {Rules, [ri]} (README contract §2 — the subtree never anchors deeper). Struck
// items become prose via mergeProseKo.
func ruleAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	ri := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		rule, err := parseInlineRuleKo(it, path)
		if err != nil {
			return nil, err
		}
		a.Rules = append(a.Rules, rule)
		endOff := absorbRuleSubtreeKo(it)
		anchors = append(anchors, doc.Block{
			Kind: doc.StmtBlock,
			Span: doc.Span{Off: it.Span.Off, Len: endOff - it.Span.Off},
			Ref:  doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{ri}},
		})
		ri++
	}
	return anchors, nil
}
