//ff:type feature=tangl type=model
//ff:what Case — a tangl:Cases entry (judgment graph + attached execution)
package ast

// Case is a single "in case of `name`" block: a judgment graph with
// nodes, defeat edges, and attached execution edges.
type Case struct {
	Name     string    `json:"name"`
	Requires []Require `json:"requires,omitempty"`
	Nodes    []Node    `json:"nodes"`
	Attacks  []Attack  `json:"attacks,omitempty"`
	Execs    []Exec    `json:"execs,omitempty"` // document declaration order (deterministic execution)
	Line     int       `json:"line"`
}
