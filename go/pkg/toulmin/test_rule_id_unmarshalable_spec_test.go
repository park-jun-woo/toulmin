//ff:type feature=engine type=model
//ff:what ruleIDUnmarshalableSpec — test helper spec whose fields cannot be JSON-marshaled
package toulmin

type ruleIDUnmarshalableSpec struct {
	Fn func()
}
