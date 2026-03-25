//ff:func feature=engine type=model control=sequence
//ff:what Get — retrieves a value from MapContext by key
package toulmin

// Get returns the value for key and whether it was found.
func (c *MapContext) Get(key string) (any, bool) {
	v, ok := c.data[key]
	return v, ok
}
