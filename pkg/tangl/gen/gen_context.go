//ff:type feature=tangl type=codegen
//ff:what genContext — shared state threaded through one Generate call
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// genContext carries the state a single Generate call threads through its
// render functions: the source Document, the resolved Definitions index
// (name -> defInfo), the feature flags that gate optional imports and
// shared helpers, and the dedup'd "checking" wrapper function names keyed
// by target case name.
type genContext struct {
	Doc           *ast.Document
	Defs          map[string]defInfo
	Flags         genFlags
	CheckWrappers map[string]string
}
