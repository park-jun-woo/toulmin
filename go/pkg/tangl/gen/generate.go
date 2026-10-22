//ff:func feature=tangl type=codegen control=sequence
//ff:what Generate — compiles a TANGEUL v0.3 Document into a single Go source file
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/materialize"
)

// Generate compiles doc into a single Go source file string: definitions,
// inline rules, case graph builders, provided endpoints, and internal
// triggers, in that document order, formatted via go/format.Source.
func Generate(doc *ast.Document) (string, error) {
	return GenerateWithOptions(doc, GenerateOptions{})
}

// GenerateOptions supplies file-definition loaders used at generation time.
type GenerateOptions struct {
	Loaders materialize.DefinitionLoaderRegistry
}

// GenerateWithOptions generates source after eagerly snapshotting file definitions.
func GenerateWithOptions(doc *ast.Document, opts GenerateOptions) (string, error) {
	materialized, err := materialize.Lists(doc, opts.Loaders)
	if err != nil {
		return "", err
	}
	lists := make(map[string][]string, len(materialized))
	for name, spec := range materialized {
		lists[name] = spec.StringList()
	}
	gc := &genContext{Doc: doc, Flags: collectFlags(doc), Rules: make(map[string]ast.InlineRule, len(doc.Rules))}
	for _, rule := range doc.Rules {
		gc.Rules[rule.Name] = rule
	}
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
