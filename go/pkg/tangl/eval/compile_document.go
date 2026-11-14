//ff:func feature=tangl type=engine control=sequence
//ff:what CompileDocument — materialize definitions, then compile every doc rule with derived rule Specs
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// CompileDocument compiles a document after eagerly materializing its definitions.
func CompileDocument(doc *ast.Document, opts CompileOptions) (map[string]func(toulmin.Context, toulmin.Specs) (bool, any), error) {
	lists, err := materialize.Lists(doc, opts.Loaders)
	if err != nil {
		return nil, err
	}
	idx := buildRuleIndex(doc)
	edges := buildEdges(doc.Rules, idx.ruleSet)
	if err := detectRuleCycle(doc.Path, edges, func(n string) int { return idx.lineOf[n] }); err != nil {
		return nil, err
	}
	compiled := map[string]evalFunc{}
	dr := docResolver{host: opts.Resolver, ruleSet: idx.ruleSet, compiled: compiled}
	if err := compileRules(topoOrder(idx.names, edges), idx.byName, dr, compiled); err != nil {
		return nil, err
	}
	derived, err := resolveRuleSpecs(doc, opts.Resolver, lists)
	if err != nil {
		return nil, err
	}
	return wrapClosuresWithSpecs(doc.Rules, compiled, derived), nil
}
