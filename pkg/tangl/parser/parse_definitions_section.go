//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseDefinitionsSection — parse the tangl:Definitions section's terms
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseDefinitionsSection parses each constant or struct term Definition.
func parseDefinitionsSection(sec section, path string) ([]ast.Definition, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	defs := make([]ast.Definition, 0, len(items))
	for _, it := range items {
		def, err := parseDefinitionItem(it, path)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	return defs, nil
}
