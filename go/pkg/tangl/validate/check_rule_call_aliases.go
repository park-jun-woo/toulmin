//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkRuleCallAliases — verifies host-call package aliases inside Rules are declared in See
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkRuleCallAliases verifies that every host Call inside a tangl:Rules
// condition whose FUNC_REF carries a package alias has that alias declared in
// the document's tangl:See section. It covers both predicate-position Calls
// (the "holds" form) and value-position Calls (the "result of" form, in
// Compare operands and an Exists source). A bare Call (no alias) is not
// resolved here; that is the consumer's two-stage lookup.
func checkRuleCallAliases(doc *ast.Document) []error {
	aliases := seeAliasSet(doc)
	var errs []error
	for _, r := range doc.Rules {
		errs = append(errs, ruleCallAliasErrs(doc, r, aliases)...)
	}
	return errs
}
