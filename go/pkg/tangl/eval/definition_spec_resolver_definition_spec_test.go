//ff:func feature=tangl type=engine control=sequence
//ff:what DefinitionSpec — test DefinitionResolver returning a derivedSpec named after the definition
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func (definitionSpecResolver) DefinitionSpec(def ast.Definition) (toulmin.Spec, error) {
	return &derivedSpec{name: def.Name}, nil
}
