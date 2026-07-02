//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what collectUsedAliases — scans the whole Document for every See alias actually referenced
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectUsedAliases scans the whole Document for every See alias
// actually referenced — by a Definitions SpecRef, a case node's "using"
// clause, or a do/undo Ref — so buildImports only imports aliases the
// generated code really calls.
func collectUsedAliases(doc *ast.Document) map[string]bool {
	used := make(map[string]bool)
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
