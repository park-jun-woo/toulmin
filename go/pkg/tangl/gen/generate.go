//ff:func feature=tangl type=codegen control=sequence
//ff:what Generate — compiles a TANGEUL v0.3 Document into a single Go source file
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// Generate compiles doc into a single Go source file string: definitions,
// inline rules, case graph builders, provided endpoints, and internal
// triggers, in that document order, formatted via go/format.Source.
func Generate(doc *ast.Document) (string, error) {
	return GenerateWithOptions(doc, GenerateOptions{})
}
