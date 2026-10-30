//ff:func feature=tangl type=engine control=sequence
//ff:what defWrap — wrap a cause error with document path, line, and definition name
package materialize

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// defWrap builds "<path>:<line>: definition "<name>": <cause>" keeping cause unwrappable.
func defWrap(doc *ast.Document, def ast.Definition, cause error) error {
	return fmt.Errorf("%s:%d: definition %q: %w", doc.Path, def.Line, def.Name, cause)
}
