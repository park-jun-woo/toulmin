//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what Lists — eagerly materialize every string-list and file-backed Definition exactly once

// Package materialize resolves syntax-directed Definition values before evaluation or generation.
package materialize

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// Lists eagerly materializes every syntax-directed list Definition exactly once.
func Lists(doc *ast.Document, registry DefinitionLoaderRegistry) (map[string]toulmin.StringListSpec, error) {
	out := make(map[string]toulmin.StringListSpec)
	for _, def := range doc.Defs {
		spec, ok, err := materializeDefinition(doc, def, registry)
		if err != nil {
			return nil, err
		}
		if ok {
			out[def.Name] = spec
		}
	}
	return out, nil
}
