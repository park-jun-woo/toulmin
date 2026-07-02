//ff:func feature=tangl type=validator control=sequence pattern=error-collection
//ff:what Validate — validates a Document, collecting every violation into one error
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// Validate runs every structural check on doc (duplicate names, dangling
// references, undo-without-do, and the checking/run/don't cycle bans) and
// collects every violation into a single error (pattern=error-collection).
// It returns nil when doc has no violations.
func Validate(doc *ast.Document) error {
	var errs []error
	errs = append(errs, checkDuplicateCaseNames(doc)...)
	errs = append(errs, checkDuplicateNodeNames(doc)...)
	errs = append(errs, checkDuplicateSeeAlias(doc)...)
	errs = append(errs, checkDuplicateDefNames(doc)...)
	errs = append(errs, checkDuplicateEndpointNames(doc)...)
	errs = append(errs, checkAttackRefs(doc)...)
	errs = append(errs, checkExecNodeRefs(doc)...)
	errs = append(errs, checkRunCaseRefs(doc)...)
	errs = append(errs, checkCheckingRefs(doc)...)
	errs = append(errs, checkEndpointCaseRefs(doc)...)
	errs = append(errs, checkInternalCaseRefs(doc)...)
	errs = append(errs, checkWithRefs(doc)...)
	errs = append(errs, checkUsingAliasRefs(doc)...)
	errs = append(errs, checkExecFuncAliasRefs(doc)...)
	errs = append(errs, checkUndoRequiresDo(doc)...)
	errs = append(errs, checkCheckingCycle(doc)...)
	errs = append(errs, checkRunCycle(doc)...)
	errs = append(errs, checkAttackCycle(doc)...)
	return newMultiError(errs)
}
