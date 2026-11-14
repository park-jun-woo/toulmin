//ff:type feature=tangl type=model
//ff:what GenerateOptions — file-definition loaders used at generation time
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"

// GenerateOptions supplies file-definition loaders used at generation time.
type GenerateOptions struct {
	Loaders materialize.DefinitionLoaderRegistry
}
