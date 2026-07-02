//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseSeeSection — parse the tangl:See section's package references
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseSeeSection parses each "see `alias` from `package`" entry.
func parseSeeSection(sec section, path string) ([]ast.See, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	sees := make([]ast.See, 0, len(items))
	for _, it := range items {
		see, err := parseSeeItem(it, path)
		if err != nil {
			return nil, err
		}
		sees = append(sees, see)
	}
	return sees, nil
}
