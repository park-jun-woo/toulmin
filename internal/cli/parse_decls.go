//ff:func feature=scanner type=parser control=iteration dimension=1
//ff:what parseDecls — converts RuleDecls to RuleMetas via annotation parsing
package cli

import (
	"github.com/park-jun-woo/toulmin/internal/scanner"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// parseDecls parses annotation lines from each RuleDecl into RuleMeta,
// setting Name to the function name.
func parseDecls(decls []scanner.RuleDecl) []toulmin.RuleMeta {
	var metas []toulmin.RuleMeta
	for _, d := range decls {
		meta := toulmin.ParseAnnotation(d.Lines)
		meta.Name = d.FuncName
		metas = append(metas, meta)
	}
	return metas
}
