//ff:func feature=tangl type=util control=sequence
//ff:what emitDefinition — dispatch a Definitions statement to the term declaration or struct field
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitDefinition routes a Definitions coordinate: a two-element path is a struct
// `has` field, a one-element path is the term declaration itself.
func emitDefinition(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	if len(ref.Path) == 2 {
		return emitStructField(e, ref, elem.(ast.Field), span)
	}
	return emitTermDecl(e, ref, elem.(ast.Definition), span)
}
