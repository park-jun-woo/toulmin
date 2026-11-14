//ff:func feature=engine type=engine control=sequence
//ff:what NewContext — creates a new MapContext
package toulmin

// NewContext creates a new empty MapContext.
func NewContext() *MapContext {
	return &MapContext{data: make(map[string]any)}
}
