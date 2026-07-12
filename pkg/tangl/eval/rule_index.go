//ff:type feature=tangl type=model
//ff:what ruleIndex — per-document lookup tables (rule set, line, node, ordered names)
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// ruleIndex holds the lookup tables CompileDoc builds once from a document's
// rules: the set of rule names, each rule's declaration line, each rule's AST
// node by name, and the rule names in document order.
type ruleIndex struct {
	ruleSet map[string]bool
	lineOf  map[string]int
	byName  map[string]ast.InlineRule
	names   []string
}
