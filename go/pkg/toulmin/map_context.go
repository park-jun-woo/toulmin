//ff:type feature=engine type=model
//ff:what MapContext — default Context implementation backed by map
package toulmin

// MapContext is the default Context implementation using a map.
type MapContext struct {
	data map[string]any
}
