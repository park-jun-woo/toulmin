//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what indexRulesByName — map each inline rule by its name (a later duplicate overwrites)
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// indexRulesByName indexes rules by name for case-binding lookup.
func indexRulesByName(rules []ast.InlineRule) map[string]ast.InlineRule {
	byName := make(map[string]ast.InlineRule, len(rules))
	for _, rule := range rules {
		byName[rule.Name] = rule
	}
	return byName
}
