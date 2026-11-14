//ff:type feature=tangl type=model
//ff:what CompileOptions — predicate/value Resolver and file-definition loaders for CompileDocument
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"

// CompileOptions supplies predicate/value resolution and file-definition loaders.
type CompileOptions struct {
	Resolver Resolver
	Loaders  materialize.DefinitionLoaderRegistry
}
