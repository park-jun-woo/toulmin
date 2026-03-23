//ff:func feature=codegen type=codegen control=iteration dimension=1
//ff:what graphVarName — converts graph name to PascalCase variable name
package toulmin

import "strings"

// graphVarName converts a graph name to a PascalCase variable name + "Graph".
func graphVarName(name string) string {
	if name == "" {
		return "Graph"
	}
	parts := strings.Split(name, "-")
	var result string
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result + "Graph"
}
