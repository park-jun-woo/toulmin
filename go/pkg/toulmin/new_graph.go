//ff:func feature=engine type=engine control=sequence
//ff:what NewGraph — creates a new Graph with the given name
package toulmin

// NewGraph creates a new Graph with the given name.
func NewGraph(name string) *Graph {
	return &Graph{name: name, roles: make(map[string]string)}
}
