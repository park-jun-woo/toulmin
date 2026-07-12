//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderPathChainKo — join backtick step names with the 의 connective (`a` 의 `b` 의 `c`)
package ko

import "strings"

// renderPathChainKo renders a data-path / spec-ref step list as its KO surface: each
// step wrapped in backticks and joined by the 의 connective (`a` 의 `b` 의 `c`). It is
// the exact inverse of parsePathKo's 의-chain, so the output re-parses to the same
// steps. It is shared by renderOperandKo for both Path and SpecRef.
func renderPathChainKo(steps []string) string {
	parts := make([]string, len(steps))
	for i, s := range steps {
		parts[i] = "`" + s + "`"
	}
	return strings.Join(parts, " 의 ")
}
