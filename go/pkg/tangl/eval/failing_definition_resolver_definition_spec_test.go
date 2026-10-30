//ff:func feature=tangl type=engine control=sequence
//ff:what DefinitionSpec — failingDefinitionResolver returns its error, or a nil Spec when err is nil
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func (f failingDefinitionResolver) DefinitionSpec(ast.Definition) (toulmin.Spec, error) {
	return nil, f.err
}
