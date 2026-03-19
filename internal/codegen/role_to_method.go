//ff:func feature=codegen type=codegen control=selection
//ff:what roleToMethod — maps YAML role to Graph method name
package codegen

// roleToMethod maps YAML role to Graph method name.
func roleToMethod(role string) string {
	switch role {
	case "warrant":
		return "Warrant"
	case "rebuttal":
		return "Rebuttal"
	case "defeater":
		return "Defeater"
	default:
		return "Warrant"
	}
}
