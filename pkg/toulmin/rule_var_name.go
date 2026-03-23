//ff:func feature=codegen type=codegen control=sequence
//ff:what ruleVarName — converts a rule name to a Go variable name
package toulmin

import "strings"

// ruleVarName converts a rule name to a Go variable name (camelCase).
func ruleVarName(name string) string {
	if len(name) == 0 {
		return "r"
	}
	return strings.ToLower(name[:1]) + name[1:]
}
