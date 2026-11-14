//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDefinitionKinds — reject definitions whose populated fields do not match their Kind
package validate

import (
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// checkDefinitionKinds reports one error per definition whose fields are invalid
// for its Kind. A well-formed file definition with a relative path additionally
// requires Document.Path to resolve against.
func checkDefinitionKinds(doc *ast.Document) []error {
	var errs []error
	for _, d := range doc.Defs {
		if definitionFieldsInvalid(d) {
			errs = append(errs, errAt(doc.Path, d.Line, "definition %q has invalid fields for %s", d.Name, d.Kind.String()))
			continue
		}
		if d.Kind == ast.FileSpecDef && !filepath.IsAbs(d.FilePath) && doc.Path == "" {
			errs = append(errs, errAt(doc.Path, d.Line, "definition %q uses a relative file path but Document.Path is empty", d.Name))
		}
	}
	return errs
}
