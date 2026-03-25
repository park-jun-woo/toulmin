//ff:func feature=codegen type=codegen control=selection
//ff:what roleToMethod — maps YAML role to Graph method name
package toulmin

// roleToMethod maps YAML role to Graph method name.
func roleToMethod(role string) string {
	switch role {
	case "rule":
		return "Rule"
	case "counter":
		return "Counter"
	case "except":
		return "Except"
	default:
		return "Rule"
	}
}
