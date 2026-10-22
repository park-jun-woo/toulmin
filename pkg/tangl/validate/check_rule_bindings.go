//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkRuleBindings — reject ambiguous rule declarations and case bindings
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

func checkRuleBindings(doc *ast.Document) []error {
	rules := make(map[string]ast.InlineRule, len(doc.Rules))
	var errs []error
	for _, rule := range doc.Rules {
		rules[rule.Name] = rule
		if (rule.Cond == nil) == (rule.Using == nil) {
			errs = append(errs, errAt(doc.Path, rule.Line, "rule %q must declare exactly one of a condition or host function", rule.Name))
		}
	}
	for _, c := range doc.Cases {
		for _, node := range c.Nodes {
			_, declared := rules[node.Name]
			if node.Checking != nil && declared {
				errs = append(errs, errAt(doc.Path, node.Line, "case %q: checking node %q conflicts with a rule declaration", c.Name, node.Name))
			}
			if declared && (node.Using != nil || len(node.With) > 0 || len(node.Grounds) > 0) {
				errs = append(errs, errAt(doc.Path, node.Line, "case %q: node %q repeats predicate properties owned by its rule declaration", c.Name, node.Name))
			}
		}
	}
	return errs
}
