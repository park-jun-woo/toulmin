//ff:func feature=codegen type=codegen control=iteration dimension=1
//ff:what formatDefeats — formats defeats slice as Go source expression
package codegen

import "fmt"

// formatDefeats returns a comma-separated, quoted string list for Go source.
func formatDefeats(defeats []string) string {
	parts := make([]string, len(defeats))
	for i, d := range defeats {
		parts[i] = fmt.Sprintf("%q", d)
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}
