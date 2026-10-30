//ff:func feature=tangl type=validator control=selection
//ff:what definitionFieldsInvalid — whether a definition populates fields that its Kind does not allow
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// definitionFieldsInvalid reports whether d carries fields foreign to its Kind:
// a const has no fields/strings/loader; a struct has no value/spec/strings/loader;
// a string list has no value/spec/fields/loader; a file definition has no
// value/spec/fields/strings but needs both a loader type and a file path. An
// unknown Kind is always invalid.
func definitionFieldsInvalid(d ast.Definition) bool {
	switch d.Kind {
	case ast.ConstDef:
		return len(d.Fields) != 0 || len(d.Strings) != 0 || d.LoadType != "" || d.FilePath != ""
	case ast.StructDef:
		return d.Value != "" || d.SpecRef != nil || len(d.Strings) != 0 || d.LoadType != "" || d.FilePath != ""
	case ast.StringListDef:
		return d.Value != "" || d.SpecRef != nil || len(d.Fields) != 0 || d.LoadType != "" || d.FilePath != ""
	case ast.FileSpecDef:
		return d.Value != "" || d.SpecRef != nil || len(d.Fields) != 0 || len(d.Strings) != 0 || d.LoadType == "" || d.FilePath == ""
	default:
		return true
	}
}
