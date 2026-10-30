//ff:func feature=tangl type=engine control=sequence
//ff:what CompileDoc — compiles every doc rule with document-local-first resolution
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
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
