//ff:func feature=tangl type=codegen control=sequence
//ff:what GenerateWithOptions — snapshot file definitions, then emit the document as one formatted Go file
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"
)

// GenerateWithOptions generates source after eagerly snapshotting file definitions.
func GenerateWithOptions(doc *ast.Document, opts GenerateOptions) (string, error) {
	materialized, err := materialize.Lists(doc, opts.Loaders)
	if err != nil {
		return "", err
	}
	lists := snapshotLists(materialized)
	gc := &genContext{Doc: doc, Flags: collectFlags(doc), Rules: indexRulesByName(doc.Rules)}
	imports, err := buildImports(gc)
	if err != nil {
		return "", err
	}
	var body strings.Builder
	gc.Defs = buildDefinitions(&body, doc.Defs, lists)
	writeRuleHelpers(&body, gc.Flags)
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
