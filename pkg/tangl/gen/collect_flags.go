//ff:func feature=tangl type=codegen control=sequence
//ff:what collectFlags — scans the Document once to decide optional imports/helpers needed
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectFlags scans the Document once to decide which optional imports
// and shared helper functions (tangl runtime, time, tanglCompare,
// tanglCaseActive) the generated file needs.
func collectFlags(doc *ast.Document) genFlags {
	return genFlags{
		NeedsTangl:      needsTangl(doc),
		NeedsTime:       needsTime(doc),
		NeedsCompare:    len(doc.Rules) > 0,
		NeedsCaseHelper: needsCaseHelper(doc),
	}
}
