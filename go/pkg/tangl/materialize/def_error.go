//ff:func feature=tangl type=engine control=sequence
//ff:what defError — format a definition error prefixed with document path, line, and definition name
package materialize

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// defError builds "<path>:<line>: definition "<name>": <message>".
func defError(doc *ast.Document, def ast.Definition, format string, args ...any) error {
	return fmt.Errorf("%s:%d: definition %q: %s", doc.Path, def.Line, def.Name, fmt.Sprintf(format, args...))
}
