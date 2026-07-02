//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseRulesSection — parse the tangl:Rules section's inline rules
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseRulesSection parses each "`name` when <condition>" inline rule.
func parseRulesSection(sec section, path string) ([]ast.InlineRule, error) {
	items, err := parseItems(sec.Lines, sec.LineOffset, path)
	if err != nil {
		return nil, err
	}
	rules := make([]ast.InlineRule, 0, len(items))
	for _, it := range items {
		r, err := parseInlineRuleItem(it, path)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}
