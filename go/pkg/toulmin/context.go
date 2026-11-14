//ff:type feature=engine type=model
//ff:what Context — evaluation context interface with key-value access
package toulmin

// Context provides key-value access to evaluation context data.
// Rule functions use Get to retrieve context values set by the caller.
type Context interface {
	Get(key string) (any, bool)
	Set(key string, value any)
}
