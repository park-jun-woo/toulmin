//ff:func feature=tangl type=codegen control=sequence
//ff:what LookupDefinitionLoader — genListRegistry resolves any loader type to genListLoader
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"

func (genListRegistry) LookupDefinitionLoader(string) (materialize.DefinitionLoader, bool) {
	return genListLoader{}, true
}
