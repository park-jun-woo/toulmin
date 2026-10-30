//ff:type feature=tangl type=interface
//ff:what DefinitionResolver — optional Resolver extension instantiating a Definition's external Spec
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// DefinitionResolver is an optional extension implemented by loaders that can
// instantiate a Definition's external Spec constructor for dynamic evaluation.
// CompileDoc uses it to mirror gen's rule-owned automatic With bindings.
type DefinitionResolver interface {
	DefinitionSpec(ast.Definition) (toulmin.Spec, error)
}
