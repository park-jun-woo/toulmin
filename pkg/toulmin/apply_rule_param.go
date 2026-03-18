//ff:func feature=annotation type=parser control=selection
//ff:what applyRuleParam — applies a single key=value pair to RuleMeta
package toulmin

import "strconv"

// applyRuleParam sets the appropriate RuleMeta field based on key and value.
func applyRuleParam(meta *RuleMeta, key, value string) {
	switch key {
	case "qualifier":
		meta.Qualifier, _ = strconv.ParseFloat(value, 64)
	case "strength":
		switch value {
		case "strict":
			meta.Strength = Strict
		case "defeasible":
			meta.Strength = Defeasible
		case "defeater":
			meta.Strength = Defeater
		}
	case "defeats":
		meta.Defeats = append(meta.Defeats, value)
	}
}
