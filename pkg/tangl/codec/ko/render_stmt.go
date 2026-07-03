//ff:func feature=tangl type=parser control=selection
//ff:what RenderStmt — regenerate one AST element as its Korean canonical statement text (doc.Codec)
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// Compile-time assertion that KoCodec fully satisfies the doc.Codec contract
// (Locale + Encode from earlier phases, RenderStmt completed here).
var _ doc.Codec = (*KoCodec)(nil)

// RenderStmt resolves ref to its AST element via doc.ResolveStmt and regenerates the
// Korean canonical statement text (README contract §3). The result is a single line
// except for a Rules inline rule, whose nested condition tree renders multi-line; the
// core renderer owns list markers and section indentation. A bare string is either the
// Subject (Sec 주제) or a Line-less run/check case name, disambiguated by ref.Sec and the
// slotKind at Path[1].
func (KoCodec) RenderStmt(d *doc.Doc, ref doc.StmtRef) (string, error) {
	el, err := doc.ResolveStmt(d, ref)
	if err != nil {
		return "", err
	}
	switch v := el.(type) {
	case string:
		if ref.Sec == codebook.SectionSubject {
			return renderSubjectKo(v), nil
		}
		if len(ref.Path) != 3 {
			return "", fmt.Errorf("ko: string element at %v has no run/check slotKind", ref.Path)
		}
		return renderRunCheckKo(ref.Path[1], v)
	case ast.See:
		return renderSeeKo(v), nil
	case ast.Definition:
		return renderDefinitionKo(v), nil
	case ast.Field:
		return renderFieldKo(v), nil
	case ast.InlineRule:
		return renderInlineRuleKo(v)
	case ast.Case:
		return renderCaseKo(v), nil
	case ast.Require:
		return renderRequireKo(v), nil
	case ast.Node:
		return renderNodeKo(v), nil
	case ast.Attack:
		return renderAttackKo(v), nil
	case ast.Supersede:
		return renderSupersedeKo(v), nil
	case ast.Exec:
		return renderExecKo(v)
	case ast.Endpoint:
		return renderEndpointKo(v), nil
	case ast.Internal:
		return renderInternalKo(v)
	default:
		return "", fmt.Errorf("ko: RenderStmt: unsupported AST element %T", el)
	}
}
