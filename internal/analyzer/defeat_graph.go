//ff:type feature=graph type=model
//ff:what DefeatGraph — extracted defeat relationship from Go source
package analyzer

// DefeatGraph represents a GraphBuilder's defeat relationships
// extracted from Go source code via AST analysis.
type DefeatGraph struct {
	Name    string
	Rules   []string
	Defeats map[string][]string // to → []from
}
