package validate

import (
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func checkDefinitionKinds(doc *ast.Document) []error {
	var errs []error
	for _, d := range doc.Defs {
		bad := false
		switch d.Kind {
		case ast.ConstDef:
			bad = len(d.Fields) != 0 || len(d.Strings) != 0 || d.LoadType != "" || d.FilePath != ""
		case ast.StructDef:
			bad = d.Value != "" || d.SpecRef != nil || len(d.Strings) != 0 || d.LoadType != "" || d.FilePath != ""
		case ast.StringListDef:
			bad = d.Value != "" || d.SpecRef != nil || len(d.Fields) != 0 || d.LoadType != "" || d.FilePath != ""
		case ast.FileSpecDef:
			bad = d.Value != "" || d.SpecRef != nil || len(d.Fields) != 0 || len(d.Strings) != 0 || d.LoadType == "" || d.FilePath == ""
			if !bad && !filepath.IsAbs(d.FilePath) && doc.Path == "" {
				errs = append(errs, errAt(doc.Path, d.Line, "definition %q uses a relative file path but Document.Path is empty", d.Name))
			}
		default:
			bad = true
		}
		if bad {
			errs = append(errs, errAt(doc.Path, d.Line, "definition %q has invalid fields for %s", d.Name, d.Kind.String()))
		}
	}
	return errs
}

func specCapable(d ast.Definition) bool {
	return d.Kind == ast.StringListDef || d.Kind == ast.FileSpecDef || (d.Kind == ast.ConstDef && d.SpecRef != nil)
}
