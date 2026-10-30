//ff:func feature=tangl type=model control=sequence
//ff:what RuleSpecNames — collect first-step SpecRef names from a rule condition
package ast

// RuleSpecNames returns the first step of every SpecRef in a condition, once,
// in source traversal order. These names form an expression rule's inferred
// definition bindings.
func RuleSpecNames(rule InlineRule) []string {
	return collectSpecNames(rule.Cond, map[string]bool{}, nil)
}
