//ff:type feature=tangl type=codegen
//ff:what defInfo — how one Definitions entry compiled to Go (const/Spec identifiers, kind)
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// defInfo records how a tangl:Definitions entry was compiled to Go: its
// constant identifier (ConstDef only), its derived Spec variable
// identifier (ConstDef with a SpecRef only), and its original Kind.
type defInfo struct {
	Const string
	Spec  string
	Kind  ast.DefKind
}
