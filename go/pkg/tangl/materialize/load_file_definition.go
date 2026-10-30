//ff:func feature=tangl type=engine control=sequence
//ff:what loadFileDefinition — resolve a file definition's loader and path, load it once, and wrap it as a list Spec
package materialize

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// loadFileDefinition materializes one FileSpecDef: a registry and a registered
// loader for def.LoadType are required, def.FilePath is resolved against the
// document, and the loader's snapshot becomes a StringListSpec named def.Name.
func loadFileDefinition(doc *ast.Document, def ast.Definition, registry DefinitionLoaderRegistry) (toulmin.StringListSpec, error) {
	if registry == nil {
		return toulmin.StringListSpec{}, defError(doc, def, "loader registry is required")
	}
	loader, ok := registry.LookupDefinitionLoader(def.LoadType)
	if !ok || loader == nil {
		return toulmin.StringListSpec{}, defError(doc, def, "loader %q is not registered", def.LoadType)
	}
	resolved, err := resolvePath(doc.Path, def.FilePath)
	if err != nil {
		return toulmin.StringListSpec{}, defWrap(doc, def, err)
	}
	snapshot, err := loader.LoadDefinition(DefinitionLoadRequest{
		DefinitionName: def.Name, TypeName: def.LoadType,
		SourcePath: def.FilePath, ResolvedPath: resolved,
	})
	if err != nil {
		return toulmin.StringListSpec{}, fmt.Errorf("%s:%d: definition %q (loader %q, path %q): %w", doc.Path, def.Line, def.Name, def.LoadType, def.FilePath, err)
	}
	return toulmin.NewStringListSpec(def.Name, snapshot.Strings), nil
}
