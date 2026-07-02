//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCasesSection — parse the tangl:Cases section's case blocks
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseCasesSection parses each "in case of `name`" block.
func parseCasesSection(sec section, path string) ([]ast.Case, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	cases := make([]ast.Case, 0, len(items))
	for _, it := range items {
		c, err := parseCaseItem(it, path)
		if err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}
