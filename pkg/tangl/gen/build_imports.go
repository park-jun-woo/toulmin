//ff:func feature=tangl type=codegen control=sequence
//ff:what buildImports — assembles the generated file's full import list
package gen

// buildImports assembles the generated file's import list: the always-
// present toulmin runtime, the tangl runtime and standard library
// packages gated by flags, and every See alias actually used by the
// document.
func buildImports(gc *genContext) ([]importSpec, error) {
	specs := []importSpec{{Alias: "toulmin", Path: "github.com/park-jun-woo/toulmin/pkg/toulmin"}}
	if gc.Flags.NeedsTangl {
		specs = append(specs, importSpec{Alias: "tangl", Path: "github.com/park-jun-woo/toulmin/pkg/tangl"})
	}
	if gc.Flags.NeedsTime {
		specs = append(specs, importSpec{Path: "time"})
	}
	if gc.Flags.NeedsCompare {
		specs = append(specs, importSpec{Path: "fmt"}, importSpec{Path: "strconv"}, importSpec{Path: "strings"})
	}
	aliasSpecs, err := buildAliasImports(gc.Doc)
	if err != nil {
		return nil, err
	}
	specs = append(specs, aliasSpecs...)
	return specs, nil
}
