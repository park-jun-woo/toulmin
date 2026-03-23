//ff:func feature=analyzer type=analyzer control=iteration dimension=1
//ff:what ExtractDefeats — extracts defeat relationships from Go source via AST
package analyzer

import (
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// ExtractDefeats parses a Go source file and extracts GraphBuilder
// defeat relationships. Returns one GraphDef per NewGraph call found.
func ExtractDefeats(path string) ([]toulmin.GraphDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	var graphs []toulmin.GraphDef
	for _, call := range findGraphCalls(f) {
		name := extractGraphName(call)
		if name == "" {
			continue
		}
		dg := &graphCollector{
			def: toulmin.GraphDef{Graph: name},
		}
		collectChain(call, dg)
		graphs = append(graphs, dg.def)
	}
	return graphs, nil
}
