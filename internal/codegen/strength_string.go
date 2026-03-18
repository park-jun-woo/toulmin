//ff:func feature=codegen type=codegen control=selection
//ff:what strengthString — converts Strength to Go expression string
package codegen

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// strengthString returns the Go source expression for a Strength value.
func strengthString(s toulmin.Strength) string {
	switch s {
	case toulmin.Strict:
		return "toulmin.Strict"
	case toulmin.Defeater:
		return "toulmin.Defeater"
	default:
		return "toulmin.Defeasible"
	}
}
