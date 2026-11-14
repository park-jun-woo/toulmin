//ff:func feature=tangl type=parser control=selection
//ff:what RenderStmt — regenerate one AST element as its English canonical statement text (doc.Codec)
package en

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// Compile-time assertion that ENCodec fully satisfies the doc.Codec contract
// (Locale + Encode from earlier phases, RenderStmt completed here).
var _ doc.Codec = (*ENCodec)(nil)

// RenderStmt resolves ref to its AST element via doc.ResolveStmt and regenerates
// the English canonical statement text (README contract §3). The result is a
// single line except for a Rules inline rule, whose nested condition tree renders
// multi-line; the core renderer owns list markers and section indentation. A bare
// string is either the Subject (Sec Subject) or a Line-less run/check case name,
// disambiguated by ref.Sec and the slotKind at Path[1].
func (ENCodec) RenderStmt(d *doc.Doc, ref doc.StmtRef) (string, error) {
	el, err := doc.ResolveStmt(d, ref)
	if err != nil {
		return "", err
	}
	switch v := el.(type) {
	case string:
		if ref.Sec == toulmin.SectionSubject {
			return renderSubject(v), nil
		}
		if len(ref.Path) != 3 {
			return "", fmt.Errorf("en: string element at %v has no run/check slotKind", ref.Path)
		}
		return renderRunCheck(ref.Path[1], v)
	case ast.See:
		return renderSee(v), nil
	case ast.Definition:
		return renderDefinition(v), nil
	case ast.Field:
		return renderField(v), nil
	case ast.InlineRule:
		return renderRule(v), nil
	case ast.Case:
		return renderCase(v), nil
	case ast.Require:
		return renderRequire(v), nil
	case ast.Node:
		return renderNode(v), nil
	case ast.Attack:
		return renderAttack(v), nil
	case ast.Supersede:
		return renderSupersede(v), nil
	case ast.Exec:
		return renderExec(v)
	case ast.Endpoint:
		return renderEndpoint(v), nil
	case ast.Internal:
		return renderInternal(v)
	default:
		return "", fmt.Errorf("en: RenderStmt: unsupported AST element %T", el)
	}
}
