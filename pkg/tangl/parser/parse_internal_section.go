//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseInternalSection — parse the tangl:Internal section's triggers
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseInternalSection parses each `on <event>` / `every <interval>` trigger.
func parseInternalSection(sec section, path string) ([]ast.Internal, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	ins := make([]ast.Internal, 0, len(items))
	for _, it := range items {
		in, err := parseInternalItem(it, path)
		if err != nil {
			return nil, err
		}
		ins = append(ins, in)
	}
	return ins, nil
}
