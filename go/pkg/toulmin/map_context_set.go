//ff:func feature=engine type=model control=sequence
//ff:what Set — stores a value in MapContext by key
package toulmin

// Set stores a value in the context.
func (c *MapContext) Set(key string, value any) {
	c.data[key] = value
}
