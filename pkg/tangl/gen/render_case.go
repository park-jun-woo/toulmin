//ff:func feature=tangl type=codegen control=sequence
//ff:what renderCase — writes one case's graph-builder function and package-level var
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderCase writes one case's graph-builder function — node
// registration, defeat edges, and do/undo/run execution attachment, in
// that document order — followed by its package-level graph var.
func renderCase(w *strings.Builder, gc *genContext, c ast.Case) error {
	subject := gc.Doc.Subject
	fnName := "new" + goIdentExported(c.Name) + "Graph"
	varName := goIdent(c.Name) + "Graph"
	writeCaseDoc(w, c)
	fmt.Fprintf(w, "func %s() *toulmin.Graph {\n", fnName)
	fmt.Fprintf(w, "\tg := toulmin.NewGraph(%s)\n", quoteGraphName(subject, c.Name))
	nodes := registerCaseNodes(w, gc, c)
	registerCaseAttacks(w, c, nodes)
	if err := registerCaseExecs(w, subject, c, nodes); err != nil {
		return err
	}
	writeUnusedNodeGuards(w, c, nodes)
	fmt.Fprintln(w, "\treturn g")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "var %s = %s()\n", varName, fnName)
	fmt.Fprintln(w)
	return nil
}
