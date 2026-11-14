//ff:func feature=tangl type=engine control=iteration dimension=2
//ff:what resolveRuleSpecs — instantiate rule-owned definition Specs through an optional loader contract
package eval

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func resolveRuleSpecs(doc *ast.Document, res Resolver, lists map[string]toulmin.StringListSpec) (map[string]toulmin.Specs, error) {
	loader, ok := res.(DefinitionResolver)
	defs := make(map[string]ast.Definition, len(doc.Defs))
	for _, def := range doc.Defs {
		defs[def.Name] = def
	}
	out := make(map[string]toulmin.Specs)
	for _, rule := range doc.Rules {
		terms := uniqueTerms(rule.With, ast.RuleSpecNames(rule))
		for _, term := range terms {
			if list, exists := lists[term]; exists {
				out[rule.Name] = append(out[rule.Name], list)
				continue
			}
			def, exists := defs[term]
			if !exists || def.SpecRef == nil || !ok {
				continue
			}
			sp, err := loader.DefinitionSpec(def)
			if err != nil {
				return nil, fmt.Errorf("rule %q definition %q: %w", rule.Name, term, err)
			}
			if sp == nil {
				return nil, fmt.Errorf("rule %q definition %q produced a nil Spec", rule.Name, term)
			}
			out[rule.Name] = append(out[rule.Name], sp)
		}
	}
	return out, nil
}
