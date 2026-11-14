//ff:func feature=tangl type=codegen control=sequence
//ff:what collectFlags — scans the Document once to decide optional imports/helpers needed
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// collectFlags scans the Document once to decide which optional imports
// and shared helper functions (tangl runtime, time, the tangl* rule
// evaluation helpers, tanglCaseActive) the generated file needs.
func collectFlags(doc *ast.Document) genFlags {
	f := genFlags{
		NeedsTangl:      needsTangl(doc),
		NeedsTime:       needsTime(doc),
		NeedsCaseHelper: needsCaseHelper(doc),
	}
	collectRuleFlags(doc.Rules, &f)
	return f
}
