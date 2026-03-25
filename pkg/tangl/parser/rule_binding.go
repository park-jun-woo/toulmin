//ff:type feature=tangl type=model
//ff:what RuleBinding — rule registration to a graph
package parser

// RuleBinding represents a rule registration: name is a role of graph using func.
type RuleBinding struct {
	Name      string
	Role      string // rule, counter, except
	Graph     string // graph name (inferred from parent in nested lists)
	Func      string // function reference (e.g., isAuth, policy.IsInRole)
	Specs     []SpecCall
	Qualifier float64 // 0 means default 1.0
	Index     int     // ordered list number (0 means unordered)
	Line      int
}
