//ff:func feature=annotation type=parser control=iteration dimension=1
//ff:what parseParams — parses key=value pairs from a rule annotation line
package toulmin

import "strings"

// parseParams splits a space-separated key=value string and applies
// each pair to the given RuleMeta.
func parseParams(meta *RuleMeta, s string) {
	for _, part := range strings.Fields(s) {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		applyRuleParam(meta, kv[0], kv[1])
	}
}
