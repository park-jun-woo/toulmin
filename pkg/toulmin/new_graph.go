//ff:func feature=engine type=engine control=sequence
//ff:what NewGraph — creates a new GraphBuilder with the given name
package toulmin

// NewGraph creates a new GraphBuilder with the given name.
func NewGraph(name string) *GraphBuilder {
	return &GraphBuilder{name: name, roles: make(map[string]string)}
}
