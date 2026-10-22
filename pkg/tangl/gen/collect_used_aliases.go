//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what collectUsedAliases — scans the whole Document for every See alias actually referenced
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectUsedAliases scans the whole Document for every See alias
// actually referenced — by a Definitions SpecRef, a case node's "using"
// clause, a do/undo Ref, or a host Call inside a tangl:Rules condition —
// so buildImports only imports aliases the generated code really calls.
func collectUsedAliases(doc *ast.Document) map[string]bool {
	used := make(map[string]bool)
	collectRuleCallAliases(doc.Rules, used)
	for _, rule := range doc.Rules {
		if rule.Using != nil && rule.Using.Alias != "" {
			used[rule.Using.Alias] = true
		}
	}
	for _, def := range doc.Defs {
		if def.SpecRef != nil && def.SpecRef.Alias != "" {
			used[def.SpecRef.Alias] = true
		}
	}
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Using != nil && n.Using.Alias != "" {
				used[n.Using.Alias] = true
			}
		}
		for _, e := range c.Execs {
			if e.Func != nil && e.Func.Alias != "" {
				used[e.Func.Alias] = true
			}
		}
	}
	return used
}
