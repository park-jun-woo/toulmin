//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseDefinitionFields — parse a Definitions struct's "has `f` as Type" child fields
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseDefinitionFields parses each StructDef child field in order, stopping
// at the first error.
func parseDefinitionFields(children []item, path string) ([]ast.Field, error) {
	fields := make([]ast.Field, 0, len(children))
	for _, child := range children {
		f, err := parseFieldItem(child, path)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}
