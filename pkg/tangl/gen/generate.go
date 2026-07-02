//ff:func feature=tangl type=codegen control=sequence
//ff:what Generate — compiles a TANGL v0.3 Document into a single Go source file
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// Generate compiles doc into a single Go source file string: definitions,
// inline rules, case graph builders, provided endpoints, and internal
// triggers, in that document order, formatted via go/format.Source.
func Generate(doc *ast.Document) (string, error) {
	gc := &genContext{Doc: doc, Flags: collectFlags(doc)}
	imports, err := buildImports(gc)
	if err != nil {
		return "", err
	}
	var body strings.Builder
	gc.Defs = buildDefinitions(&body, doc.Defs)
	if gc.Flags.NeedsCompare {
		writeCompareHelper(&body)
	}
	if gc.Flags.NeedsCaseHelper {
		writeCaseHelper(&body)
	}
	if err := buildRules(&body, doc.Rules); err != nil {
		return "", err
	}
	gc.CheckWrappers = buildCheckingWrappers(&body, doc)
	if err := buildCases(&body, gc); err != nil {
		return "", err
	}
	buildProvides(&body, gc)
	buildInternals(&body, doc)
	var out strings.Builder
	writeHeader(&out, doc.Subject)
	writeImportBlock(&out, imports)
	fmt.Fprintln(&out)
	out.WriteString(body.String())
	return formatSource(out.String())
}
