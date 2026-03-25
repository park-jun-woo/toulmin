//ff:type feature=tangl type=model
//ff:what EvalDecl — evaluation declaration
package parser

// EvalDecl represents an evaluation: name is results of evaluating graph.
type EvalDecl struct {
	Name     string
	Graph    string
	Trace    bool
	Duration bool
	Line     int
}
