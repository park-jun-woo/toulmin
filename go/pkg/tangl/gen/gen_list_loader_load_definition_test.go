//ff:func feature=tangl type=codegen control=sequence
//ff:what LoadDefinition — genListLoader always yields the snapshot ["loaded"]
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"

func (genListLoader) LoadDefinition(materialize.DefinitionLoadRequest) (materialize.DefinitionSnapshot, error) {
	return materialize.DefinitionSnapshot{Strings: []string{"loaded"}}, nil
}
