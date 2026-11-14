//ff:func feature=tangl type=engine control=selection
//ff:what materializeDefinition — turn one string-list or file Definition into a list Spec (other kinds skipped)
package materialize

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// materializeDefinition returns the StringListSpec for a StringListDef (inline
// values) or a FileSpecDef (loaded through registry). The bool is false for any
// other Kind, which Lists leaves out of its result.
func materializeDefinition(doc *ast.Document, def ast.Definition, registry DefinitionLoaderRegistry) (toulmin.StringListSpec, bool, error) {
	switch def.Kind {
	case ast.StringListDef:
		return toulmin.NewStringListSpec(def.Name, def.Strings), true, nil
	case ast.FileSpecDef:
		spec, err := loadFileDefinition(doc, def, registry)
		return spec, err == nil, err
	default:
		return toulmin.StringListSpec{}, false, nil
	}
}
