//ff:func feature=tangl type=engine control=sequence
//ff:what CompileDoc — compiles every doc rule with document-local-first resolution
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// CompileDoc compiles every tangl:Rules entry of doc into a rule-func closure.
// A bare predicate Call resolves to a sibling document rule first, else to
// res.Predicate. Sibling references are checked for cycles (rejected as a
// compile error) and compiled in dependency-first order so every closure is
// fully bound with no runtime lookup path.
func CompileDoc(doc *ast.Document, res Resolver) (map[string]func(toulmin.Context, toulmin.Specs) (bool, any), error) {
	return CompileDocument(doc, CompileOptions{Resolver: res})
}

// CompileOptions supplies predicate/value resolution and file-definition loaders.
type CompileOptions struct {
	Resolver Resolver
	Loaders  materialize.DefinitionLoaderRegistry
}

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
