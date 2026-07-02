//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseProvidesSection — parse the tangl:Provides section's endpoints
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseProvidesSection parses each "provides `name`" endpoint entry.
func parseProvidesSection(sec section, path string) ([]ast.Endpoint, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	eps := make([]ast.Endpoint, 0, len(items))
	for _, it := range items {
		ep, err := parseEndpointItem(it, path)
		if err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}
	return eps, nil
}
